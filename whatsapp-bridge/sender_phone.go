package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// messages.sender keeps the user part WhatsApp addressed the message with (a phone number or,
// in LID-addressed groups, a LID): reactions and quoted replies rebuild the message key from it.
// messages.sender_phone (fork) holds the sender's phone number when it is known, so the database
// can be filtered by number.

// senderIdentity resolves phone numbers and LIDs with the whatsmeow LID map
type senderIdentity struct {
	pnForLID lidResolver
	lidForPN lidResolver
	own      types.JID // own phone-number JID, for messages sent by this account
}

func clientSenderIdentity(client *whatsmeow.Client) senderIdentity {
	id := senderIdentity{pnForLID: clientLIDResolver(client)}
	if client == nil || client.Store == nil {
		return id
	}
	if client.Store.ID != nil {
		id.own = client.Store.ID.ToNonAD()
	}
	if client.Store.LIDs != nil {
		id.lidForPN = func(pn types.JID) (types.JID, bool) {
			lid, err := client.Store.LIDs.GetLIDForPN(context.Background(), pn.ToNonAD())
			return lid, err == nil && !lid.IsEmpty()
		}
	}
	return id
}

// phoneOf returns the phone number of a sender known by its JID and optional alternate JID
func (s senderIdentity) phoneOf(jid, alt types.JID) string {
	pn, _ := splitPNAndLID(jid, alt, s.pnForLID)
	return pn.User
}

// messagePhone returns the phone number of the sender of a message, "" when unknown
func (s senderIdentity) messagePhone(sender, alt types.JID, isFromMe bool) string {
	if isFromMe && !s.own.IsEmpty() {
		return s.own.User
	}
	return s.phoneOf(sender, alt)
}

// storedSenderPhone returns the phone number of a stored sender value, or "" when it cannot be
// told for sure: a bare user part is a LID if the LID map knows it, a phone number if the LID map
// knows it as one, and is left unresolved otherwise.
func (s senderIdentity) storedSenderPhone(sender string, isFromMe bool) string {
	if isFromMe && !s.own.IsEmpty() {
		return s.own.User
	}
	sender = strings.TrimSpace(sender)
	if sender == "" {
		return ""
	}
	if strings.Contains(sender, "@") {
		jid, err := types.ParseJID(sender)
		if err != nil {
			return ""
		}
		return s.phoneOf(jid, types.EmptyJID)
	}
	user := strings.SplitN(sender, ":", 2)[0]
	if s.pnForLID != nil {
		if pn, ok := s.pnForLID(types.NewJID(user, types.HiddenUserServer)); ok {
			return pn.User
		}
	}
	if s.lidForPN != nil {
		if _, ok := s.lidForPN(types.NewJID(user, types.DefaultUserServer)); ok {
			return user
		}
	}
	return ""
}

// addColumnIfMissing adds a column to a table created by an earlier version of the bridge
func addColumnIfMissing(db *sql.DB, table, column, definition string) error {
	rows, err := db.Query("SELECT name FROM pragma_table_info(?)", table)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		if name == column {
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	rows.Close()
	_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, definition))
	return err
}

// SetSenderPhone records the phone number of a stored message's sender
func (store *MessageStore) SetSenderPhone(id, chatJID, phone string) error {
	if phone == "" {
		return nil
	}
	_, err := store.db.Exec("UPDATE messages SET sender_phone = ? WHERE id = ? AND chat_jid = ?", phone, id, chatJID)
	return err
}

// BackfillSenderPhones fills sender_phone for stored messages that do not have it yet, when the
// number can be resolved. It returns the number of updated messages.
func (store *MessageStore) BackfillSenderPhones(ident senderIdentity) (int64, error) {
	type key struct {
		sender   string
		isFromMe bool
	}
	rows, err := store.db.Query("SELECT DISTINCT sender, is_from_me FROM messages WHERE sender_phone IS NULL")
	if err != nil {
		return 0, err
	}
	var keys []key
	for rows.Next() {
		var k key
		var sender sql.NullString
		if err := rows.Scan(&sender, &k.isFromMe); err != nil {
			rows.Close()
			return 0, err
		}
		k.sender = sender.String
		keys = append(keys, k)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}

	var updated int64
	for _, k := range keys {
		phone := ident.storedSenderPhone(k.sender, k.isFromMe)
		if phone == "" {
			continue
		}
		res, err := store.db.Exec("UPDATE messages SET sender_phone = ? WHERE sender_phone IS NULL AND sender = ? AND is_from_me = ?", phone, k.sender, k.isFromMe)
		if err != nil {
			return updated, err
		}
		n, _ := res.RowsAffected()
		updated += n
	}
	return updated, nil
}
