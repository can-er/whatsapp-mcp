//go:build cgo

package main

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// Fictitious identities (ADR-0007)
const (
	testLIDKnown   = "123456789000001" // LID whose phone number is in the LID map
	testLIDUnknown = "123456789000002" // LID without mapping
	testPNOfLID    = "15550000011"
	testPNMapped   = "15550000012" // phone number the LID map knows (as a phone number)
	testPNOwn      = "15550000013"
)

func testSenderIdentity() senderIdentity {
	return senderIdentity{
		pnForLID: func(lid types.JID) (types.JID, bool) {
			if lid.User == testLIDKnown && lid.Server == types.HiddenUserServer {
				return types.NewJID(testPNOfLID, types.DefaultUserServer), true
			}
			return types.EmptyJID, false
		},
		lidForPN: func(pn types.JID) (types.JID, bool) {
			if pn.User == testPNMapped && pn.Server == types.DefaultUserServer {
				return types.NewJID("123456789000003", types.HiddenUserServer), true
			}
			return types.EmptyJID, false
		},
		own: types.NewJID(testPNOwn, types.DefaultUserServer),
	}
}

func TestSenderPhoneOfMessage(t *testing.T) {
	ident := testSenderIdentity()
	lid := func(u string) types.JID { return types.NewJID(u, types.HiddenUserServer) }
	pn := func(u string) types.JID { return types.NewJID(u, types.DefaultUserServer) }
	cases := []struct {
		name     string
		sender   types.JID
		alt      types.JID
		fromMe   bool
		expected string
	}{
		{"LID with phone alternate", lid(testLIDUnknown), pn("15550000014"), false, "15550000014"},
		{"LID resolved by the LID map", lid(testLIDKnown), types.EmptyJID, false, testPNOfLID},
		{"unknown LID", lid(testLIDUnknown), types.EmptyJID, false, ""},
		{"phone number", pn("15550000015"), types.EmptyJID, false, "15550000015"},
		{"sent by this account", lid(testLIDUnknown), types.EmptyJID, true, testPNOwn},
	}
	for _, c := range cases {
		if got := ident.messagePhone(c.sender, c.alt, c.fromMe); got != c.expected {
			t.Errorf("%s: phone %q, want %q", c.name, got, c.expected)
		}
	}
}

func senderColumns(t *testing.T, store *MessageStore, id string) (sender string, phone sql.NullString) {
	t.Helper()
	if err := store.db.QueryRow("SELECT sender, sender_phone FROM messages WHERE id = ?", id).Scan(&sender, &phone); err != nil {
		t.Fatalf("message %s: %v", id, err)
	}
	return sender, phone
}

func TestSenderPhoneStoredFromLiveMessage(t *testing.T) {
	store := newCaptionStore(t)
	group, _ := types.ParseJID(captionGroup)
	evt := &events.Message{
		Info: types.MessageInfo{
			MessageSource: types.MessageSource{
				Chat:      group,
				Sender:    types.NewJID(testLIDUnknown, types.HiddenUserServer),
				SenderAlt: types.NewJID("15550000016", types.DefaultUserServer),
				IsGroup:   true,
			},
			ID:        "LIVE1",
			Timestamp: time.Date(2026, 11, 5, 8, 29, 20, 0, time.UTC),
		},
		Message: &waE2E.Message{Conversation: proto.String("bug report")},
	}
	handleMessage(nil, store, evt, parseChatAllowlist(captionGroup), nil, waLog.Noop)
	sender, phone := senderColumns(t, store, "LIVE1")
	if sender != testLIDUnknown {
		t.Fatalf("sender changed: %q", sender)
	}
	if phone.String != "15550000016" {
		t.Fatalf("sender_phone %v", phone)
	}
}

func TestSenderPhoneMigrationAndBackfill(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.MkdirAll(storePath(), 0700); err != nil {
		t.Fatal(err)
	}
	// A database created by an earlier version: no sender_phone column
	old, err := sql.Open("sqlite3", sqliteDSN(storePath("messages.db")))
	if err != nil {
		t.Fatal(err)
	}
	_, err = old.Exec(`
		CREATE TABLE chats (jid TEXT PRIMARY KEY, name TEXT, last_message_time TIMESTAMP);
		CREATE TABLE messages (id TEXT, chat_jid TEXT, sender TEXT, content TEXT, timestamp TIMESTAMP,
			is_from_me BOOLEAN, media_type TEXT, filename TEXT, url TEXT, media_key BLOB, file_sha256 BLOB,
			file_enc_sha256 BLOB, file_length INTEGER, PRIMARY KEY (id, chat_jid));`)
	if err != nil {
		t.Fatal(err)
	}
	rows := map[string]struct {
		sender string
		fromMe bool
	}{
		"FULL_LID":     {testLIDKnown + "@lid", false},
		"BARE_LID":     {testLIDKnown, false},
		"FULL_PN":      {"15550000017@s.whatsapp.net", false},
		"BARE_PN":      {testPNMapped, false},
		"UNKNOWN_BARE": {testLIDUnknown, false},
		"UNKNOWN_LID":  {testLIDUnknown + "@lid", false},
		"MINE":         {testPNOwn, true},
	}
	for id, r := range rows {
		if _, err := old.Exec("INSERT INTO messages (id, chat_jid, sender, content, timestamp, is_from_me) VALUES (?, ?, ?, 'x', ?, ?)",
			id, captionGroup, r.sender, time.Now(), r.fromMe); err != nil {
			t.Fatal(err)
		}
	}
	old.Close()

	store, err := NewMessageStore()
	if err != nil {
		t.Fatalf("opening an older database: %v", err)
	}
	defer store.Close()
	// Opening twice must not fail on the existing column
	if again, err := NewMessageStore(); err != nil {
		t.Fatalf("reopening: %v", err)
	} else {
		again.Close()
	}

	n, err := store.BackfillSenderPhones(testSenderIdentity())
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"FULL_LID":     testPNOfLID,
		"BARE_LID":     testPNOfLID,
		"FULL_PN":      "15550000017",
		"BARE_PN":      testPNMapped,
		"UNKNOWN_BARE": "",
		"UNKNOWN_LID":  "",
		"MINE":         testPNOwn,
	}
	if n != 5 {
		t.Errorf("backfilled %d rows, want 5", n)
	}
	for id, phone := range want {
		sender, got := senderColumns(t, store, id)
		if sender != rows[id].sender {
			t.Errorf("%s: sender changed to %q", id, sender)
		}
		if got.String != phone || got.Valid != (phone != "") {
			t.Errorf("%s: sender_phone %v, want %q", id, got, phone)
		}
	}
	// A second run has nothing left to resolve
	if n, err := store.BackfillSenderPhones(testSenderIdentity()); err != nil || n != 0 {
		t.Fatalf("second backfill: %d, %v", n, err)
	}
}

func TestSenderPhoneStoredFromHistorySync(t *testing.T) {
	store := newCaptionStore(t)
	handleHistorySync(nil, store, captionHistorySync(), parseChatAllowlist(captionGroup), waLog.Noop)
	if sender, phone := senderColumns(t, store, "HIST1"); sender != "15550000003@s.whatsapp.net" || phone.String != "15550000003" {
		t.Fatalf("sender %q sender_phone %v", sender, phone)
	}
}
