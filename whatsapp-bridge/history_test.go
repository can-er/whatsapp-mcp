//go:build cgo

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const historyChat = "15550000001@s.whatsapp.net"

func newHistoryHandler(t *testing.T) (http.Handler, []string) {
	t.Helper()
	t.Chdir(t.TempDir())
	store, err := NewMessageStore()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })

	base := time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC)
	other := "15550000002@s.whatsapp.net"
	store.StoreChat(historyChat, "Test Chat", base)
	store.StoreChat(other, "Other", base)

	var ids []string
	for i := 0; i < 120; i++ {
		ts := base.Add(time.Duration(i) * time.Minute)
		if i >= 60 && i < 90 {
			ts = base.Add(60 * time.Minute) // a burst: 30 messages in the same second
		}
		if i%2 == 1 {
			ts = ts.In(time.FixedZone("CET", 3600)) // stored with another UTC offset
		}
		id := fmt.Sprintf("M%03d", i)
		ids = append(ids, id)
		if err := store.StoreMessage(id, historyChat, "15550000001", "text "+id, ts, i%3 == 0, "", "", "", nil, nil, nil, 0); err != nil {
			t.Fatal(err)
		}
	}
	// Media: one downloaded, one not
	store.StoreMessage("IMG1", historyChat, "15550000001", "", base.Add(-time.Hour), false, "image", "image_1.jpg", "", nil, nil, nil, 0)
	store.StoreMessage("DOC1", historyChat, "15550000001", "", base.Add(-2*time.Hour), false, "document", "../../report.pdf", "", nil, nil, nil, 0)
	if err := os.MkdirAll(filepath.Join("store", historyChat), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("store", historyChat, "image_1.jpg"), []byte("jpeg"), 0600); err != nil {
		t.Fatal(err)
	}
	store.StoreMessage("HIDDEN", other, "15550000002", "secret", base, false, "", "", "", nil, nil, nil, 0)

	handler := requireAPIAuth(apiTokens{bridge: testToken}, newAPIMux(nil, store, parseChatAllowlist("15550000001"), apiOptions{}))
	return handler, append(ids, "IMG1", "DOC1")
}

func getHistory(t *testing.T, handler http.Handler, query url.Values, token string) (int, historyResponse) {
	t.Helper()
	headers := map[string]string{}
	if token != "" {
		headers["Authorization"] = "Bearer " + token
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, newTestRequest(http.MethodGet, "/api/history?"+query.Encode(), "", headers))
	var resp historyResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	return rec.Code, resp
}

func TestHistoryPagination(t *testing.T) {
	handler, all := newHistoryHandler(t)

	var got []*MessageEvent
	before := ""
	for pages := 0; ; pages++ {
		if pages > 10 {
			t.Fatal("pagination does not end")
		}
		q := url.Values{"chat_jid": {"15550000001"}, "limit": {"50"}}
		if before != "" {
			q.Set("before", before)
		}
		code, resp := getHistory(t, handler, q, testToken)
		if code != http.StatusOK {
			t.Fatalf("status %d", code)
		}
		got = append(got, resp.Messages...)
		if resp.NextBefore == nil {
			break
		}
		before = *resp.NextBefore
	}

	if len(got) != len(all) {
		t.Fatalf("got %d messages, want %d", len(got), len(all))
	}
	seen := map[string]bool{}
	for i, m := range got {
		if seen[m.ID] {
			t.Fatalf("duplicate %s", m.ID)
		}
		seen[m.ID] = true
		if m.ChatJID != historyChat || m.Type != "message" || m.ReceivedAt != nil {
			t.Fatalf("unexpected event: %+v", m)
		}
		if i > 0 {
			prev := got[i-1]
			if m.Timestamp.After(prev.Timestamp) || (m.Timestamp.Equal(prev.Timestamp) && m.ID > prev.ID) {
				t.Fatalf("order broken at %d: %s after %s", i, m.ID, prev.ID)
			}
		}
	}
	byID := map[string]*MessageEvent{}
	for _, m := range got {
		byID[m.ID] = m
	}
	if m := byID["M007"]; m.Text != "text M007" || m.SenderPhone == nil || *m.SenderPhone != "15550000001" || m.SenderJID != historyChat {
		t.Errorf("text message: %+v", m)
	}
	img := byID["IMG1"].Media
	if img == nil || img.Status != "ok" || img.Path == nil || filepath.Base(*img.Path) != "image_1.jpg" || img.Size != 4 || img.SHA256 == "" {
		t.Errorf("downloaded media: %+v", img)
	}
	doc := byID["DOC1"].Media
	if doc == nil || doc.Status != "not_downloaded" || doc.Path != nil || doc.OriginalFilename != "../../report.pdf" {
		t.Errorf("not downloaded media: %+v", doc)
	}

	// An RFC 3339 time is accepted as cursor
	q := url.Values{"chat_jid": {historyChat}, "before": {"2026-10-01T08:00:30Z"}}
	if code, resp := getHistory(t, handler, q, testToken); code != http.StatusOK || len(resp.Messages) != 3 || resp.Messages[0].ID != "M000" {
		t.Errorf("RFC 3339 cursor: status %d, %d messages", code, len(resp.Messages))
	}
}

func TestHistoryRejections(t *testing.T) {
	handler, _ := newHistoryHandler(t)
	cases := []struct {
		name  string
		query url.Values
		token string
		want  int
	}{
		{"chat outside allowlist", url.Values{"chat_jid": {"15550000002"}}, testToken, http.StatusForbidden},
		{"no token", url.Values{"chat_jid": {"15550000001"}}, "", http.StatusUnauthorized},
		{"missing chat", url.Values{}, testToken, http.StatusBadRequest},
		{"limit zero", url.Values{"chat_jid": {"15550000001"}, "limit": {"0"}}, testToken, http.StatusBadRequest},
		{"limit too large", url.Values{"chat_jid": {"15550000001"}, "limit": {"501"}}, testToken, http.StatusBadRequest},
		{"limit not a number", url.Values{"chat_jid": {"15550000001"}, "limit": {"abc"}}, testToken, http.StatusBadRequest},
		{"bad cursor", url.Values{"chat_jid": {"15550000001"}, "before": {"%%%"}}, testToken, http.StatusBadRequest},
	}
	for _, tc := range cases {
		if code, _ := getHistory(t, handler, tc.query, tc.token); code != tc.want {
			t.Errorf("%s: status %d, want %d", tc.name, code, tc.want)
		}
	}
}
