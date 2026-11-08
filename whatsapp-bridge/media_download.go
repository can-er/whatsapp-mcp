package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/proto/waMmsRetry"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// POST /api/download (fork): the media is looked up in this order
//  1. the copy already downloaded by /api/download (store/<chat>/),
//  2. the copy written by the message events (WHATSAPP_INBOX_DIR/media/<id>.<ext>), if its SHA-256
//     matches the stored one,
//  3. the WhatsApp CDN, with the direct path WhatsApp gave (signed query included); when the CDN
//     answers 403/404/410 the phone is asked to upload the media again (media retry) and the
//     download is retried with the new path.

const mediaRetryTimeout = 30 * time.Second

// mediaRetryWaiter hands the MediaRetry events to the downloads waiting for them
type mediaRetryWaiter struct {
	mu      sync.Mutex
	waiting map[string]chan *events.MediaRetry
}

func newMediaRetryWaiter() *mediaRetryWaiter {
	return &mediaRetryWaiter{waiting: make(map[string]chan *events.MediaRetry)}
}

// register starts waiting for the retry of a message; done must be called when no longer waiting
func (w *mediaRetryWaiter) register(messageID string) (<-chan *events.MediaRetry, func()) {
	ch := make(chan *events.MediaRetry, 1)
	w.mu.Lock()
	w.waiting[messageID] = ch
	w.mu.Unlock()
	return ch, func() {
		w.mu.Lock()
		if w.waiting[messageID] == ch {
			delete(w.waiting, messageID)
		}
		w.mu.Unlock()
	}
}

// deliver passes a MediaRetry event to the download waiting for it, if any
func (w *mediaRetryWaiter) deliver(evt *events.MediaRetry) bool {
	if w == nil || evt == nil {
		return false
	}
	w.mu.Lock()
	ch := w.waiting[evt.MessageID]
	w.mu.Unlock()
	if ch == nil {
		return false
	}
	select {
	case ch <- evt:
		return true
	default:
		return false
	}
}

// storedMedia is the media part of a stored message
type storedMedia struct {
	MediaType     string
	Filename      string
	URL           string
	DirectPath    string
	MediaKey      []byte
	FileSHA256    []byte
	FileEncSHA256 []byte
	FileLength    uint64
	Sender        string
	IsFromMe      bool
}

func (store *MessageStore) GetStoredMedia(id, chatJID string) (storedMedia, error) {
	var m storedMedia
	var mediaType, filename, url, directPath, sender sql.NullString
	var fileLength sql.NullInt64
	err := store.db.QueryRow(
		`SELECT media_type, filename, url, direct_path, media_key, file_sha256, file_enc_sha256, file_length, sender, is_from_me
		FROM messages WHERE id = ? AND chat_jid = ?`, id, chatJID,
	).Scan(&mediaType, &filename, &url, &directPath, &m.MediaKey, &m.FileSHA256, &m.FileEncSHA256, &fileLength, &sender, &m.IsFromMe)
	m.MediaType, m.Filename, m.URL, m.DirectPath, m.Sender = mediaType.String, filename.String, url.String, directPath.String, sender.String
	if fileLength.Valid && fileLength.Int64 > 0 {
		m.FileLength = uint64(fileLength.Int64)
	}
	return m, err
}

// SetDirectPath records the CDN direct path of a stored media message
func (store *MessageStore) SetDirectPath(id, chatJID, directPath string) error {
	if directPath == "" {
		return nil
	}
	_, err := store.db.Exec("UPDATE messages SET direct_path = ? WHERE id = ? AND chat_jid = ?", directPath, id, chatJID)
	return err
}

// mediaDirectPath returns the CDN direct path of the media of a message
func mediaDirectPath(msg *waE2E.Message) string {
	switch {
	case msg.GetImageMessage() != nil:
		return msg.GetImageMessage().GetDirectPath()
	case msg.GetVideoMessage() != nil:
		return msg.GetVideoMessage().GetDirectPath()
	case msg.GetAudioMessage() != nil:
		return msg.GetAudioMessage().GetDirectPath()
	case msg.GetDocumentMessage() != nil:
		return msg.GetDocumentMessage().GetDirectPath()
	}
	return ""
}

// extractDirectPathFromURL returns the direct path of a WhatsApp media URL (rows stored without
// direct_path). The query holds the signature of the path (oh, oe) and must be kept.
func extractDirectPathFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Path == "" {
		return ""
	}
	if u.RawQuery == "" {
		return u.EscapedPath()
	}
	return u.EscapedPath() + "?" + u.RawQuery
}

// mediaFetcher downloads stored media; its network parts are replaceable in tests
type mediaFetcher struct {
	download     func(ctx context.Context, msg whatsmeow.DownloadableMessage) ([]byte, error)
	requestRetry func(ctx context.Context, info *types.MessageInfo, mediaKey []byte) (string, error)
	resolve      lidResolver
	inboxDir     string
	retryTimeout time.Duration
}

func newClientMediaFetcher(client *whatsmeow.Client, retries *mediaRetryWaiter, inboxDir string) *mediaFetcher {
	f := &mediaFetcher{
		download: func(ctx context.Context, msg whatsmeow.DownloadableMessage) ([]byte, error) {
			return client.Download(ctx, msg)
		},
		resolve:      clientLIDResolver(client),
		inboxDir:     inboxDir,
		retryTimeout: mediaRetryTimeout,
	}
	if retries != nil {
		f.requestRetry = func(ctx context.Context, info *types.MessageInfo, mediaKey []byte) (string, error) {
			ch, done := retries.register(info.ID)
			defer done()
			if err := client.SendMediaRetryReceipt(ctx, info, mediaKey); err != nil {
				return "", fmt.Errorf("failed to ask the phone to upload the media again: %w", err)
			}
			select {
			case evt := <-ch:
				notif, err := whatsmeow.DecryptMediaRetryNotification(evt, mediaKey)
				if err != nil {
					return "", fmt.Errorf("media retry failed: %w", err)
				}
				if notif.GetResult() != waMmsRetry.MediaRetryNotification_SUCCESS {
					return "", fmt.Errorf("media retry failed: %s", notif.GetResult())
				}
				return notif.GetDirectPath(), nil
			case <-ctx.Done():
				return "", fmt.Errorf("the phone did not upload the media again: %w", ctx.Err())
			}
		}
	}
	return f
}

// findInboxCopy returns the copy of a message's media written by the message events, if its
// SHA-256 matches the stored one (any copy when no hash is stored)
func (f *mediaFetcher) findInboxCopy(messageID string, wantSHA256 []byte) string {
	idPart := sanitizeFilename(messageID)
	if f.inboxDir == "" || idPart == "" {
		return ""
	}
	dir := filepath.Join(f.inboxDir, "media")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if !e.Type().IsRegular() || !strings.HasPrefix(e.Name(), idPart+".") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		if len(wantSHA256) > 0 {
			sum, err := sha256File(path)
			if err != nil || !bytes.Equal(sum, wantSHA256) {
				continue
			}
		}
		if abs, err := filepath.Abs(path); err == nil {
			return abs
		}
	}
	return ""
}

func sha256File(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func isExpiredMediaError(err error) bool {
	return errors.Is(err, whatsmeow.ErrMediaDownloadFailedWith403) ||
		errors.Is(err, whatsmeow.ErrMediaDownloadFailedWith404) ||
		errors.Is(err, whatsmeow.ErrMediaDownloadFailedWith410)
}

// fetch returns the media type, file name and absolute path of a stored message's media
func (f *mediaFetcher) fetch(ctx context.Context, store *MessageStore, messageID, chatJID string) (string, string, string, error) {
	m, err := store.GetStoredMedia(messageID, chatJID)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to find message: %v", err)
	}
	if m.MediaType == "" {
		return "", "", "", fmt.Errorf("not a media message")
	}

	// Local path, guaranteed to stay inside the chat directory
	chatDir, localPath, err := localMediaPath(chatJID, m.MediaType, messageID, m.Filename)
	if err != nil {
		return "", "", "", err
	}
	filename := filepath.Base(localPath)
	absPath, err := filepath.Abs(localPath)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get absolute path: %v", err)
	}
	if _, err := os.Stat(localPath); err == nil {
		return m.MediaType, filename, absPath, nil
	}
	if inboxPath := f.findInboxCopy(messageID, m.FileSHA256); inboxPath != "" {
		return m.MediaType, filepath.Base(inboxPath), inboxPath, nil
	}

	directPath := m.DirectPath
	if directPath == "" {
		directPath = extractDirectPathFromURL(m.URL)
	}
	if directPath == "" || len(m.MediaKey) == 0 || len(m.FileSHA256) == 0 || len(m.FileEncSHA256) == 0 || m.FileLength == 0 {
		return "", "", "", fmt.Errorf("incomplete media information for download")
	}
	var waMediaType whatsmeow.MediaType
	switch m.MediaType {
	case "image":
		waMediaType = whatsmeow.MediaImage
	case "video":
		waMediaType = whatsmeow.MediaVideo
	case "audio":
		waMediaType = whatsmeow.MediaAudio
	case "document":
		waMediaType = whatsmeow.MediaDocument
	default:
		return "", "", "", fmt.Errorf("unsupported media type: %s", m.MediaType)
	}
	downloader := &MediaDownloader{
		URL:           m.URL,
		DirectPath:    directPath,
		MediaKey:      m.MediaKey,
		FileLength:    m.FileLength,
		FileSHA256:    m.FileSHA256,
		FileEncSHA256: m.FileEncSHA256,
		MediaType:     waMediaType,
	}

	data, err := f.download(ctx, downloader)
	if err != nil && isExpiredMediaError(err) && f.requestRetry != nil {
		// The CDN copy expired: ask the phone that sent or received the message to upload it again
		chat, perr := types.ParseJID(chatJID)
		if perr != nil {
			return "", "", "", fmt.Errorf("failed to download media: %v", err)
		}
		info := &types.MessageInfo{
			MessageSource: types.MessageSource{
				Chat:     chat,
				Sender:   storedSenderJID(m.Sender, f.resolve),
				IsFromMe: m.IsFromMe,
				IsGroup:  chat.Server == types.GroupServer,
			},
			ID: messageID,
		}
		retryCtx, cancel := context.WithTimeout(ctx, f.retryTimeout)
		newPath, rerr := f.requestRetry(retryCtx, info, m.MediaKey)
		cancel()
		if rerr != nil {
			return "", "", "", fmt.Errorf("failed to download media: %v; %v", err, rerr)
		}
		if newPath == "" {
			return "", "", "", fmt.Errorf("failed to download media: %v; media retry returned no path", err)
		}
		if serr := store.SetDirectPath(messageID, chatJID, newPath); serr != nil {
			return "", "", "", fmt.Errorf("failed to store the new media path: %v", serr)
		}
		downloader.DirectPath = newPath
		data, err = f.download(ctx, downloader)
	}
	if err != nil {
		return "", "", "", fmt.Errorf("failed to download media: %v", err)
	}

	if err := ensurePrivateDir(chatDir); err != nil {
		return "", "", "", fmt.Errorf("failed to create chat directory: %v", err)
	}
	if err := os.WriteFile(localPath, data, 0600); err != nil {
		return "", "", "", fmt.Errorf("failed to save media file: %v", err)
	}
	return m.MediaType, filename, absPath, nil
}
