//go:build cgo

package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

const (
	fetchGroup      = "123456789012345678@g.us"
	fetchStoredPath = "/v/t62.7118-24/stored_n.enc?ccb=11-4&oh=01_STORED&oe=6A000000"
	fetchNewPath    = "/v/t62.7118-24/renewed_n.enc?ccb=11-4&oh=01_RENEWED&oe=6B000000"
)

var fetchContent = []byte("fictitious jpeg bytes")

type fakeCDN struct {
	paths   []string
	answers []error // one per call; nil = success
	retries []*types.MessageInfo
	retry   func() (string, error)
}

func (c *fakeCDN) fetcher(inboxDir string) *mediaFetcher {
	return &mediaFetcher{
		download: func(ctx context.Context, msg whatsmeow.DownloadableMessage) ([]byte, error) {
			c.paths = append(c.paths, msg.GetDirectPath())
			var err error
			if len(c.answers) > 0 {
				err, c.answers = c.answers[0], c.answers[1:]
			}
			if err != nil {
				return nil, err
			}
			return fetchContent, nil
		},
		requestRetry: func(ctx context.Context, info *types.MessageInfo, mediaKey []byte) (string, error) {
			c.retries = append(c.retries, info)
			if _, ok := ctx.Deadline(); !ok {
				return "", errors.New("media retry without deadline")
			}
			return c.retry()
		},
		inboxDir:     inboxDir,
		retryTimeout: time.Second,
	}
}

func newFetchStore(t *testing.T, directPath string) *MessageStore {
	t.Helper()
	t.Chdir(t.TempDir())
	store, err := NewMessageStore()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	sum := sha256.Sum256(fetchContent)
	store.StoreChat(fetchGroup, "Test Group", time.Now())
	err = store.StoreMessage("IMG1", fetchGroup, "15550000021", "", time.Now(), false, "image", "image_1.jpg",
		"https://mmg.whatsapp.net/v/t62.7118-24/from_url_n.enc?ccb=11-4&oh=01_URL&oe=6A000000&mms3=true",
		[]byte("media-key"), sum[:], []byte("enc-sha"), uint64(len(fetchContent)))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SetDirectPath("IMG1", fetchGroup, directPath); err != nil {
		t.Fatal(err)
	}
	return store
}

func storedDirectPath(t *testing.T, store *MessageStore) string {
	t.Helper()
	m, err := store.GetStoredMedia("IMG1", fetchGroup)
	if err != nil {
		t.Fatal(err)
	}
	return m.DirectPath
}

func TestMediaDownloadUsesStoredDirectPath(t *testing.T) {
	store := newFetchStore(t, fetchStoredPath)
	cdn := &fakeCDN{}
	mediaType, _, path, err := cdn.fetcher("").fetch(context.Background(), store, "IMG1", fetchGroup)
	if err != nil || mediaType != "image" {
		t.Fatalf("fetch: %q %v", mediaType, err)
	}
	if len(cdn.paths) != 1 || cdn.paths[0] != fetchStoredPath {
		t.Fatalf("downloaded with %q", cdn.paths)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != string(fetchContent) {
		t.Fatalf("saved file: %v", err)
	}
	if runtime.GOOS != "windows" {
		if info, _ := os.Stat(path); info.Mode().Perm() != 0600 {
			t.Fatalf("file mode %v", info.Mode().Perm())
		}
	}
	// A second request is served from the saved copy
	if _, _, again, err := cdn.fetcher("").fetch(context.Background(), store, "IMG1", fetchGroup); err != nil || again != path {
		t.Fatalf("second fetch: %q %v", again, err)
	}
}

func TestMediaDownloadFallsBackToURLWithQuery(t *testing.T) {
	store := newFetchStore(t, "") // row stored by an earlier version: no direct_path
	cdn := &fakeCDN{}
	if _, _, _, err := cdn.fetcher("").fetch(context.Background(), store, "IMG1", fetchGroup); err != nil {
		t.Fatal(err)
	}
	if want := "/v/t62.7118-24/from_url_n.enc?ccb=11-4&oh=01_URL&oe=6A000000&mms3=true"; cdn.paths[0] != want {
		t.Fatalf("downloaded with %q, want %q", cdn.paths[0], want)
	}
}

func TestMediaDownloadFromInbox(t *testing.T) {
	store := newFetchStore(t, fetchStoredPath)
	inbox := t.TempDir()
	if err := os.MkdirAll(filepath.Join(inbox, "media"), 0750); err != nil {
		t.Fatal(err)
	}
	inboxFile := filepath.Join(inbox, "media", "IMG1.jpg")
	if err := os.WriteFile(inboxFile, fetchContent, 0640); err != nil {
		t.Fatal(err)
	}
	cdn := &fakeCDN{answers: []error{whatsmeow.ErrMediaDownloadFailedWith403}}
	_, _, path, err := cdn.fetcher(inbox).fetch(context.Background(), store, "IMG1", fetchGroup)
	if err != nil {
		t.Fatal(err)
	}
	if path != inboxFile || len(cdn.paths) != 0 {
		t.Fatalf("path %q, CDN calls %d", path, len(cdn.paths))
	}

	// Another file under the same name is not this media
	if err := os.WriteFile(inboxFile, []byte("other bytes"), 0640); err != nil {
		t.Fatal(err)
	}
	cdn = &fakeCDN{}
	_, _, path, err = cdn.fetcher(inbox).fetch(context.Background(), store, "IMG1", fetchGroup)
	if err != nil || path == inboxFile || len(cdn.paths) != 1 {
		t.Fatalf("path %q, CDN calls %d, err %v", path, len(cdn.paths), err)
	}
}

func TestMediaDownloadRetriesExpiredMedia(t *testing.T) {
	for _, expired := range []error{whatsmeow.ErrMediaDownloadFailedWith403, whatsmeow.ErrMediaDownloadFailedWith404, whatsmeow.ErrMediaDownloadFailedWith410} {
		t.Run(expired.Error(), func(t *testing.T) {
			store := newFetchStore(t, fetchStoredPath)
			cdn := &fakeCDN{answers: []error{expired, nil}, retry: func() (string, error) { return fetchNewPath, nil }}
			if _, _, _, err := cdn.fetcher("").fetch(context.Background(), store, "IMG1", fetchGroup); err != nil {
				t.Fatal(err)
			}
			if len(cdn.retries) != 1 || len(cdn.paths) != 2 || cdn.paths[1] != fetchNewPath {
				t.Fatalf("retries %d, paths %q", len(cdn.retries), cdn.paths)
			}
			info := cdn.retries[0]
			if info.ID != "IMG1" || info.Chat.String() != fetchGroup || !info.IsGroup || info.IsFromMe ||
				info.Sender != types.NewJID("15550000021", types.DefaultUserServer) {
				t.Fatalf("retry receipt for %+v", info)
			}
			if got := storedDirectPath(t, store); got != fetchNewPath {
				t.Fatalf("stored direct path %q", got)
			}
		})
	}
}

func TestMediaDownloadRetryFailures(t *testing.T) {
	store := newFetchStore(t, fetchStoredPath)

	// The phone no longer has the media
	cdn := &fakeCDN{answers: []error{whatsmeow.ErrMediaDownloadFailedWith403}, retry: func() (string, error) {
		return "", whatsmeow.ErrMediaNotAvailableOnPhone
	}}
	_, _, _, err := cdn.fetcher("").fetch(context.Background(), store, "IMG1", fetchGroup)
	if err == nil || !strings.Contains(err.Error(), "403") || !strings.Contains(err.Error(), "no longer available on phone") {
		t.Fatalf("error %v", err)
	}
	if got := storedDirectPath(t, store); got != fetchStoredPath {
		t.Fatalf("direct path changed to %q", got)
	}

	// Errors other than an expired copy are not retried
	cdn = &fakeCDN{answers: []error{whatsmeow.ErrInvalidMediaHMAC}, retry: func() (string, error) { return fetchNewPath, nil }}
	if _, _, _, err := cdn.fetcher("").fetch(context.Background(), store, "IMG1", fetchGroup); err == nil || len(cdn.retries) != 0 {
		t.Fatalf("err %v, retries %d", err, len(cdn.retries))
	}

	// Without a WhatsApp connection to ask (no retry function), the CDN error is returned
	f := (&fakeCDN{answers: []error{whatsmeow.ErrMediaDownloadFailedWith410}}).fetcher("")
	f.requestRetry = nil
	if _, _, _, err := f.fetch(context.Background(), store, "IMG1", fetchGroup); err == nil || !strings.Contains(err.Error(), "410") {
		t.Fatalf("err %v", err)
	}

	// Unknown and non-media messages
	if _, _, _, err := cdn.fetcher("").fetch(context.Background(), store, "NOPE", fetchGroup); err == nil {
		t.Fatal("unknown message downloaded")
	}
}

func TestDirectPathStoredFromLiveMessage(t *testing.T) {
	store := newCaptionStore(t)
	group, _ := types.ParseJID(captionGroup)
	evt := &events.Message{
		Info: types.MessageInfo{
			MessageSource: types.MessageSource{Chat: group, Sender: types.NewJID("15550000022", types.DefaultUserServer), IsGroup: true},
			ID:            "IMG2",
			Timestamp:     time.Now(),
		},
		Message: &waE2E.Message{ImageMessage: &waE2E.ImageMessage{DirectPath: proto.String(fetchStoredPath), Mimetype: proto.String("image/jpeg")}},
	}
	handleMessage(nil, store, evt, parseChatAllowlist(captionGroup), nil, waLog.Noop)
	m, err := store.GetStoredMedia("IMG2", captionGroup)
	if err != nil || m.DirectPath != fetchStoredPath {
		t.Fatalf("direct_path %q, %v", m.DirectPath, err)
	}
}
