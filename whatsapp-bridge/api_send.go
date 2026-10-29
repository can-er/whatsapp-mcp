package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Fork-specific send API hardening (ADR-0006): dedicated send token, idempotent sends through
// client_msg_id, append-only audit log, and a health endpoint.

// sendFunc sends a message and returns the outcome, a status message and the WhatsApp message ID
type sendFunc func(recipient, message, mediaPath string) (success bool, status string, messageID string)

// healthStatus is returned by GET /api/health
type healthStatus struct {
	Connected   bool       `json:"connected"`
	LoggedIn    bool       `json:"logged_in"`
	LastEventAt *time.Time `json:"last_event_at"`
}

// apiOptions holds the replaceable parts of the REST API (the defaults use the WhatsApp client)
type apiOptions struct {
	send    sendFunc
	sendLog *sendLog
	health  func() healthStatus
}

// apiTokens are the secrets accepted by the REST API. When send is set, only it opens /api/send,
// so a process holding the bridge token can read (download, health) but not send.
type apiTokens struct {
	bridge string
	send   string
}

// sendRecord is one line of the audit log. The message text itself is never logged, only its hash.
type sendRecord struct {
	Time        time.Time `json:"time"`
	Recipient   string    `json:"recipient"`
	ClientMsgID string    `json:"client_msg_id,omitempty"`
	TextSHA256  string    `json:"text_sha256"`
	HasMedia    bool      `json:"has_media"`
	Success     bool      `json:"success"`
	MessageID   string    `json:"message_id,omitempty"`
	Status      string    `json:"status"`
}

// sendLog appends every send to an audit file and remembers successful client_msg_id values,
// rebuilt from the file on start, so a retried request never sends a second message.
type sendLog struct {
	mu   sync.Mutex
	path string
	done map[string]sendRecord
	now  func() time.Time
}

func openSendLog(path string) (*sendLog, error) {
	l := &sendLog{path: path, done: make(map[string]sendRecord), now: time.Now}

	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open send audit log: %v", err)
	}
	defer f.Close()
	if err := os.Chmod(path, 0600); err != nil {
		return nil, fmt.Errorf("failed to restrict send audit log permissions: %v", err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		var rec sendRecord
		if json.Unmarshal(scanner.Bytes(), &rec) == nil && rec.Success && rec.ClientMsgID != "" {
			l.done[rec.ClientMsgID] = rec
		}
	}
	return l, scanner.Err()
}

func (l *sendLog) append(rec sendRecord) error {
	line, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(l.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(line, '\n'))
	return err
}

// Send sends through send unless clientMsgID was already sent successfully, in which case the
// recorded result is returned with duplicate=true. Sends are serialized, so concurrent retries
// of the same request cannot both go through.
func (l *sendLog) Send(clientMsgID, recipient, text, mediaPath string, send sendFunc) (rec sendRecord, duplicate bool) {
	if l == nil {
		success, status, id := send(recipient, text, mediaPath)
		return sendRecord{Recipient: recipient, Success: success, Status: status, MessageID: id}, false
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if clientMsgID != "" {
		if prev, ok := l.done[clientMsgID]; ok {
			return prev, true
		}
	}

	success, status, id := send(recipient, text, mediaPath)
	sum := sha256.Sum256([]byte(text))
	rec = sendRecord{
		Time:        l.now().UTC(),
		Recipient:   recipient,
		ClientMsgID: clientMsgID,
		TextSHA256:  hex.EncodeToString(sum[:]),
		HasMedia:    mediaPath != "",
		Success:     success,
		MessageID:   id,
		Status:      status,
	}
	if err := l.append(rec); err != nil {
		fmt.Printf("Failed to write send audit log: %v\n", err)
	}
	if success && clientMsgID != "" {
		l.done[clientMsgID] = rec
	}
	return rec, false
}
