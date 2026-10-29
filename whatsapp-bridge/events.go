package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// Event emission for local automation (fork-specific): every real-time message of an allowed chat is
// written as a JSON file to WHATSAPP_INBOX_DIR (media first) and announced to WHATSAPP_EVENT_WEBHOOK_URL.
// The format is documented in the README ("Message events").

const eventSchemaVersion = 1

// EventMedia describes the media attached to a message event
type EventMedia struct {
	Status           string  `json:"status"` // "ok" or "failed"
	Error            string  `json:"error,omitempty"`
	Type             string  `json:"type"` // image, video, audio, document
	Mime             string  `json:"mime"`
	Path             *string `json:"path"`
	Size             int64   `json:"size"`
	SHA256           string  `json:"sha256,omitempty"`
	Voice            bool    `json:"voice"`
	DurationS        *uint32 `json:"duration_s"`
	OriginalFilename string  `json:"original_filename,omitempty"` // informative only, never used as a path
}

// MessageEvent is the JSON document written to the inbox and posted to the webhook
type MessageEvent struct {
	Version         int         `json:"version"`
	Type            string      `json:"type"`
	ID              string      `json:"id"`
	ChatJID         string      `json:"chat_jid"`
	ChatLID         *string     `json:"chat_lid"`
	ChatName        string      `json:"chat_name"`
	SenderJID       string      `json:"sender_jid"`
	SenderPhone     *string     `json:"sender_phone"`
	SenderLID       *string     `json:"sender_lid"`
	SenderName      string      `json:"sender_name"`
	IsFromMe        bool        `json:"is_from_me"`
	Timestamp       time.Time   `json:"timestamp"`
	ReceivedAt      time.Time   `json:"received_at"`
	Text            string      `json:"text"`
	QuotedMessageID *string     `json:"quoted_message_id"`
	Media           *EventMedia `json:"media"`
	// Only for type edit, revoke and reaction: the message the update refers to
	TargetMessageID *string `json:"target_message_id,omitempty"`
	// Only for type reaction: the emoji, empty when the reaction was removed
	Reaction *string `json:"reaction,omitempty"`
}

// messageUpdate detects edits, deletions and reactions, which refer to an earlier message.
// kind is "" for a regular message.
func messageUpdate(msg *waE2E.Message) (kind, target, text string, reaction *string) {
	if pm := msg.GetProtocolMessage(); pm != nil {
		switch pm.GetType() {
		case waE2E.ProtocolMessage_MESSAGE_EDIT:
			return "edit", pm.GetKey().GetID(), messageText(pm.GetEditedMessage()), nil
		case waE2E.ProtocolMessage_REVOKE:
			return "revoke", pm.GetKey().GetID(), "", nil
		}
	}
	if r := msg.GetReactionMessage(); r != nil {
		emoji := r.GetText()
		return "reaction", r.GetKey().GetID(), "", &emoji
	}
	return "", "", "", nil
}

// lidResolver maps a LID JID to the phone-number JID it belongs to, if known
type lidResolver func(lid types.JID) (types.JID, bool)

func clientLIDResolver(client *whatsmeow.Client) lidResolver {
	return func(lid types.JID) (types.JID, bool) {
		if client == nil || client.Store == nil || client.Store.LIDs == nil {
			return types.EmptyJID, false
		}
		pn, err := client.Store.LIDs.GetPNForLID(context.Background(), lid.ToNonAD())
		return pn, err == nil && !pn.IsEmpty()
	}
}

// splitPNAndLID returns the phone-number form and the LID form of a user JID, when known
func splitPNAndLID(jid, alt types.JID, resolve lidResolver) (pn, lid types.JID) {
	for _, j := range []types.JID{jid, alt} {
		switch j.Server {
		case types.DefaultUserServer:
			if pn.IsEmpty() {
				pn = j.ToNonAD()
			}
		case types.HiddenUserServer:
			if lid.IsEmpty() {
				lid = j.ToNonAD()
			}
		}
	}
	if pn.IsEmpty() && !lid.IsEmpty() && resolve != nil {
		if resolved, ok := resolve(lid); ok {
			pn = resolved.ToNonAD()
		}
	}
	return pn, lid
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// messageMedia returns the downloadable part of a message and its event description (without download result)
func messageMedia(msg *waE2E.Message) (whatsmeow.DownloadableMessage, *EventMedia) {
	duration := func(s uint32) *uint32 {
		if s == 0 {
			return nil
		}
		return &s
	}
	switch {
	case msg.GetImageMessage() != nil:
		m := msg.GetImageMessage()
		return m, &EventMedia{Type: "image", Mime: m.GetMimetype(), Size: int64(m.GetFileLength())}
	case msg.GetVideoMessage() != nil:
		m := msg.GetVideoMessage()
		return m, &EventMedia{Type: "video", Mime: m.GetMimetype(), Size: int64(m.GetFileLength()), DurationS: duration(m.GetSeconds())}
	case msg.GetAudioMessage() != nil:
		m := msg.GetAudioMessage()
		return m, &EventMedia{Type: "audio", Mime: m.GetMimetype(), Size: int64(m.GetFileLength()), Voice: m.GetPTT(), DurationS: duration(m.GetSeconds())}
	case msg.GetDocumentMessage() != nil:
		m := msg.GetDocumentMessage()
		return m, &EventMedia{Type: "document", Mime: m.GetMimetype(), Size: int64(m.GetFileLength()), OriginalFilename: m.GetFileName()}
	}
	return nil, nil
}

func messageText(msg *waE2E.Message) string {
	switch {
	case msg.GetConversation() != "":
		return msg.GetConversation()
	case msg.GetExtendedTextMessage() != nil:
		return msg.GetExtendedTextMessage().GetText()
	case msg.GetImageMessage() != nil:
		return msg.GetImageMessage().GetCaption()
	case msg.GetVideoMessage() != nil:
		return msg.GetVideoMessage().GetCaption()
	case msg.GetDocumentMessage() != nil:
		return msg.GetDocumentMessage().GetCaption()
	}
	return ""
}

func quotedMessageID(msg *waE2E.Message) string {
	for _, ci := range []*waE2E.ContextInfo{
		msg.GetExtendedTextMessage().GetContextInfo(),
		msg.GetImageMessage().GetContextInfo(),
		msg.GetVideoMessage().GetContextInfo(),
		msg.GetAudioMessage().GetContextInfo(),
		msg.GetDocumentMessage().GetContextInfo(),
	} {
		if id := ci.GetStanzaID(); id != "" {
			return id
		}
	}
	return ""
}

var mimeExtensions = map[string]string{
	"image/jpeg": "jpg", "image/png": "png", "image/webp": "webp", "image/gif": "gif",
	"video/mp4": "mp4", "video/3gpp": "3gp", "video/quicktime": "mov",
	"audio/ogg": "ogg", "audio/mpeg": "mp3", "audio/mp4": "m4a", "audio/aac": "aac", "audio/amr": "amr",
	"application/pdf": "pdf",
}

// mediaExtension picks a file extension from the MIME type, then from the sender's file name
func mediaExtension(media *EventMedia) string {
	base := strings.TrimSpace(strings.SplitN(media.Mime, ";", 2)[0])
	if ext, ok := mimeExtensions[strings.ToLower(base)]; ok {
		return ext
	}
	ext := strings.TrimPrefix(filepath.Ext(sanitizeFilename(media.OriginalFilename)), ".")
	if ext != "" && len(ext) <= 8 && strings.Trim(strings.ToLower(ext), "abcdefghijklmnopqrstuvwxyz0123456789") == "" {
		return strings.ToLower(ext)
	}
	return "bin"
}

// buildMessageEvent converts a whatsmeow message into an event; the media download is done separately.
func buildMessageEvent(evt *events.Message, chatName string, resolve lidResolver, receivedAt time.Time) (*MessageEvent, whatsmeow.DownloadableMessage) {
	info := evt.Info
	chatPN, chatLID := splitPNAndLID(info.Chat, types.EmptyJID, resolve)
	chatJID := info.Chat.ToNonAD()
	if !chatPN.IsEmpty() {
		chatJID = chatPN
	}
	senderPN, senderLID := splitPNAndLID(info.Sender, info.SenderAlt, resolve)

	ev := &MessageEvent{
		Version:         eventSchemaVersion,
		Type:            "message",
		ID:              info.ID,
		ChatJID:         chatJID.String(),
		ChatName:        chatName,
		SenderJID:       info.Sender.ToNonAD().String(),
		SenderName:      info.PushName,
		IsFromMe:        info.IsFromMe,
		Timestamp:       info.Timestamp.UTC(),
		ReceivedAt:      receivedAt.UTC(),
		Text:            messageText(evt.Message),
		QuotedMessageID: optionalString(quotedMessageID(evt.Message)),
	}
	if !chatLID.IsEmpty() {
		ev.ChatLID = optionalString(chatLID.String())
	}
	if !senderPN.IsEmpty() {
		ev.SenderPhone = optionalString(senderPN.User)
	}
	if !senderLID.IsEmpty() {
		ev.SenderLID = optionalString(senderLID.String())
	}

	if kind, target, text, reaction := messageUpdate(evt.Message); kind != "" {
		ev.Type, ev.Text, ev.Reaction = kind, text, reaction
		ev.TargetMessageID = &target
		return ev, nil
	}

	downloadable, media := messageMedia(evt.Message)
	ev.Media = media
	return ev, downloadable
}

// eventSink writes message events to the inbox directory and notifies the webhook
type eventSink struct {
	inboxDir     string
	webhookURL   string
	webhookToken string
	retention    time.Duration
	httpClient   *http.Client
	retryDelay   time.Duration
	download     func(ctx context.Context, msg whatsmeow.DownloadableMessage) ([]byte, error)
	resolve      lidResolver
	now          func() time.Time
	queue        chan eventJob
	logf         func(format string, args ...interface{})
}

type eventJob struct {
	evt      *events.Message
	chatName string
}

// validateWebhookURL only accepts http(s) URLs pointing to the loopback interface
func validateWebhookURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("WHATSAPP_EVENT_WEBHOOK_URL must be an http(s) URL")
	}
	host := u.Hostname()
	if host == "localhost" {
		return nil
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return nil
	}
	return fmt.Errorf("WHATSAPP_EVENT_WEBHOOK_URL must point to localhost or a loopback address")
}

// newEventSinkFromEnv returns nil when neither WHATSAPP_INBOX_DIR nor WHATSAPP_EVENT_WEBHOOK_URL is set
func newEventSinkFromEnv(client *whatsmeow.Client) (*eventSink, error) {
	inboxDir := strings.TrimSpace(os.Getenv("WHATSAPP_INBOX_DIR"))
	webhookURL := strings.TrimSpace(os.Getenv("WHATSAPP_EVENT_WEBHOOK_URL"))
	if inboxDir == "" && webhookURL == "" {
		return nil, nil
	}
	if webhookURL != "" {
		if err := validateWebhookURL(webhookURL); err != nil {
			return nil, err
		}
	}
	retentionDays := 30
	if raw := strings.TrimSpace(os.Getenv("WHATSAPP_INBOX_RETENTION_DAYS")); raw != "" {
		days, err := strconv.Atoi(raw)
		if err != nil || days < 0 {
			return nil, fmt.Errorf("WHATSAPP_INBOX_RETENTION_DAYS must be a number of days (0 disables the purge)")
		}
		retentionDays = days
	}
	if inboxDir != "" {
		abs, err := filepath.Abs(inboxDir)
		if err != nil {
			return nil, fmt.Errorf("invalid WHATSAPP_INBOX_DIR: %v", err)
		}
		inboxDir = abs
	}

	sink := &eventSink{
		inboxDir:     inboxDir,
		webhookURL:   webhookURL,
		webhookToken: strings.TrimSpace(os.Getenv("WHATSAPP_EVENT_WEBHOOK_TOKEN")),
		retention:    time.Duration(retentionDays) * 24 * time.Hour,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		retryDelay:   time.Second,
		download: func(ctx context.Context, msg whatsmeow.DownloadableMessage) ([]byte, error) {
			return client.Download(ctx, msg)
		},
		resolve: clientLIDResolver(client),
		now:     time.Now,
		logf:    func(format string, args ...interface{}) { fmt.Printf(format+"\n", args...) },
	}
	if err := sink.prepareInbox(); err != nil {
		return nil, err
	}
	return sink, nil
}

// prepareInbox creates the inbox sub-directories. The inbox root is never chmod-ed, so a setgid bit
// set by the administrator (to share files with the pipeline's group) is preserved and inherited.
func (s *eventSink) prepareInbox() error {
	if s.inboxDir == "" {
		return nil
	}
	for _, sub := range []string{"events", "media", "tmp"} {
		if err := os.MkdirAll(filepath.Join(s.inboxDir, sub), 0750); err != nil {
			return fmt.Errorf("failed to create inbox directory: %v", err)
		}
	}
	return nil
}

// Start runs the event worker (preserving message order) and the retention purge
func (s *eventSink) Start(ctx context.Context) {
	s.queue = make(chan eventJob, 256)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case job := <-s.queue:
				if _, err := s.Process(ctx, job.evt, job.chatName); err != nil {
					s.logf("Failed to emit event for message %s: %v", job.evt.Info.ID, err)
				}
			}
		}
	}()
	if s.inboxDir != "" && s.retention > 0 {
		go func() {
			ticker := time.NewTicker(time.Hour)
			defer ticker.Stop()
			for {
				s.Purge()
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}
			}
		}()
	}
}

// Enqueue schedules an event without blocking the WhatsApp event handler on media downloads
func (s *eventSink) Enqueue(evt *events.Message, chatName string) {
	if s == nil || s.queue == nil {
		return
	}
	s.queue <- eventJob{evt: evt, chatName: chatName}
}

// writeAtomic writes data to dir/name through the inbox tmp directory, so watchers of dir only
// ever see complete files (inotify IN_MOVED_TO).
func (s *eventSink) writeAtomic(dir, name string, data []byte) (string, error) {
	tmp, err := os.CreateTemp(filepath.Join(s.inboxDir, "tmp"), name+".*.tmp")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())
	if err := tmp.Chmod(0640); err != nil {
		tmp.Close()
		return "", err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	final := filepath.Join(s.inboxDir, dir, name)
	if err := os.Rename(tmp.Name(), final); err != nil {
		return "", err
	}
	return final, nil
}

// Process downloads the media (if any), writes it, then writes the event and notifies the webhook.
// It returns the event as written.
func (s *eventSink) Process(ctx context.Context, evt *events.Message, chatName string) (*MessageEvent, error) {
	ev, downloadable := buildMessageEvent(evt, chatName, s.resolve, s.now())
	idPart := sanitizeFilename(ev.ID)
	if idPart == "" {
		idPart = strconv.FormatInt(ev.ReceivedAt.UnixNano(), 10)
	}

	if ev.Media != nil {
		ev.Media.Status = "ok"
		if s.inboxDir == "" {
			ev.Media.Status, ev.Media.Error = "failed", "WHATSAPP_INBOX_DIR is not set"
		} else if data, err := s.download(ctx, downloadable); err != nil {
			ev.Media.Status, ev.Media.Error = "failed", err.Error()
		} else {
			sum := sha256.Sum256(data)
			path, err := s.writeAtomic("media", idPart+"."+mediaExtension(ev.Media), data)
			if err != nil {
				ev.Media.Status, ev.Media.Error = "failed", err.Error()
			} else {
				ev.Media.Path = &path
				ev.Media.Size = int64(len(data))
				ev.Media.SHA256 = hex.EncodeToString(sum[:])
			}
		}
	}

	body, err := json.Marshal(ev)
	if err != nil {
		return nil, err
	}
	if err := s.deliver(ev.Timestamp.Format("20060102T150405Z")+"_"+idPart+".json", body); err != nil {
		return nil, err
	}
	return ev, nil
}

// deliver writes an event to the inbox (if configured) and posts it to the webhook (if configured)
func (s *eventSink) deliver(name string, body []byte) error {
	if s.inboxDir != "" {
		if _, err := s.writeAtomic("events", name, body); err != nil {
			return fmt.Errorf("failed to write event: %v", err)
		}
	}
	if s.webhookURL != "" {
		go s.notify(body)
	}
	return nil
}

// notify posts the event to the webhook, retrying a few times; the inbox stays the durable record
func (s *eventSink) notify(body []byte) bool {
	delay := s.retryDelay
	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequest(http.MethodPost, s.webhookURL, bytes.NewReader(body))
		if err != nil {
			s.logf("Invalid webhook request: %v", err)
			return false
		}
		req.Header.Set("Content-Type", "application/json")
		if s.webhookToken != "" {
			req.Header.Set("Authorization", "Bearer "+s.webhookToken)
		}
		resp, err := s.httpClient.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return true
			}
			err = fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		s.logf("Webhook attempt %d failed: %v", attempt, err)
		if attempt < 3 {
			time.Sleep(delay)
			delay *= 2
		}
	}
	return false
}

// Purge deletes inbox files older than the retention period (and stale temporary files)
func (s *eventSink) Purge() {
	if s.inboxDir == "" || s.retention <= 0 {
		return
	}
	now := s.now()
	for sub, maxAge := range map[string]time.Duration{"events": s.retention, "media": s.retention, "tmp": time.Hour} {
		dir := filepath.Join(s.inboxDir, sub)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil || !info.Mode().IsRegular() {
				continue
			}
			if now.Sub(info.ModTime()) > maxAge {
				os.Remove(filepath.Join(dir, entry.Name()))
			}
		}
	}
}
