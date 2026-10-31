//go:build cgo

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPurgeDisallowedChats(t *testing.T) {
	t.Chdir(t.TempDir())
	store, err := NewMessageStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	allowedChat := "15550000001@s.whatsapp.net"
	otherChat := "15550000002@s.whatsapp.net"
	group := "120363000000000002@g.us"
	broken := "not-a-jid"
	ts := time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC)
	for i, chat := range []string{allowedChat, otherChat, group, broken} {
		store.StoreChat(chat, "Chat", ts)
		for j := 0; j < 3; j++ {
			content := "kept message"
			if chat != allowedChat {
				content = "PRIVATE-SECRET-TEXT"
			}
			id := string(rune('A'+i)) + string(rune('0'+j))
			if err := store.StoreMessage(id, chat, "15550000009", content, ts, false, "", "", "", nil, nil, nil, 0); err != nil {
				t.Fatal(err)
			}
		}
	}
	// Downloaded media of an allowed and of a disallowed chat, plus store files that must stay
	for _, chat := range []string{allowedChat, otherChat} {
		os.MkdirAll(storePath(chat), 0700)
		os.WriteFile(storePath(chat, "photo.jpg"), []byte("x"), 0600)
	}
	os.WriteFile(storePath("bridge_token"), []byte("t"), 0600)

	allowlist := parseChatAllowlist("15550000001")
	result, err := store.PurgeDisallowedChats(allowlist.Allows, storePath())
	if err != nil {
		t.Fatal(err)
	}
	if result.Chats != 3 || result.Messages != 9 || result.MediaDirs != 1 {
		t.Errorf("unexpected result %+v", result)
	}

	var chats, messages int
	store.db.QueryRow("SELECT COUNT(*) FROM chats").Scan(&chats)
	store.db.QueryRow("SELECT COUNT(*) FROM messages WHERE chat_jid = ?", allowedChat).Scan(&messages)
	if chats != 1 || messages != 3 {
		t.Errorf("allowed chat altered: %d chats, %d messages", chats, messages)
	}
	if _, err := os.Stat(storePath(allowedChat, "photo.jpg")); err != nil {
		t.Error("media of the allowed chat removed")
	}
	if _, err := os.Stat(storePath(otherChat)); !os.IsNotExist(err) {
		t.Error("media of the disallowed chat kept")
	}
	if _, err := os.Stat(storePath("bridge_token")); err != nil {
		t.Error("store files must not be removed")
	}

	// The deleted text must not remain in the database files (free pages, WAL)
	files, _ := filepath.Glob(storePath("messages.db*"))
	for _, f := range files {
		data, _ := os.ReadFile(f)
		if bytes.Contains(data, []byte("PRIVATE-SECRET-TEXT")) {
			t.Errorf("deleted text still present in %s", f)
		}
	}
}
