package main

import (
	"context"
	"testing"

	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func TestEventKinds(t *testing.T) {
	key := &waCommon.MessageKey{ID: proto.String("3EB0ORIGINAL")}
	editType, revokeType := waE2E.ProtocolMessage_MESSAGE_EDIT, waE2E.ProtocolMessage_REVOKE

	cases := []struct {
		name     string
		msg      *waE2E.Message
		kind     string
		text     string
		reaction *string
	}{
		{"edit", &waE2E.Message{ProtocolMessage: &waE2E.ProtocolMessage{
			Type: &editType, Key: key, EditedMessage: &waE2E.Message{Conversation: proto.String("corrected text")},
		}}, "edit", "corrected text", nil},
		{"revoke", &waE2E.Message{ProtocolMessage: &waE2E.ProtocolMessage{Type: &revokeType, Key: key}}, "revoke", "", nil},
		{"reaction", &waE2E.Message{ReactionMessage: &waE2E.ReactionMessage{Key: key, Text: proto.String("👍")}}, "reaction", "", proto.String("👍")},
		{"reaction removed", &waE2E.Message{ReactionMessage: &waE2E.ReactionMessage{Key: key, Text: proto.String("")}}, "reaction", "", proto.String("")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ev, downloadable := buildMessageEvent(newTestMessage("3EB0UPDATE", tc.msg), "", testResolver, testNow)
			if ev.Type != tc.kind || ev.TargetMessageID == nil || *ev.TargetMessageID != "3EB0ORIGINAL" || ev.Text != tc.text {
				t.Fatalf("type=%s target=%v text=%q", ev.Type, ev.TargetMessageID, ev.Text)
			}
			if (tc.reaction == nil) != (ev.Reaction == nil) || (tc.reaction != nil && *ev.Reaction != *tc.reaction) {
				t.Fatalf("reaction=%v, want %v", ev.Reaction, tc.reaction)
			}
			if ev.Media != nil || downloadable != nil {
				t.Fatal("updates carry no media")
			}
		})
	}

	// Regular messages keep type "message" without the update fields
	ev, _ := buildMessageEvent(newTestMessage("3EB0PLAIN", &waE2E.Message{Conversation: proto.String("hi")}), "", testResolver, testNow)
	if ev.Type != "message" || ev.TargetMessageID != nil || ev.Reaction != nil {
		t.Fatalf("plain message: %+v", ev)
	}
}

func TestEventUpdateWrittenToInbox(t *testing.T) {
	s := newTestSink(t, nil)
	msg := newTestMessage("3EB0REACT", &waE2E.Message{ReactionMessage: &waE2E.ReactionMessage{
		Key: &waCommon.MessageKey{ID: proto.String("3EB0PREVIEW")}, Text: proto.String("👍"),
	}})
	if _, err := s.Process(context.Background(), msg, ""); err != nil {
		t.Fatal(err)
	}
	ev, _ := readSingleEvent(t, s)
	if ev["type"] != "reaction" || ev["target_message_id"] != "3EB0PREVIEW" || ev["reaction"] != "👍" || ev["media"] != nil {
		t.Fatalf("unexpected event: %v", ev)
	}
}
