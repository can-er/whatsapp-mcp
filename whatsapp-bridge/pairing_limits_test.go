package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func TestDeliverRecreatesMissingEventsDir(t *testing.T) {
	sink := newTestSink(t, nil)
	if err := os.RemoveAll(filepath.Join(sink.inboxDir, "events")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(sink.inboxDir, "tmp")); err != nil {
		t.Fatal(err)
	}
	if err := sink.deliver("20261029T093002Z_x.json", []byte(`{"version":1}`)); err != nil {
		t.Fatalf("deliver must recreate the removed directories: %v", err)
	}
	readSingleEvent(t, sink)
}

func TestDeliverCallsWebhookWhenInboxFails(t *testing.T) {
	received := make(chan []byte, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received <- body
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	sink := newTestSink(t, nil)
	sink.webhookURL = server.URL
	// The inbox root is replaced by a regular file: the inbox cannot be written at all
	if err := os.RemoveAll(sink.inboxDir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sink.inboxDir, nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := sink.deliver("20261029T093002Z_x.json", []byte(`{"version":1}`)); err == nil {
		t.Error("the inbox failure must be reported")
	}
	select {
	case body := <-received:
		if string(body) != `{"version":1}` {
			t.Errorf("unexpected webhook body %s", body)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("webhook not called after an inbox failure")
	}
}

func TestPairingRequestLimit(t *testing.T) {
	clock := testNow
	n := newPairingNotifier(nil, "15550000001", 5)
	n.now = func() time.Time { return clock }

	for i := 0; i < 4; i++ {
		if !n.ShouldRequestCode() {
			t.Fatalf("request %d must be allowed", i+1)
		}
		n.CodeReady(fmt.Sprintf("CODE-%04d", i))
		clock = clock.Add(pairingSessionLifetime)
	}
	// A failed request counts too
	n.RequestFailed(errors.New("network error"))
	if !n.LimitReached() || n.ShouldRequestCode() {
		t.Error("no request may follow the 5th one")
	}
	if n.requests != 5 {
		t.Errorf("requests = %d, want 5", n.requests)
	}
}

func TestPairingRateLimitBackoff(t *testing.T) {
	clock := testNow
	n := newPairingNotifier(nil, "15550000001", 100)
	n.now = func() time.Time { return clock }

	rateLimited := fmt.Errorf("failed: %w", whatsmeow.ErrIQRateOverLimit)
	if !isRateLimited(rateLimited) || !isRateLimited(&whatsmeow.IQError{Code: 429}) || isRateLimited(errors.New("timeout")) {
		t.Fatal("rate-limit detection is wrong")
	}
	want := []time.Duration{5 * time.Minute, 10 * time.Minute, 20 * time.Minute, 40 * time.Minute, time.Hour, time.Hour}
	for i, w := range want {
		if got := n.RequestFailed(rateLimited); got != w {
			t.Errorf("backoff %d = %s, want %s", i+1, got, w)
		}
		if n.Wait() != w || n.ShouldRequestCode() {
			t.Errorf("no request may be made during the backoff (%d)", i+1)
		}
	}
	clock = clock.Add(time.Hour)
	if n.Wait() != 0 || !n.ShouldRequestCode() {
		t.Error("a request is allowed after the backoff")
	}
	// A delivered code resets the backoff; a non rate-limit error waits the base delay
	n.CodeReady("AAAA-AAAA")
	if got := n.RequestFailed(rateLimited); got != pairingRetryBackoff {
		t.Errorf("backoff after a code = %s, want %s", got, pairingRetryBackoff)
	}
	if got := n.RequestFailed(errors.New("timeout")); got != pairingRetryBackoff {
		t.Errorf("backoff after another error = %s, want %s", got, pairingRetryBackoff)
	}
}

func TestPairingStoppedEvent(t *testing.T) {
	received := make(chan []byte, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received <- body
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	sink := newTestSink(t, nil)
	sink.webhookURL = server.URL
	n := newPairingNotifier(sink, "+1 555 000 0001", 5)
	n.now = func() time.Time { return testNow }

	err := n.Stopped("max_codes", "device not linked after 5 pairing code requests", waLog.Noop)
	if !errors.Is(err, errPairingStopped) {
		t.Fatalf("the error must wrap errPairingStopped (exit status %d): %v", exitPairingStopped, err)
	}
	ev, path := readSingleEvent(t, sink)
	if filepath.Base(path) != "20261029T093002Z_pairing-stopped-20261029T093002Z.json" {
		t.Errorf("file name %s", filepath.Base(path))
	}
	if ev["type"] != "pairing_stopped" || ev["reason"] != "max_codes" || ev["phone_hint"] != "…01" || ev["version"] != 1.0 {
		t.Errorf("unexpected event %v", ev)
	}
	select {
	case body := <-received:
		var got map[string]interface{}
		json.Unmarshal(body, &got)
		if got["type"] != "pairing_stopped" {
			t.Errorf("unexpected webhook body %s", body)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("webhook not called")
	}
}

func TestPairMaxCodes(t *testing.T) {
	t.Setenv("WHATSAPP_PAIR_MAX_CODES", "")
	if pairMaxCodes() != 5 {
		t.Error("default must be 5 requests")
	}
	t.Setenv("WHATSAPP_PAIR_MAX_CODES", "3")
	if pairMaxCodes() != 3 {
		t.Error("WHATSAPP_PAIR_MAX_CODES not applied")
	}
	t.Setenv("WHATSAPP_PAIR_MAX_CODES", "0")
	if pairMaxCodes() != 5 {
		t.Error("invalid values fall back to the default")
	}
}
