package main

import (
	"testing"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func TestCaptionStoredAsContent(t *testing.T) {
	cases := []struct {
		name string
		msg  *waE2E.Message
		want string
	}{
		{"text", &waE2E.Message{Conversation: proto.String("hello")}, "hello"},
		{"extended text", &waE2E.Message{ExtendedTextMessage: &waE2E.ExtendedTextMessage{Text: proto.String("quoted reply")}}, "quoted reply"},
		{"image caption", &waE2E.Message{ImageMessage: &waE2E.ImageMessage{Caption: proto.String("2. Package management")}}, "2. Package management"},
		{"video caption", &waE2E.Message{VideoMessage: &waE2E.VideoMessage{Caption: proto.String("screen recording")}}, "screen recording"},
		{"document caption", &waE2E.Message{DocumentMessage: &waE2E.DocumentMessage{Caption: proto.String("invoice"), FileName: proto.String("a.pdf")}}, "invoice"},
		{"image without caption", &waE2E.Message{ImageMessage: &waE2E.ImageMessage{}}, ""},
		{"audio", &waE2E.Message{AudioMessage: &waE2E.AudioMessage{}}, ""},
		{"nil", nil, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := extractTextContent(c.msg); got != c.want {
				t.Fatalf("extractTextContent = %q, want %q", got, c.want)
			}
			// The stored content and the exported event text must not diverge
			if c.msg != nil {
				if ev := messageText(c.msg); ev != c.want {
					t.Fatalf("messageText = %q, want %q", ev, c.want)
				}
			}
		})
	}

	// A media message without caption is still stored, as media
	mediaType, _, _, _, _, _, _ := extractMediaInfo(&waE2E.Message{ImageMessage: &waE2E.ImageMessage{}})
	if mediaType != "image" {
		t.Fatalf("media type = %q", mediaType)
	}
}
