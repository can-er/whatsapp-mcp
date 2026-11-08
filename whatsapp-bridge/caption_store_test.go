//go:build cgo

package main

import (
	"testing"
	"time"

	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/proto/waHistorySync"
	"go.mau.fi/whatsmeow/proto/waWeb"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

const captionGroup = "123456789012345678@g.us"

func newCaptionStore(t *testing.T) *MessageStore {
	t.Helper()
	t.Chdir(t.TempDir())
	store, err := NewMessageStore()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	// A known chat name: GetChatName then never needs a WhatsApp client
	if err := store.StoreChat(captionGroup, "Test Group", time.Now()); err != nil {
		t.Fatal(err)
	}
	return store
}

func storedContent(t *testing.T, store *MessageStore, id string) (content, mediaType string) {
	t.Helper()
	err := store.db.QueryRow("SELECT content, media_type FROM messages WHERE id = ? AND chat_jid = ?", id, captionGroup).Scan(&content, &mediaType)
	if err != nil {
		t.Fatalf("message %s not stored: %v", id, err)
	}
	return content, mediaType
}

func TestCaptionStoredFromLiveMessage(t *testing.T) {
	store := newCaptionStore(t)
	group, _ := types.ParseJID(captionGroup)
	for id, msg := range map[string]*waE2E.Message{
		"IMG": {ImageMessage: &waE2E.ImageMessage{Caption: proto.String("2. Package management"), Mimetype: proto.String("image/jpeg")}},
		"VID": {VideoMessage: &waE2E.VideoMessage{Caption: proto.String("screen recording"), Mimetype: proto.String("video/mp4")}},
		"DOC": {DocumentMessage: &waE2E.DocumentMessage{Caption: proto.String("invoice"), FileName: proto.String("a.pdf")}},
		"RAW": {ImageMessage: &waE2E.ImageMessage{Mimetype: proto.String("image/jpeg")}},
	} {
		evt := &events.Message{
			Info: types.MessageInfo{
				MessageSource: types.MessageSource{Chat: group, Sender: types.NewJID("15550000003", types.DefaultUserServer), IsGroup: true},
				ID:            id,
				Timestamp:     time.Date(2026, 11, 5, 8, 29, 20, 0, time.UTC),
			},
			Message: msg,
		}
		handleMessage(nil, store, evt, parseChatAllowlist(captionGroup), nil, waLog.Noop)
	}
	for id, want := range map[string][2]string{
		"IMG": {"2. Package management", "image"},
		"VID": {"screen recording", "video"},
		"DOC": {"invoice", "document"},
		"RAW": {"", "image"},
	} {
		content, mediaType := storedContent(t, store, id)
		if content != want[0] || mediaType != want[1] {
			t.Errorf("%s: content %q media %q, want %q %q", id, content, mediaType, want[0], want[1])
		}
	}
}

func TestCaptionStoredFromHistorySync(t *testing.T) {
	store := newCaptionStore(t)
	sync := &events.HistorySync{Data: &waHistorySync.HistorySync{
		SyncType: waHistorySync.HistorySync_RECENT.Enum(),
		Conversations: []*waHistorySync.Conversation{{
			ID:   proto.String(captionGroup),
			Name: proto.String("Test Group"),
			Messages: []*waHistorySync.HistorySyncMsg{{
				Message: &waWeb.WebMessageInfo{
					Key: &waCommon.MessageKey{
						RemoteJID:   proto.String(captionGroup),
						FromMe:      proto.Bool(false),
						ID:          proto.String("HIST1"),
						Participant: proto.String("15550000003@s.whatsapp.net"),
					},
					MessageTimestamp: proto.Uint64(uint64(time.Date(2026, 11, 7, 14, 48, 0, 0, time.UTC).Unix())),
					Message:          &waE2E.Message{VideoMessage: &waE2E.VideoMessage{Caption: proto.String("history caption")}},
				},
			}},
		}},
	}}
	handleHistorySync(nil, store, sync, parseChatAllowlist(captionGroup), waLog.Noop)
	if content, mediaType := storedContent(t, store, "HIST1"); content != "history caption" || mediaType != "video" {
		t.Fatalf("content %q media %q", content, mediaType)
	}
}
