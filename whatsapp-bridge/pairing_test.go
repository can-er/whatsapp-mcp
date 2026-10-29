package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPairingPhoneHint(t *testing.T) {
	for in, want := range map[string]string{"+1 555 000 0001": "…01", "15550000042": "…42", "7": "", "": ""} {
		if got := phoneHint(in); got != want {
			t.Errorf("phoneHint(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPairingCodeDeliveredToInboxAndWebhook(t *testing.T) {
	received := make(chan []byte, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received <- body
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	sink := newTestSink(t, nil)
	sink.webhookURL = server.URL
	clock := testNow
	n := newPairingNotifier(sink, "+1 555 000 0001")
	n.now = func() time.Time { return clock }

	if err := n.CodeReady("ABCD-EFGH"); err != nil {
		t.Fatal(err)
	}
	ev, path := readSingleEvent(t, sink)
	if filepath.Base(path) != "20261029T093002Z_pair-20261029T093002Z.json" {
		t.Errorf("file name %s", filepath.Base(path))
	}
	want := map[string]interface{}{
		"version": 1.0, "type": "pairing_code", "id": "pair-20261029T093002Z", "code": "ABCD-EFGH",
		"expires_at": "2026-10-29T09:32:42Z", "phone_hint": "…01",
	}
	for k, v := range want {
		if ev[k] != v {
			t.Errorf("%s = %#v, want %#v", k, ev[k], v)
		}
	}
	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "5550000001") {
		t.Error("the event must not contain the full phone number")
	}
	select {
	case body := <-received:
		if string(body) != string(data) {
			t.Errorf("webhook body differs from the inbox event: %s", body)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("webhook not called")
	}

	// Paired event
	clock = clock.Add(90 * time.Second)
	if err := n.Paired(); err != nil {
		t.Fatal(err)
	}
	paired, err := os.ReadFile(filepath.Join(sink.inboxDir, "events", "20261029T093132Z_paired-20261029T093132Z.json"))
	if err != nil {
		t.Fatal(err)
	}
	var pe map[string]interface{}
	json.Unmarshal(paired, &pe)
	if pe["type"] != "paired" || pe["id"] != "paired-20261029T093132Z" || pe["version"] != 1.0 || pe["code"] != nil {
		t.Errorf("unexpected paired event: %v", pe)
	}
}

func TestPairingCodeMinInterval(t *testing.T) {
	clock := testNow
	n := newPairingNotifier(nil, "15550000001")
	n.now = func() time.Time { return clock }

	if !n.ShouldRequestCode() {
		t.Fatal("the first code must be requested")
	}
	n.CodeReady("AAAA-AAAA")
	clock = clock.Add(59 * time.Second)
	if n.ShouldRequestCode() {
		t.Error("a new code must not be requested within 60 s")
	}
	clock = clock.Add(time.Second)
	if !n.ShouldRequestCode() {
		t.Error("a new code may be requested after 60 s")
	}
	if err := n.Paired(); err != nil {
		t.Errorf("without a sink, events are only printed: %v", err)
	}
}

func TestPairTimeout(t *testing.T) {
	t.Setenv("WHATSAPP_PAIR_TIMEOUT_MINUTES", "")
	if pairTimeout() != 30*time.Minute {
		t.Error("default pairing timeout must be 30 minutes")
	}
	t.Setenv("WHATSAPP_PAIR_TIMEOUT_MINUTES", "5")
	if pairTimeout() != 5*time.Minute {
		t.Error("WHATSAPP_PAIR_TIMEOUT_MINUTES not applied")
	}
}
