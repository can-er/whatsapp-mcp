//go:build cgo

package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func newTestContainer(t *testing.T) *sqlstore.Container {
	t.Helper()
	dsn := "file:" + filepath.Join(t.TempDir(), "whatsapp.db") + "?_foreign_keys=on"
	container, err := sqlstore.New(context.Background(), "sqlite3", dsn, waLog.Noop)
	if err != nil {
		t.Fatalf("sqlstore.New: %v", err)
	}
	t.Cleanup(func() { container.Close() })
	return container
}

// The session store opens with the current whatsmeow API and yields a usable (unpaired) client,
// exactly as main() does, without connecting to WhatsApp.
func TestSQLStoreOpens(t *testing.T) {
	device, err := newTestContainer(t).GetFirstDevice(context.Background())
	if err != nil || device == nil {
		t.Fatalf("GetFirstDevice: %v", err)
	}
	client := whatsmeow.NewClient(device, waLog.Noop)
	if client == nil || client.IsConnected() || client.Store.ID != nil {
		t.Fatal("expected a new, unpaired and disconnected client")
	}
}

// Chats addressed by LID are allowed when the phone number they map to is listed.
func TestChatAllowlistResolvesLID(t *testing.T) {
	ctx := context.Background()
	container := newTestContainer(t)
	device := container.NewDevice()
	device.LIDs = container.LIDMap

	lid := types.NewJID("123456789012345", types.HiddenUserServer)
	pn := types.NewJID("15550000001", types.DefaultUserServer)
	if err := device.LIDs.PutLIDMapping(ctx, lid, pn); err != nil {
		t.Fatal(err)
	}
	client := whatsmeow.NewClient(device, waLog.Noop)
	a := parseChatAllowlist("15550000001")

	if !a.AllowsChat(client, lid) || !a.AllowsChatString(client, lid.String()) {
		t.Error("LID of an allowed phone number must be allowed")
	}
	if a.AllowsChat(client, types.NewJID("999999999999999", types.HiddenUserServer)) {
		t.Error("unknown LID must not be allowed")
	}
	if a.AllowsChat(nil, lid) {
		t.Error("without a device store the LID cannot be resolved and must not be allowed")
	}
}

// An edit replaces the stored text of the original message
func TestUpdateMessageContent(t *testing.T) {
	t.Chdir(t.TempDir())
	store, err := NewMessageStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	chat := "15550000001@s.whatsapp.net"
	if err := store.StoreChat(chat, "Test", time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := store.StoreMessage("3EB0ORIGINAL", chat, "15550000001", "typo", time.Now(), false, "", "", "", nil, nil, nil, 0); err != nil {
		t.Fatal(err)
	}
	if err := store.UpdateMessageContent("3EB0ORIGINAL", chat, "fixed"); err != nil {
		t.Fatal(err)
	}
	msgs, err := store.GetMessages(chat, 10)
	if err != nil || len(msgs) != 1 || msgs[0].Content != "fixed" {
		t.Fatalf("messages after edit: %+v, %v", msgs, err)
	}
}
