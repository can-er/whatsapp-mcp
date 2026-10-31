//go:build cgo

package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"go.mau.fi/whatsmeow/types"
)

const (
	reactGroup = "120363000000000001@g.us"
	reactChat  = "15550000001@s.whatsapp.net"
	otherChat  = "15550000002@s.whatsapp.net"
)

type fakeTargets struct {
	mu        sync.Mutex
	reactions []string
	replies   []string
}

func (f *fakeTargets) react(chat types.JID, target storedMessage, emoji string) (bool, string, string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.reactions = append(f.reactions, chat.String()+"|"+target.ID+"|"+target.Sender+"|"+emoji)
	return true, "Reaction sent", "3EB0REACT"
}

func (f *fakeTargets) reply(chat types.JID, text string, quoted storedMessage) (bool, string, string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.replies = append(f.replies, chat.String()+"|"+quoted.ID+"|"+text)
	return true, "Reply sent", "3EB0REPLY"
}

func newReactHandler(t *testing.T) (http.Handler, *fakeTargets, string) {
	t.Helper()
	t.Chdir(t.TempDir())
	store, err := NewMessageStore()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	ts := time.Date(2026, 10, 31, 9, 0, 0, 0, time.UTC)
	for _, m := range []struct {
		id, chat, sender, content string
		fromMe                    bool
	}{
		{"3EB0GROUP1", reactGroup, "15550000003", "the preview is broken", false},
		{"3EB0MINE01", reactGroup, "", "fix deployed", true},
		{"3EB0DIRECT", reactChat, "15550000001", "hello", false},
		{"3EB0OTHER1", otherChat, "15550000002", "not allowed", false},
	} {
		store.StoreChat(m.chat, "Chat", ts)
		if err := store.StoreMessage(m.id, m.chat, m.sender, m.content, ts, m.fromMe, "", "", "", nil, nil, nil, 0); err != nil {
			t.Fatal(err)
		}
	}

	logPath := storePath("send_audit.jsonl")
	log, err := openSendLog(logPath)
	if err != nil {
		t.Fatal(err)
	}
	fake := &fakeTargets{}
	opts := apiOptions{sendLog: log, react: fake.react, reply: fake.reply,
		send: func(recipient, message, mediaPath string) (bool, string, string) {
			t.Error("a reply must not go through the plain send")
			return false, "", ""
		}}
	allowed := parseChatAllowlist("15550000001," + reactGroup)
	tokens := apiTokens{bridge: testToken, send: testSendToken}
	return requireAPIAuth(tokens, newAPIMux(nil, store, allowed, opts)), fake, logPath
}

func postJSON(handler http.Handler, path, token, body string) (int, SendMessageResponse) {
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, newTestRequest(http.MethodPost, path, body, map[string]string{
		"Authorization": "Bearer " + token, "Content-Type": "application/json",
	}))
	var resp SendMessageResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	return rec.Code, resp
}

func TestReactAPI(t *testing.T) {
	handler, fake, logPath := newReactHandler(t)

	body := `{"chat_jid":"` + reactGroup + `","message_id":"3EB0GROUP1","emoji":"🤖","client_msg_id":"react-1"}`
	code, resp := postJSON(handler, "/api/react", testSendToken, body)
	if code != http.StatusOK || !resp.Success || resp.MessageID != "3EB0REACT" || resp.Duplicate {
		t.Fatalf("react: %d %+v", code, resp)
	}
	if code, resp = postJSON(handler, "/api/react", testSendToken, body); code != http.StatusOK || !resp.Duplicate {
		t.Errorf("retry: %d %+v, want duplicate", code, resp)
	}
	if len(fake.reactions) != 1 || fake.reactions[0] != reactGroup+"|3EB0GROUP1|15550000003|🤖" {
		t.Errorf("reactions sent: %v", fake.reactions)
	}

	for name, tc := range map[string]struct {
		token, body string
		want        int
	}{
		"bridge token":    {testToken, body, http.StatusUnauthorized},
		"not allowed":     {testSendToken, `{"chat_jid":"` + otherChat + `","message_id":"3EB0OTHER1","emoji":"👍"}`, http.StatusForbidden},
		"unknown message": {testSendToken, `{"chat_jid":"` + reactGroup + `","message_id":"3EB0NOPE00","emoji":"👍"}`, http.StatusNotFound},
		"other chat id":   {testSendToken, `{"chat_jid":"` + reactGroup + `","message_id":"3EB0DIRECT","emoji":"👍"}`, http.StatusNotFound},
		"empty emoji":     {testSendToken, `{"chat_jid":"` + reactGroup + `","message_id":"3EB0GROUP1","emoji":""}`, http.StatusBadRequest},
		"long emoji":      {testSendToken, `{"chat_jid":"` + reactGroup + `","message_id":"3EB0GROUP1","emoji":"` + strings.Repeat("x", 65) + `"}`, http.StatusBadRequest},
		"missing id":      {testSendToken, `{"chat_jid":"` + reactGroup + `","emoji":"👍"}`, http.StatusBadRequest},
	} {
		if code, _ := postJSON(handler, "/api/react", tc.token, tc.body); code != tc.want {
			t.Errorf("%s: status %d, want %d", name, code, tc.want)
		}
	}
	if len(fake.reactions) != 1 {
		t.Errorf("rejected requests must not send: %v", fake.reactions)
	}

	// Audit: kind and target, never the emoji itself
	data, _ := os.ReadFile(logPath)
	if strings.Contains(string(data), "🤖") {
		t.Error("the audit log must not contain the emoji")
	}
	var rec map[string]interface{}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	scanner.Scan()
	json.Unmarshal(scanner.Bytes(), &rec)
	if rec["kind"] != "reaction" || rec["target_message_id"] != "3EB0GROUP1" || rec["client_msg_id"] != "react-1" {
		t.Errorf("unexpected audit record %v", rec)
	}
}

func TestReplyAPI(t *testing.T) {
	handler, fake, logPath := newReactHandler(t)

	body := `{"recipient":"` + reactGroup + `","message":"Preview ready","reply_to_message_id":"3EB0GROUP1","client_msg_id":"reply-1"}`
	code, resp := postJSON(handler, "/api/send", testSendToken, body)
	if code != http.StatusOK || !resp.Success || resp.MessageID != "3EB0REPLY" {
		t.Fatalf("reply: %d %+v", code, resp)
	}
	if _, resp = postJSON(handler, "/api/send", testSendToken, body); !resp.Duplicate {
		t.Error("a retried reply must be a duplicate")
	}
	// A direct chat given as a phone number finds the message stored under its JID
	if code, _ := postJSON(handler, "/api/send", testSendToken, `{"recipient":"+15550000001","message":"ok","reply_to_message_id":"3EB0DIRECT"}`); code != http.StatusOK {
		t.Errorf("reply in a direct chat: status %d", code)
	}
	want := []string{reactGroup + "|3EB0GROUP1|Preview ready", reactChat + "|3EB0DIRECT|ok"}
	if strings.Join(fake.replies, ",") != strings.Join(want, ",") {
		t.Errorf("replies sent: %v", fake.replies)
	}

	for name, tc := range map[string]struct {
		body string
		want int
	}{
		"other chat":      {`{"recipient":"` + reactGroup + `","message":"x","reply_to_message_id":"3EB0DIRECT"}`, http.StatusBadRequest},
		"unknown message": {`{"recipient":"` + reactGroup + `","message":"x","reply_to_message_id":"3EB0NOPE00"}`, http.StatusBadRequest},
		"with media":      {`{"recipient":"` + reactGroup + `","message":"x","media_path":"a.jpg","reply_to_message_id":"3EB0GROUP1"}`, http.StatusBadRequest},
		"not allowed":     {`{"recipient":"` + otherChat + `","message":"x","reply_to_message_id":"3EB0OTHER1"}`, http.StatusForbidden},
	} {
		if code, _ := postJSON(handler, "/api/send", testSendToken, tc.body); code != tc.want {
			t.Errorf("%s: status %d, want %d", name, code, tc.want)
		}
	}
	if len(fake.replies) != 2 {
		t.Errorf("rejected replies must not send: %v", fake.replies)
	}
	data, _ := os.ReadFile(logPath)
	if strings.Contains(string(data), "Preview ready") || !strings.Contains(string(data), `"kind":"reply"`) {
		t.Errorf("unexpected audit log: %s", data)
	}
}

func TestReplyAndReactionMessages(t *testing.T) {
	group, _ := types.ParseJID(reactGroup)
	direct, _ := types.ParseJID(reactChat)
	now := time.Date(2026, 10, 31, 9, 0, 0, 0, time.UTC)

	// Reaction to someone else's message in a group: participant set
	r := reactionMessage(group, storedMessage{ID: "3EB0GROUP1", Sender: "15550000003"}, "🤖", nil, now).GetReactionMessage()
	if r.GetText() != "🤖" || r.GetKey().GetID() != "3EB0GROUP1" || r.GetKey().GetRemoteJID() != reactGroup ||
		r.GetKey().GetFromMe() || r.GetKey().GetParticipant() != "15550000003@s.whatsapp.net" || r.GetSenderTimestampMS() != now.UnixMilli() {
		t.Errorf("unexpected group reaction %v", r)
	}
	// Own message: from_me, no participant; direct chat: no participant
	if k := reactionMessage(group, storedMessage{ID: "3EB0MINE01", IsFromMe: true}, "👍", nil, now).GetReactionMessage().GetKey(); !k.GetFromMe() || k.Participant != nil {
		t.Errorf("unexpected key for an own message %v", k)
	}
	if k := reactionMessage(direct, storedMessage{ID: "3EB0DIRECT", Sender: "15550000001"}, "👍", nil, now).GetReactionMessage().GetKey(); k.GetFromMe() || k.Participant != nil {
		t.Errorf("unexpected key in a direct chat %v", k)
	}

	// Quoted reply
	quoted := storedMessage{ID: "3EB0GROUP1", Sender: "15550000003", Content: "the preview is broken"}
	m := quotedReplyMessage("Preview ready", quoted, types.NewJID("15550000003", types.DefaultUserServer)).GetExtendedTextMessage()
	ci := m.GetContextInfo()
	if m.GetText() != "Preview ready" || ci.GetStanzaID() != "3EB0GROUP1" || ci.GetParticipant() != "15550000003@s.whatsapp.net" ||
		ci.GetQuotedMessage().GetConversation() != "the preview is broken" {
		t.Errorf("unexpected reply %v", m)
	}
	// Author of a quoted message stored as sent by another participant
	if p := quotedParticipant(nil, group, quoted); p.String() != "15550000003@s.whatsapp.net" {
		t.Errorf("quoted participant %s", p)
	}
	if !validEmoji("👍🏽") || !validEmoji("👨‍👩‍👧") || validEmoji("a b") || validEmoji("\n") {
		t.Error("emoji validation is wrong")
	}
}
