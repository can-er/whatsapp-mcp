package main

import (
	"context"
	"os"
	"strings"

	"go.mau.fi/whatsmeow"
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

// AllowsChat is Allows for a chat that may be addressed by its LID (hidden user ID): the phone
// number the LID maps to in the device store is checked as well, so listing phone numbers works
// for LID-addressed chats too.
func (a *chatAllowlist) AllowsChat(client *whatsmeow.Client, jid types.JID) bool {
	if a.Allows(jid) {
		return true
	}
	if jid.Server != types.HiddenUserServer || client == nil || client.Store == nil || client.Store.LIDs == nil {
		return false
	}
	pn, err := client.Store.LIDs.GetPNForLID(context.Background(), jid.ToNonAD())
	return err == nil && !pn.IsEmpty() && a.Allows(pn)
}

// AllowsChatString is AllowsChat for a JID or phone number given as a string
func (a *chatAllowlist) AllowsChatString(client *whatsmeow.Client, chat string) bool {
	if a == nil {
		return true
	}
	normalized, ok := normalizeChatJID(chat)
	if !ok {
		return false
	}
	jid, err := types.ParseJID(normalized)
	return err == nil && a.AllowsChat(client, jid)
}
