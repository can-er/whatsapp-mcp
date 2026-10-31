package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.mau.fi/whatsmeow/types"
)

// purgeResult counts what PurgeDisallowedChats deleted
type purgeResult struct {
	Chats     int64
	Messages  int64
	MediaDirs int
}

// PurgeDisallowedChats deletes the chats, messages and downloaded media of every chat that allows
// rejects (WHATSAPP_PURGE_DISALLOWED_CHATS), then compacts the database so the deleted content does
// not stay in free pages. Chat JIDs that cannot be parsed are deleted too.
func (store *MessageStore) PurgeDisallowedChats(allows func(types.JID) bool, mediaDir string) (purgeResult, error) {
	var result purgeResult
	rows, err := store.db.Query("SELECT jid FROM chats UNION SELECT DISTINCT chat_jid FROM messages")
	if err != nil {
		return result, err
	}
	var purge []string
	for rows.Next() {
		var chat string
		if err := rows.Scan(&chat); err != nil {
			rows.Close()
			return result, err
		}
		if jid, err := types.ParseJID(chat); err != nil || jid.Server == "" || !allows(jid) {
			purge = append(purge, chat)
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return result, err
	}

	tx, err := store.db.Begin()
	if err != nil {
		return result, err
	}
	for _, chat := range purge {
		res, err := tx.Exec("DELETE FROM messages WHERE chat_jid = ?", chat)
		if err != nil {
			tx.Rollback()
			return result, err
		}
		n, _ := res.RowsAffected()
		result.Messages += n
		if res, err = tx.Exec("DELETE FROM chats WHERE jid = ?", chat); err != nil {
			tx.Rollback()
			return result, err
		}
		n, _ = res.RowsAffected()
		result.Chats += n
	}
	if err := tx.Commit(); err != nil {
		return result, err
	}
	if _, err := store.db.Exec("VACUUM"); err != nil {
		return result, fmt.Errorf("failed to compact the database: %v", err)
	}
	if _, err := store.db.Exec("PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
		return result, err
	}

	// Media are downloaded in <store>/<chat JID>/; other entries of the store are files
	entries, err := os.ReadDir(mediaDir)
	if err != nil {
		return result, err
	}
	for _, entry := range entries {
		if !entry.IsDir() || !strings.Contains(entry.Name(), "@") {
			continue
		}
		if jid, err := types.ParseJID(entry.Name()); err == nil && jid.Server != "" && allows(jid) {
			continue
		}
		if err := os.RemoveAll(filepath.Join(mediaDir, entry.Name())); err != nil {
			return result, err
		}
		result.MediaDirs++
	}
	return result, nil
}
