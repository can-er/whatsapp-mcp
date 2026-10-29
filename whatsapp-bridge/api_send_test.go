package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

const testSendToken = "send-token-0123456789abcdef0123456789"

// fakeSender records calls instead of sending WhatsApp messages
type fakeSender struct {
	mu    sync.Mutex
	calls int
	fail  bool
}

func (f *fakeSender) send(recipient, message, mediaPath string) (bool, string, string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	if f.fail {
		return false, "Error sending message: offline", ""
	}
	return true, "Message sent to " + recipient, "3EB0SENT" + strings.Repeat("X", f.calls)
}

func postSend(t *testing.T, handler http.Handler, token, body string) (int, SendMessageResponse) {
	t.Helper()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, newTestRequest(http.MethodPost, "/api/send", body, map[string]string{
		"Authorization": "Bearer " + token, "Content-Type": "application/json",
	}))
	var resp SendMessageResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	return rec.Code, resp
}

func newSendHandler(t *testing.T, sender *fakeSender) (http.Handler, string) {
	t.Helper()
	logPath := filepath.Join(t.TempDir(), "send_audit.jsonl")
	log, err := openSendLog(logPath)
	if err != nil {
		t.Fatal(err)
	}
	opts := apiOptions{send: sender.send, sendLog: log, health: func() healthStatus {
		last := time.Date(2026, 10, 29, 9, 30, 0, 0, time.UTC)
		return healthStatus{Connected: true, LoggedIn: true, LastEventAt: &last}
	}}
	tokens := apiTokens{bridge: testToken, send: testSendToken}
	return requireAPIAuth(tokens, newAPIMux(nil, nil, parseChatAllowlist("15550000001"), opts)), logPath
}

func TestSendTokenSeparation(t *testing.T) {
	sender := &fakeSender{}
	handler, _ := newSendHandler(t, sender)
	body := `{"recipient":"15550000001","message":"preview ready"}`

	if code, _ := postSend(t, handler, testToken, body); code != http.StatusUnauthorized {
		t.Errorf("bridge token on /api/send: status %d, want 401", code)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, newTestRequest(http.MethodPost, "/api//send", body, map[string]string{
		"Authorization": "Bearer " + testToken, "Content-Type": "application/json",
	}))
	if rec.Code == http.StatusOK {
		t.Errorf("bridge token on an unclean /api/send path must not send")
	}
	if code, resp := postSend(t, handler, testSendToken, body); code != http.StatusOK || !resp.Success || resp.MessageID == "" {
		t.Errorf("send token on /api/send: status %d, %+v", code, resp)
	}

	// The send token only opens /api/send
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, newTestRequest(http.MethodGet, "/api/health", "", map[string]string{"Authorization": "Bearer " + testSendToken}))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("send token on /api/health: status %d, want 401", rec.Code)
	}
	if sender.calls != 1 {
		t.Errorf("sender called %d times, want 1", sender.calls)
	}
}

func TestSendIdempotency(t *testing.T) {
	sender := &fakeSender{}
	handler, _ := newSendHandler(t, sender)
	body := `{"recipient":"15550000001","message":"preview ready","client_msg_id":"ks-12-preview"}`

	_, first := postSend(t, handler, testSendToken, body)
	code, second := postSend(t, handler, testSendToken, body)
	if sender.calls != 1 {
		t.Fatalf("sender called %d times for the same client_msg_id, want 1", sender.calls)
	}
	if code != http.StatusOK || !second.Duplicate || second.MessageID != first.MessageID || first.Duplicate {
		t.Fatalf("first %+v, second %+v (status %d)", first, second, code)
	}

	postSend(t, handler, testSendToken, `{"recipient":"15550000001","message":"other","client_msg_id":"ks-13-preview"}`)
	postSend(t, handler, testSendToken, `{"recipient":"15550000001","message":"no key"}`)
	postSend(t, handler, testSendToken, `{"recipient":"15550000001","message":"no key"}`)
	if sender.calls != 4 {
		t.Fatalf("sender called %d times, want 4 (new key + 2 requests without key)", sender.calls)
	}

	// Concurrent retries of the same request send only once
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			postSend(t, handler, testSendToken, `{"recipient":"15550000001","message":"x","client_msg_id":"ks-14"}`)
		}()
	}
	wg.Wait()
	if sender.calls != 5 {
		t.Fatalf("sender called %d times after concurrent retries, want 5", sender.calls)
	}
}

func TestSendFailureIsNotCached(t *testing.T) {
	sender := &fakeSender{fail: true}
	handler, _ := newSendHandler(t, sender)
	body := `{"recipient":"15550000001","message":"preview ready","client_msg_id":"ks-15"}`

	if code, resp := postSend(t, handler, testSendToken, body); code != http.StatusInternalServerError || resp.Success {
		t.Fatalf("failed send: status %d, %+v", code, resp)
	}
	sender.fail = false
	if code, resp := postSend(t, handler, testSendToken, body); code != http.StatusOK || resp.Duplicate {
		t.Fatalf("retry after failure: status %d, %+v", code, resp)
	}
	if sender.calls != 2 {
		t.Fatalf("sender called %d times, want 2", sender.calls)
	}
}

func TestSendAuditLog(t *testing.T) {
	sender := &fakeSender{}
	handler, logPath := newSendHandler(t, sender)
	text := "preview ready: confidential details"
	postSend(t, handler, testSendToken, `{"recipient":"15550000001","message":"`+text+`","client_msg_id":"ks-16"}`)

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "confidential") {
		t.Fatal("the audit log must not contain the message text")
	}
	var rec sendRecord
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(data))), &rec); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(text))
	if rec.TextSHA256 != hex.EncodeToString(sum[:]) || rec.Recipient != "15550000001" || rec.ClientMsgID != "ks-16" ||
		!rec.Success || rec.MessageID == "" || rec.Time.IsZero() {
		t.Fatalf("unexpected audit record: %+v", rec)
	}
	if runtime.GOOS != "windows" {
		if info, _ := os.Stat(logPath); info.Mode().Perm() != 0600 {
			t.Errorf("audit log mode %o, want 600", info.Mode().Perm())
		}
	}

	// Idempotency survives a restart: the log is replayed on open
	reopened, err := openSendLog(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, dup := reopened.Send("ks-16", "15550000001", text, "", sender.send); !dup || sender.calls != 1 {
		t.Fatalf("after restart: duplicate=%v, calls=%d", dup, sender.calls)
	}
}

func TestSendRespectsAllowlist(t *testing.T) {
	sender := &fakeSender{}
	handler, _ := newSendHandler(t, sender)
	if code, _ := postSend(t, handler, testSendToken, `{"recipient":"15550000009","message":"x"}`); code != http.StatusForbidden || sender.calls != 0 {
		t.Fatalf("recipient outside the allowlist: status %d, calls %d", code, sender.calls)
	}
}

func TestHealth(t *testing.T) {
	handler, _ := newSendHandler(t, &fakeSender{})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, newTestRequest(http.MethodGet, "/api/health", "", map[string]string{"Authorization": "Bearer " + testToken}))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var got map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &got)
	if got["connected"] != true || got["logged_in"] != true || got["last_event_at"] != "2026-10-29T09:30:00Z" {
		t.Fatalf("unexpected health: %v", got)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, newTestRequest(http.MethodGet, "/api/health", "", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("health without token: status %d, want 401", rec.Code)
	}
}
