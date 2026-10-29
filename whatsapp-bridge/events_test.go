package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var (
	testPN      = types.NewJID("15550000001", types.DefaultUserServer)
	testLID     = types.NewJID("123456789012345", types.HiddenUserServer)
	testMsgTime = time.Date(2026, 10, 29, 9, 30, 0, 0, time.UTC)
	testNow     = time.Date(2026, 10, 29, 9, 30, 2, 0, time.UTC)
)

func testResolver(lid types.JID) (types.JID, bool) {
	if lid.User == testLID.User {
		return testPN, true
	}
	return types.EmptyJID, false
}

func newTestMessage(id string, msg *waE2E.Message) *events.Message {
	return &events.Message{
		Info: types.MessageInfo{
			MessageSource: types.MessageSource{Chat: testPN, Sender: testPN},
			ID:            id,
			PushName:      "Test Sender",
			Timestamp:     testMsgTime,
		},
		Message: msg,
	}
}

func newTestSink(t *testing.T, download func(context.Context, whatsmeow.DownloadableMessage) ([]byte, error)) *eventSink {
	t.Helper()
	s := &eventSink{
		inboxDir:   t.TempDir(),
		retention:  30 * 24 * time.Hour,
		httpClient: &http.Client{Timeout: 2 * time.Second},
		retryDelay: time.Millisecond,
		download:   download,
		resolve:    testResolver,
		now:        func() time.Time { return testNow },
		logf:       t.Logf,
	}
	if err := s.prepareInbox(); err != nil {
		t.Fatal(err)
	}
	return s
}

func readSingleEvent(t *testing.T, s *eventSink) (map[string]interface{}, string) {
	t.Helper()
	files, _ := filepath.Glob(filepath.Join(s.inboxDir, "events", "*.json"))
	if len(files) != 1 {
		t.Fatalf("expected 1 event file, got %v", files)
	}
	if tmp, _ := os.ReadDir(filepath.Join(s.inboxDir, "tmp")); len(tmp) != 0 {
		t.Fatalf("temporary files left behind: %v", tmp)
	}
	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	var ev map[string]interface{}
	if err := json.Unmarshal(data, &ev); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	return ev, files[0]
}

func TestEventTextMessage(t *testing.T) {
	s := newTestSink(t, nil)
	msg := newTestMessage("3EB0TEXT", &waE2E.Message{ExtendedTextMessage: &waE2E.ExtendedTextMessage{
		Text:        proto.String("the booking page shows an error"),
		ContextInfo: &waE2E.ContextInfo{StanzaID: proto.String("3EB0QUOTED")},
	}})

	if _, err := s.Process(context.Background(), msg, "Test Chat"); err != nil {
		t.Fatal(err)
	}
	ev, path := readSingleEvent(t, s)

	if filepath.Base(path) != "20261029T093000Z_3EB0TEXT.json" {
		t.Errorf("event file name = %s", filepath.Base(path))
	}
	want := map[string]interface{}{
		"version": 1.0, "type": "message", "id": "3EB0TEXT", "chat_jid": "15550000001@s.whatsapp.net",
		"chat_name": "Test Chat", "sender_jid": "15550000001@s.whatsapp.net", "sender_phone": "15550000001",
		"sender_name": "Test Sender", "is_from_me": false, "timestamp": "2026-10-29T09:30:00Z",
		"received_at": "2026-10-29T09:30:02Z", "text": "the booking page shows an error",
		"quoted_message_id": "3EB0QUOTED", "media": nil, "sender_lid": nil, "chat_lid": nil,
	}
	for k, v := range want {
		if ev[k] != v {
			t.Errorf("%s = %#v, want %#v", k, ev[k], v)
		}
	}
	if runtime.GOOS != "windows" {
		info, _ := os.Stat(path)
		if info.Mode().Perm() != 0640 {
			t.Errorf("event file mode %o, want 640", info.Mode().Perm())
		}
	}
}

func TestEventLIDSenderResolvedToPhone(t *testing.T) {
	msg := newTestMessage("3EB0LID", &waE2E.Message{Conversation: proto.String("hi")})
	msg.Info.Chat, msg.Info.Sender = testLID, testLID

	ev, _ := buildMessageEvent(msg, "", testResolver, testNow)
	if ev.ChatJID != "15550000001@s.whatsapp.net" || ev.ChatLID == nil || *ev.ChatLID != "123456789012345@lid" {
		t.Errorf("chat_jid=%s chat_lid=%v", ev.ChatJID, ev.ChatLID)
	}
	if ev.SenderPhone == nil || *ev.SenderPhone != "15550000001" || ev.SenderLID == nil {
		t.Errorf("sender_phone=%v sender_lid=%v", ev.SenderPhone, ev.SenderLID)
	}

	unknown := newTestMessage("3EB0LID2", &waE2E.Message{Conversation: proto.String("hi")})
	unknown.Info.Sender = types.NewJID("999999999999999", types.HiddenUserServer)
	ev, _ = buildMessageEvent(unknown, "", testResolver, testNow)
	if ev.SenderPhone != nil || ev.SenderLID == nil || *ev.SenderLID != "999999999999999@lid" {
		t.Errorf("unresolved LID: sender_phone=%v sender_lid=%v", ev.SenderPhone, ev.SenderLID)
	}
}

func TestEventMediaWrittenBeforeEvent(t *testing.T) {
	audio := []byte("OggS fake voice note")
	var s *eventSink
	s = newTestSink(t, func(ctx context.Context, msg whatsmeow.DownloadableMessage) ([]byte, error) {
		// The event must not exist yet while the media is being fetched
		if files, _ := filepath.Glob(filepath.Join(s.inboxDir, "events", "*")); len(files) != 0 {
			t.Error("event written before its media")
		}
		return audio, nil
	})
	msg := newTestMessage("3EB0VOICE", &waE2E.Message{AudioMessage: &waE2E.AudioMessage{
		Mimetype: proto.String("audio/ogg; codecs=opus"), PTT: proto.Bool(true), Seconds: proto.Uint32(7),
		FileLength: proto.Uint64(uint64(len(audio))),
	}})

	if _, err := s.Process(context.Background(), msg, ""); err != nil {
		t.Fatal(err)
	}
	ev, _ := readSingleEvent(t, s)
	media := ev["media"].(map[string]interface{})
	sum := sha256.Sum256(audio)
	path, _ := media["path"].(string)
	if media["status"] != "ok" || media["type"] != "audio" || media["voice"] != true || media["duration_s"] != 7.0 ||
		media["sha256"] != hex.EncodeToString(sum[:]) || filepath.Base(path) != "3EB0VOICE.ogg" || !filepath.IsAbs(path) {
		t.Fatalf("unexpected media: %#v", media)
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != string(audio) {
		t.Fatalf("media file content: %q, %v", data, err)
	}
}

func TestEventMediaDownloadFailure(t *testing.T) {
	s := newTestSink(t, func(context.Context, whatsmeow.DownloadableMessage) ([]byte, error) {
		return nil, errors.New("media expired")
	})
	msg := newTestMessage("3EB0DOC", &waE2E.Message{DocumentMessage: &waE2E.DocumentMessage{
		Mimetype: proto.String("application/octet-stream"), FileName: proto.String("../../evil.sh"),
	}})
	if _, err := s.Process(context.Background(), msg, ""); err != nil {
		t.Fatal(err)
	}
	ev, _ := readSingleEvent(t, s)
	media := ev["media"].(map[string]interface{})
	if media["status"] != "failed" || media["error"] != "media expired" || media["path"] != nil {
		t.Fatalf("unexpected media: %#v", media)
	}
	if files, _ := os.ReadDir(filepath.Join(s.inboxDir, "media")); len(files) != 0 {
		t.Fatalf("no media file expected, got %v", files)
	}
}

func TestMediaExtensionNeverUsesSenderPath(t *testing.T) {
	cases := []struct {
		media EventMedia
		want  string
	}{
		{EventMedia{Mime: "image/jpeg"}, "jpg"},
		{EventMedia{Mime: "audio/ogg; codecs=opus"}, "ogg"},
		{EventMedia{Mime: "application/octet-stream", OriginalFilename: "../../evil.sh"}, "sh"},
		{EventMedia{Mime: "application/octet-stream", OriginalFilename: "noextension"}, "bin"},
		{EventMedia{OriginalFilename: "x.verylongextension"}, "bin"},
	}
	for _, tc := range cases {
		if got := mediaExtension(&tc.media); got != tc.want {
			t.Errorf("mediaExtension(%+v) = %q, want %q", tc.media, got, tc.want)
		}
	}
}

func TestWebhookRetriesWithToken(t *testing.T) {
	var calls int32
	received := make(chan []byte, 3)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer hook-secret" || r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("bad headers: %v", r.Header)
		}
		body, _ := io.ReadAll(r.Body)
		received <- body
		if atomic.AddInt32(&calls, 1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	s := newTestSink(t, nil)
	s.webhookURL, s.webhookToken = server.URL, "hook-secret"
	if !s.notify([]byte(`{"id":"x"}`)) {
		t.Fatal("notify should succeed on the third attempt")
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Fatalf("calls = %d, want 3", calls)
	}
	if body := <-received; string(body) != `{"id":"x"}` {
		t.Fatalf("body = %s", body)
	}
}

func TestWebhookURLMustBeLoopback(t *testing.T) {
	for _, u := range []string{"http://127.0.0.1:8090/hook", "http://localhost:9/x", "https://[::1]:443/x"} {
		if err := validateWebhookURL(u); err != nil {
			t.Errorf("%s rejected: %v", u, err)
		}
	}
	for _, u := range []string{"http://example.com/hook", "http://10.0.0.5/hook", "ftp://127.0.0.1/x", "127.0.0.1:8090", "http://localhost.example.com/"} {
		if err := validateWebhookURL(u); err == nil {
			t.Errorf("%s accepted", u)
		}
	}

	t.Setenv("WHATSAPP_INBOX_DIR", "")
	t.Setenv("WHATSAPP_EVENT_WEBHOOK_URL", "http://example.com/hook")
	if _, err := newEventSinkFromEnv(nil); err == nil {
		t.Error("a non-loopback webhook must prevent the bridge from starting")
	}
	t.Setenv("WHATSAPP_EVENT_WEBHOOK_URL", "")
	if s, err := newEventSinkFromEnv(nil); s != nil || err != nil {
		t.Errorf("events must be disabled when nothing is configured: %v %v", s, err)
	}
}

func TestInboxRetentionPurge(t *testing.T) {
	s := newTestSink(t, nil)
	old := filepath.Join(s.inboxDir, "events", "old.json")
	recent := filepath.Join(s.inboxDir, "media", "recent.jpg")
	staleTmp := filepath.Join(s.inboxDir, "tmp", "x.tmp")
	for _, f := range []string{old, recent, staleTmp} {
		if err := os.WriteFile(f, []byte("x"), 0640); err != nil {
			t.Fatal(err)
		}
	}
	os.Chtimes(old, testNow.Add(-31*24*time.Hour), testNow.Add(-31*24*time.Hour))
	os.Chtimes(recent, testNow.Add(-29*24*time.Hour), testNow.Add(-29*24*time.Hour))
	os.Chtimes(staleTmp, testNow.Add(-2*time.Hour), testNow.Add(-2*time.Hour))

	s.Purge()
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Error("event older than the retention was not purged")
	}
	if _, err := os.Stat(staleTmp); !os.IsNotExist(err) {
		t.Error("stale temporary file was not purged")
	}
	if _, err := os.Stat(recent); err != nil {
		t.Error("recent media was purged")
	}
}
