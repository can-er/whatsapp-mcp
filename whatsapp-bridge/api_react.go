package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// Reactions and quoted replies (fork-specific, ADR-0006). The target message is looked up in the
// bridge store, so the pipeline only passes a chat and a message ID; the store gives the exact chat
// JID and the sender (the participant, for groups).

// reactFunc sends a reaction to target in chat
type reactFunc func(chat types.JID, target storedMessage, emoji string) (success bool, status string, messageID string)

// replyFunc sends text to chat, quoting the quoted message
type replyFunc func(chat types.JID, text string, quoted storedMessage) (success bool, status string, messageID string)

// ReactRequest is the body of POST /api/react
type ReactRequest struct {
	ChatJID     string `json:"chat_jid"`
	MessageID   string `json:"message_id"`
	Emoji       string `json:"emoji"`
	ClientMsgID string `json:"client_msg_id,omitempty"`
}

// maxEmojiBytes bounds a reaction (an emoji with modifiers and joiners stays far below)
const maxEmojiBytes = 64

// errUnknownMessage is returned when the target message is not in the store for that chat
var errUnknownMessage = errors.New("message not found in this chat")

// FindMessage returns the stored message with this ID in one of the chat JIDs (phone number and LID
// forms of the same chat), or errUnknownMessage
func (store *MessageStore) FindMessage(chatJIDs []string, id string) (storedMessage, error) {
	var m storedMessage
	if len(chatJIDs) == 0 || id == "" {
		return m, errUnknownMessage
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(chatJIDs)), ",")
	args := []interface{}{id}
	for _, jid := range chatJIDs {
		args = append(args, jid)
	}
	err := store.db.QueryRow(`SELECT id, chat_jid, COALESCE(sender, ''), COALESCE(content, ''), timestamp, is_from_me,
			COALESCE(media_type, ''), COALESCE(filename, '')
		FROM messages WHERE id = ? AND chat_jid IN (`+placeholders+`) LIMIT 1`, args...).
		Scan(&m.ID, &m.ChatJID, &m.Sender, &m.Content, &m.Timestamp, &m.IsFromMe, &m.MediaType, &m.Filename)
	if errors.Is(err, sql.ErrNoRows) {
		return m, errUnknownMessage
	}
	return m, err
}

// findChatMessage looks up a message of an (allowed) chat given as a phone number or JID
func findChatMessage(client *whatsmeow.Client, store *MessageStore, chat, id string) (storedMessage, error) {
	if store == nil {
		return storedMessage{}, fmt.Errorf("message store unavailable")
	}
	return store.FindMessage(historyChatJIDs(client, chat), id)
}

// validEmoji accepts a short reaction without control characters
func validEmoji(emoji string) bool {
	if emoji == "" || len(emoji) > maxEmojiBytes {
		return false
	}
	for _, r := range emoji {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// messageKey mirrors whatsmeow's BuildMessageKey: the participant is only set for others' messages
// in group-like chats
func messageKey(chat, sender types.JID, fromMe bool, id string) *waProto.MessageKey {
	key := &waProto.MessageKey{RemoteJID: proto.String(chat.String()), FromMe: proto.Bool(fromMe), ID: proto.String(id)}
	if !fromMe && !sender.IsEmpty() && chat.Server != types.DefaultUserServer && chat.Server != types.HiddenUserServer &&
		chat.Server != types.MessengerServer {
		key.Participant = proto.String(sender.ToNonAD().String())
	}
	return key
}

// reactionMessage builds the reaction to target
func reactionMessage(chat types.JID, target storedMessage, emoji string, resolve lidResolver, now time.Time) *waProto.Message {
	sender := storedSenderJID(target.Sender, resolve)
	return &waProto.Message{ReactionMessage: &waProto.ReactionMessage{
		Key:               messageKey(chat, sender, target.IsFromMe, target.ID),
		Text:              proto.String(emoji),
		SenderTimestampMS: proto.Int64(now.UnixMilli()),
	}}
}

// quotedReplyMessage builds a text message quoting quoted, sent by participant
func quotedReplyMessage(text string, quoted storedMessage, participant types.JID) *waProto.Message {
	return &waProto.Message{ExtendedTextMessage: &waProto.ExtendedTextMessage{
		Text: proto.String(text),
		ContextInfo: &waProto.ContextInfo{
			StanzaID:      proto.String(quoted.ID),
			Participant:   proto.String(participant.ToNonAD().String()),
			QuotedMessage: &waProto.Message{Conversation: proto.String(quoted.Content)},
		},
	}}
}

// quotedParticipant is the author of the quoted message: the account itself (LID form in LID chats)
// or the stored sender
func quotedParticipant(client *whatsmeow.Client, chat types.JID, quoted storedMessage) types.JID {
	if quoted.IsFromMe {
		if client != nil && client.Store != nil {
			if chat.Server == types.HiddenUserServer && !client.Store.LID.IsEmpty() {
				return client.Store.LID
			}
			if client.Store.ID != nil {
				return *client.Store.ID
			}
		}
		return types.EmptyJID
	}
	return storedSenderJID(quoted.Sender, clientLIDResolver(client))
}

// sendReaction is the reactFunc used with a WhatsApp client
func sendReaction(client *whatsmeow.Client) reactFunc {
	return func(chat types.JID, target storedMessage, emoji string) (bool, string, string) {
		if client == nil || !client.IsConnected() {
			return false, "Not connected to WhatsApp", ""
		}
		msg := reactionMessage(chat, target, emoji, clientLIDResolver(client), time.Now())
		resp, err := client.SendMessage(context.Background(), chat, msg)
		if err != nil {
			return false, fmt.Sprintf("Error sending reaction: %v", err), ""
		}
		return true, "Reaction sent", resp.ID
	}
}

// sendReply is the replyFunc used with a WhatsApp client
func sendReply(client *whatsmeow.Client) replyFunc {
	return func(chat types.JID, text string, quoted storedMessage) (bool, string, string) {
		if client == nil || !client.IsConnected() {
			return false, "Not connected to WhatsApp", ""
		}
		msg := quotedReplyMessage(text, quoted, quotedParticipant(client, chat, quoted))
		resp, err := client.SendMessage(context.Background(), chat, msg)
		if err != nil {
			return false, fmt.Sprintf("Error sending message: %v", err), ""
		}
		return true, "Reply sent", resp.ID
	}
}

// writeSendResult answers a send, reply or reaction request
func writeSendResult(w http.ResponseWriter, rec sendRecord, duplicate bool) {
	w.Header().Set("Content-Type", "application/json")
	if !rec.Success {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(SendMessageResponse{
		Success:   rec.Success,
		Message:   rec.Status,
		MessageID: rec.MessageID,
		Duplicate: duplicate,
	})
}

// reactHandler serves POST /api/react
func reactHandler(client *whatsmeow.Client, store *MessageStore, allowed *chatAllowlist, opts apiOptions) http.HandlerFunc {
	react := opts.react
	if react == nil {
		react = sendReaction(client)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req ReactRequest
		if !decodeJSONBody(w, r, &req) {
			return
		}
		if req.ChatJID == "" || req.MessageID == "" {
			http.Error(w, "chat_jid and message_id are required", http.StatusBadRequest)
			return
		}
		if !validEmoji(req.Emoji) {
			http.Error(w, "emoji must be a short non-empty reaction", http.StatusBadRequest)
			return
		}
		if !allowed.AllowsChatString(client, req.ChatJID) {
			http.Error(w, "Chat is not in WHATSAPP_ALLOWED_CHATS", http.StatusForbidden)
			return
		}
		target, err := findChatMessage(client, store, req.ChatJID, req.MessageID)
		if errors.Is(err, errUnknownMessage) {
			http.Error(w, "Message not found in this chat", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Failed to look up the message", http.StatusInternalServerError)
			return
		}
		chat, err := types.ParseJID(target.ChatJID)
		if err != nil {
			http.Error(w, "Invalid stored chat", http.StatusInternalServerError)
			return
		}

		rec := sendRecord{Kind: "reaction", Recipient: req.ChatJID, TargetMessageID: target.ID, TextSHA256: textSHA256(req.Emoji)}
		rec, duplicate := opts.sendLog.Do(req.ClientMsgID, rec, func() (bool, string, string) {
			return react(chat, target, req.Emoji)
		})
		writeSendResult(w, rec, duplicate)
	}
}

// handleReply sends a /api/send request that has reply_to_message_id
func handleReply(w http.ResponseWriter, client *whatsmeow.Client, store *MessageStore, opts apiOptions, req SendMessageRequest) {
	if req.MediaPath != "" {
		http.Error(w, "reply_to_message_id is only supported for text messages", http.StatusBadRequest)
		return
	}
	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}
	quoted, err := findChatMessage(client, store, req.Recipient, req.ReplyToMessageID)
	if errors.Is(err, errUnknownMessage) {
		http.Error(w, "reply_to_message_id is not a stored message of this chat", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Failed to look up the quoted message", http.StatusInternalServerError)
		return
	}
	chat, err := types.ParseJID(quoted.ChatJID)
	if err != nil {
		http.Error(w, "Invalid stored chat", http.StatusInternalServerError)
		return
	}
	reply := opts.reply
	if reply == nil {
		reply = sendReply(client)
	}

	rec := sendRecord{Kind: "reply", Recipient: req.Recipient, TargetMessageID: quoted.ID, TextSHA256: textSHA256(req.Message)}
	rec, duplicate := opts.sendLog.Do(req.ClientMsgID, rec, func() (bool, string, string) {
		return reply(chat, req.Message, quoted)
	})
	writeSendResult(w, rec, duplicate)
}
