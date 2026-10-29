package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// GET /api/history (fork-specific, ADR-0006): stored messages of one allowed chat, newest first,
// in the same format as message events, paginated with an opaque cursor.

const (
	defaultHistoryLimit = 50
	maxHistoryLimit     = 500
)

type historyResponse struct {
	Messages   []*MessageEvent `json:"messages"`
	NextBefore *string         `json:"next_before"`
}

// historyCursor positions a page: messages strictly older than (Time, ID)
type historyCursor struct {
	Time time.Time
	ID   string
}

func (c historyCursor) encode() string {
	return base64.RawURLEncoding.EncodeToString([]byte(c.Time.UTC().Format(time.RFC3339Nano) + "|" + c.ID))
}

// parseHistoryCursor accepts an RFC 3339 time or a next_before value; "" means the newest messages
func parseHistoryCursor(raw string) (*historyCursor, error) {
	if raw == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339Nano, raw); err == nil {
		// Every message of that instant is older than an ID that sorts after all others
		return &historyCursor{Time: t, ID: "\U0010FFFF"}, nil
	}
	data, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid before cursor")
	}
	ts, id, ok := strings.Cut(string(data), "|")
	t, err := time.Parse(time.RFC3339Nano, ts)
	if !ok || err != nil {
		return nil, fmt.Errorf("invalid before cursor")
	}
	return &historyCursor{Time: t, ID: id}, nil
}

type storedMessage struct {
	ID, ChatJID, Sender, Content, MediaType, Filename, ChatName string
	Timestamp                                                   time.Time
	IsFromMe                                                    bool
}

// HistoryPage returns up to limit messages of the given chat JIDs older than the cursor, newest first,
// and the cursor of the next page (nil on the last page).
func (store *MessageStore) HistoryPage(chatJIDs []string, before *historyCursor, limit int) ([]storedMessage, *historyCursor, error) {
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(chatJIDs)), ",")
	query := `SELECT m.id, m.chat_jid, COALESCE(m.sender, ''), COALESCE(m.content, ''), m.timestamp, m.is_from_me,
			COALESCE(m.media_type, ''), COALESCE(m.filename, ''), COALESCE(c.name, '')
		FROM messages m LEFT JOIN chats c ON c.jid = m.chat_jid
		WHERE m.chat_jid IN (` + placeholders + `)`
	args := make([]interface{}, 0, len(chatJIDs)+4)
	for _, jid := range chatJIDs {
		args = append(args, jid)
	}
	// julianday() compares instants even if timestamps were stored with different UTC offsets
	if before != nil {
		ts := before.Time.UTC().Format("2006-01-02T15:04:05.999999999Z")
		query += ` AND (julianday(m.timestamp) < julianday(?) OR (julianday(m.timestamp) = julianday(?) AND m.id < ?))`
		args = append(args, ts, ts, before.ID)
	}
	query += ` ORDER BY julianday(m.timestamp) DESC, m.id DESC LIMIT ?`
	args = append(args, limit+1)

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var page []storedMessage
	for rows.Next() {
		var m storedMessage
		if err := rows.Scan(&m.ID, &m.ChatJID, &m.Sender, &m.Content, &m.Timestamp, &m.IsFromMe, &m.MediaType, &m.Filename, &m.ChatName); err != nil {
			return nil, nil, err
		}
		page = append(page, m)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	var next *historyCursor
	if len(page) > limit {
		page = page[:limit]
		last := page[len(page)-1]
		next = &historyCursor{Time: last.Timestamp, ID: last.ID}
	}
	return page, next, nil
}

// historyChatJIDs returns the stored forms of a chat: its phone-number JID and, when known, its LID
func historyChatJIDs(client *whatsmeow.Client, chat string) []string {
	normalized, ok := normalizeChatJID(chat)
	if !ok {
		return nil
	}
	jids := []string{normalized}
	jid, err := types.ParseJID(normalized)
	if err != nil || client == nil || client.Store == nil || client.Store.LIDs == nil {
		return jids
	}
	switch jid.Server {
	case types.DefaultUserServer:
		if lid, err := client.Store.LIDs.GetLIDForPN(context.Background(), jid); err == nil && !lid.IsEmpty() {
			jids = append(jids, lid.ToNonAD().String())
		}
	case types.HiddenUserServer:
		if pn, err := client.Store.LIDs.GetPNForLID(context.Background(), jid); err == nil && !pn.IsEmpty() {
			jids = append(jids, pn.ToNonAD().String())
		}
	}
	return jids
}

// storedSenderJID rebuilds the sender JID: rows only hold the user part (phone number or LID)
func storedSenderJID(sender string, resolve lidResolver) types.JID {
	if strings.Contains(sender, "@") {
		if jid, err := types.ParseJID(sender); err == nil {
			return jid.ToNonAD()
		}
	}
	user := strings.SplitN(sender, ":", 2)[0]
	if lid := types.NewJID(user, types.HiddenUserServer); resolve != nil {
		if _, ok := resolve(lid); ok {
			return lid
		}
	}
	return types.NewJID(user, types.DefaultUserServer)
}

// historyEvent converts a stored message to the message event format
func historyEvent(m storedMessage, resolve lidResolver) *MessageEvent {
	chat, _ := types.ParseJID(m.ChatJID)
	chatPN, chatLID := splitPNAndLID(chat, types.EmptyJID, resolve)
	chatJID := m.ChatJID
	if !chatPN.IsEmpty() {
		chatJID = chatPN.String()
	}
	sender := storedSenderJID(m.Sender, resolve)
	senderPN, senderLID := splitPNAndLID(sender, types.EmptyJID, resolve)

	ev := &MessageEvent{
		Version:   eventSchemaVersion,
		Type:      "message",
		ID:        m.ID,
		ChatJID:   chatJID,
		ChatName:  m.ChatName,
		SenderJID: sender.String(),
		IsFromMe:  m.IsFromMe,
		Timestamp: m.Timestamp.UTC(),
		Text:      m.Content,
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

	if m.MediaType != "" {
		media := &EventMedia{Type: m.MediaType, Status: "not_downloaded"}
		if m.MediaType == "document" {
			media.OriginalFilename = m.Filename
		}
		if _, path, err := localMediaPath(m.ChatJID, m.MediaType, m.ID, m.Filename); err == nil {
			if abs, err := filepath.Abs(path); err == nil {
				if size, sum, err := fileSHA256(abs); err == nil {
					media.Status, media.Path, media.Size, media.SHA256 = "ok", &abs, size, sum
				}
			}
		}
		ev.Media = media
	}
	return ev
}

func fileSHA256(path string) (int64, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, "", err
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return 0, "", err
	}
	return n, hex.EncodeToString(h.Sum(nil)), nil
}
