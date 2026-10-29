package main

import (
	"os"
	"path/filepath"
	"strings"
)

// storePath returns a path inside the data directory: WHATSAPP_STORE_DIR, or "store" relative to
// the working directory (default). It holds the session keys, messages, API token, send audit and media.
func storePath(elem ...string) string {
	dir := strings.TrimSpace(os.Getenv("WHATSAPP_STORE_DIR"))
	if dir == "" {
		dir = "store"
	}
	return filepath.Join(append([]string{dir}, elem...)...)
}

// sqliteDSN builds a SQLite URI for a database file
func sqliteDSN(path string) string {
	return "file:" + filepath.ToSlash(path) + "?_foreign_keys=on"
}

// ensurePrivateDir creates dir (and its parents) if needed and restricts it to the owner.
// The store directory holds the WhatsApp session keys and the message history.
func ensurePrivateDir(dir string) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	// MkdirAll keeps the mode of an existing directory, so tighten it explicitly
	return os.Chmod(dir, 0700)
}

// ensurePrivateFile creates path if needed and restricts it to the owner, keeping its content.
// Creating SQLite databases this way before opening them also gives their -wal/-shm/-journal
// files the same mode, as SQLite copies the database file permissions.
func ensurePrivateFile(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}
