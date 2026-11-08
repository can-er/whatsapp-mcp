package main

import (
	"testing"

	"go.mau.fi/whatsmeow/types/events"
)

func TestDirectPathFromURL(t *testing.T) {
	cases := map[string]string{
		// The signature of the path (oh, oe) must survive
		"https://mmg.whatsapp.net/v/t62.7118-24/12345_67890_n.enc?ccb=11-4&oh=01_SIGNATURE&oe=6A000000&_nc_sid=5e03e0&mms3=true": "/v/t62.7118-24/12345_67890_n.enc?ccb=11-4&oh=01_SIGNATURE&oe=6A000000&_nc_sid=5e03e0&mms3=true",
		"https://mmg.whatsapp.net/o1/v/t24/f2/m231/AQExample?ccb=9-4&oh=01_X&oe=6A000001":                                       "/o1/v/t24/f2/m231/AQExample?ccb=9-4&oh=01_X&oe=6A000001",
		"https://mmg.whatsapp.net/d/f/AbCd.enc": "/d/f/AbCd.enc",
		"":                                      "",
		"::not a url":                           "",
	}
	for raw, want := range cases {
		if got := extractDirectPathFromURL(raw); got != want {
			t.Errorf("extractDirectPathFromURL(%q) = %q, want %q", raw, got, want)
		}
	}
}

func TestMediaRetryWaiter(t *testing.T) {
	w := newMediaRetryWaiter()
	ch, done := w.register("MSG1")
	if w.deliver(&events.MediaRetry{MessageID: "OTHER"}) {
		t.Fatal("delivered a retry nobody waits for")
	}
	if !w.deliver(&events.MediaRetry{MessageID: "MSG1"}) {
		t.Fatal("retry not delivered")
	}
	if evt := <-ch; evt.MessageID != "MSG1" {
		t.Fatalf("got %q", evt.MessageID)
	}
	done()
	if w.deliver(&events.MediaRetry{MessageID: "MSG1"}) {
		t.Fatal("delivered after the download stopped waiting")
	}
	var nilWaiter *mediaRetryWaiter
	if nilWaiter.deliver(&events.MediaRetry{MessageID: "MSG1"}) {
		t.Fatal("nil waiter delivered")
	}
}
