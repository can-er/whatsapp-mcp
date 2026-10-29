//go:build cgo

package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func TestStoreDirDefault(t *testing.T) {
	t.Setenv("WHATSAPP_STORE_DIR", "")
	if storePath() != "store" || storePath("messages.db") != filepath.Join("store", "messages.db") {
		t.Fatalf("default store path changed: %q", storePath("messages.db"))
	}
}

// Every file of the bridge goes to WHATSAPP_STORE_DIR, nothing to ./store
func TestStoreDirFromEnv(t *testing.T) {
	work := t.TempDir()
	t.Chdir(work)
	data := filepath.Join(t.TempDir(), "state dir", "ks-wa-bridge")
	t.Setenv("WHATSAPP_STORE_DIR", data)
	t.Setenv("WHATSAPP_BRIDGE_TOKEN", "")

	store, err := NewMessageStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := ensurePrivateFile(storePath("whatsapp.db")); err != nil {
		t.Fatal(err)
	}
	container, err := sqlstore.New(context.Background(), "sqlite3", sqliteDSN(storePath("whatsapp.db")), waLog.Noop)
	if err != nil {
		t.Fatalf("session store in WHATSAPP_STORE_DIR: %v", err)
	}
	container.Close()
	if _, err := loadOrCreateBridgeToken(storePath("bridge_token")); err != nil {
		t.Fatal(err)
	}
	if _, err := openSendLog(storePath("send_audit.jsonl")); err != nil {
		t.Fatal(err)
	}
	chatDir, mediaPath, err := localMediaPath("15550000001@s.whatsapp.net", "image", "3EB0IMG", "photo.jpg")
	if err != nil || filepath.Dir(chatDir) != data || filepath.Dir(mediaPath) != chatDir {
		t.Fatalf("media path %q not under %q: %v", mediaPath, data, err)
	}

	for _, f := range []string{"messages.db", "whatsapp.db", "bridge_token", "send_audit.jsonl"} {
		info, err := os.Stat(filepath.Join(data, f))
		if err != nil {
			t.Fatalf("%s not created in WHATSAPP_STORE_DIR: %v", f, err)
		}
		if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
			t.Errorf("%s mode %o, want 600", f, info.Mode().Perm())
		}
	}
	if info, _ := os.Stat(data); runtime.GOOS != "windows" && info.Mode().Perm() != 0700 {
		t.Errorf("store dir mode %o, want 700", info.Mode().Perm())
	}
	if _, err := os.Stat(filepath.Join(work, "store")); !os.IsNotExist(err) {
		t.Fatal("nothing must be written to ./store when WHATSAPP_STORE_DIR is set")
	}
}
