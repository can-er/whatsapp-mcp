package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.mau.fi/whatsmeow/types"
)

func TestChatAllowlistParsing(t *testing.T) {
	if parseChatAllowlist("") != nil || parseChatAllowlist(" , ,") != nil {
		t.Fatal("an empty list must allow every chat (nil allowlist)")
	}

	a := parseChatAllowlist(" +15550000001 , 15550000002@s.whatsapp.net,120363000000000001@g.us, not a jid@@ ")
	if len(a.jids) != 3 {
		t.Fatalf("parsed %d entries, want 3: %v", len(a.jids), a.jids)
	}

	allowed := []types.JID{
		types.NewJID("15550000001", types.DefaultUserServer),
		types.NewJID("15550000002", types.DefaultUserServer),
		types.NewJID("120363000000000001", types.GroupServer),
		{User: "15550000001", Server: types.DefaultUserServer, Device: 3}, // message from a linked device
	}
	for _, jid := range allowed {
		if !a.Allows(jid) {
			t.Errorf("Allows(%s) = false, want true", jid)
		}
	}

	denied := []types.JID{
		types.NewJID("32400000000", types.DefaultUserServer),
		types.NewJID("120363000000000002", types.GroupServer),
		types.NewJID("15550000001", types.GroupServer),
	}
	for _, jid := range denied {
		if a.Allows(jid) {
			t.Errorf("Allows(%s) = true, want false", jid)
		}
	}

	for chat, want := range map[string]bool{
		"15550000001":                true,
		"+15550000001":               true,
		"15550000001@s.whatsapp.net": true,
		"32400000000":                false,
		"":                           false,
		"../../etc":                  false,
	} {
		if got := a.AllowsChatString(nil, chat); got != want {
			t.Errorf("AllowsChatString(%q) = %v, want %v", chat, got, want)
		}
	}

	var everyone *chatAllowlist
	if !everyone.Allows(types.NewJID("1", types.DefaultUserServer)) || !everyone.AllowsChatString(nil, "1") {
		t.Fatal("nil allowlist must allow every chat")
	}
}

// Chats outside the allowlist cannot be reached through the API; the request is rejected
// before the (nil) WhatsApp client is used.
func TestChatAllowlistAPI(t *testing.T) {
	handler := requireAPIAuth(testToken, newAPIMux(nil, nil, parseChatAllowlist("15550000001")))

	requests := map[string]string{
		"/api/send":     `{"recipient":"32400000000","message":"hi"}`,
		"/api/download": `{"message_id":"ABC","chat_jid":"32400000000@s.whatsapp.net"}`,
	}
	for path, body := range requests {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, newTestRequest(http.MethodPost, path, body, validHeaders()))
		if rec.Code != http.StatusForbidden {
			t.Errorf("%s: status = %d, want 403", path, rec.Code)
		}
	}
}
