package main

import (
	"os"
	"strings"

	"go.mau.fi/whatsmeow/types"
)

// chatAllowlist restricts which chats are stored and reachable through the API.
// A nil allowlist allows every chat (default behaviour).
type chatAllowlist struct {
	jids map[string]bool
}

// loadChatAllowlist reads WHATSAPP_ALLOWED_CHATS
func loadChatAllowlist() *chatAllowlist {
	return parseChatAllowlist(os.Getenv("WHATSAPP_ALLOWED_CHATS"))
}

// parseChatAllowlist parses a comma-separated list of JIDs or phone numbers.
// It returns nil (every chat allowed) when the list is empty.
func parseChatAllowlist(raw string) *chatAllowlist {
	jids := make(map[string]bool)
	for _, entry := range strings.Split(raw, ",") {
		if jid, ok := normalizeChatJID(entry); ok {
			jids[jid] = true
		}
	}
	if len(jids) == 0 {
		return nil
	}
	return &chatAllowlist{jids: jids}
}

// normalizeChatJID turns a JID or a phone number (with or without +, spaces) into a chat JID
// without device part, e.g. "+1 555 000 0005" -> "15550000005@s.whatsapp.net"
func normalizeChatJID(entry string) (string, bool) {
	entry = strings.ReplaceAll(strings.TrimSpace(entry), " ", "")
	if entry == "" {
		return "", false
	}
	if !strings.Contains(entry, "@") {
		entry = strings.TrimPrefix(entry, "+") + "@" + types.DefaultUserServer
	}
	jid, err := types.ParseJID(entry)
	if err != nil || jid.User == "" || jid.Server == "" {
		return "", false
	}
	return jid.ToNonAD().String(), true
}

// Allows reports whether the chat may be stored and used through the API
func (a *chatAllowlist) Allows(jid types.JID) bool {
	if a == nil {
		return true
	}
	return a.jids[jid.ToNonAD().String()]
}

// AllowsString is Allows for a JID or phone number given as a string
func (a *chatAllowlist) AllowsString(chat string) bool {
	if a == nil {
		return true
	}
	jid, ok := normalizeChatJID(chat)
	return ok && a.jids[jid]
}
