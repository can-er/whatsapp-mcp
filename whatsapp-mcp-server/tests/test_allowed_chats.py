import os
import sqlite3
import sys

import pytest

sys.path.insert(0, os.path.join(os.path.dirname(os.path.abspath(__file__)), ".."))

import whatsapp  # noqa: E402

ALLOWED = "32400000001@s.whatsapp.net"
HIDDEN = "32400000002@s.whatsapp.net"
GROUP = "120363000000000001@g.us"

# Same schema as the bridge (whatsapp-bridge/main.go, NewMessageStore)
SCHEMA = """
CREATE TABLE chats (jid TEXT PRIMARY KEY, name TEXT, last_message_time TIMESTAMP);
CREATE TABLE messages (
    id TEXT, chat_jid TEXT, sender TEXT, content TEXT, timestamp TIMESTAMP, is_from_me BOOLEAN,
    media_type TEXT, filename TEXT, url TEXT, media_key BLOB, file_sha256 BLOB, file_enc_sha256 BLOB,
    file_length INTEGER, PRIMARY KEY (id, chat_jid), FOREIGN KEY (chat_jid) REFERENCES chats(jid)
);
"""


@pytest.fixture
def db(tmp_path, monkeypatch):
    path = tmp_path / "messages.db"
    conn = sqlite3.connect(path)
    conn.executescript(SCHEMA)
    for i, (jid, name) in enumerate([(ALLOWED, "Alice"), (HIDDEN, "Bob"), (GROUP, "Team")]):
        ts = f"2026-10-0{i + 1} 10:00:00+00:00"
        conn.execute("INSERT INTO chats VALUES (?, ?, ?)", (jid, name, ts))
        conn.execute(
            "INSERT INTO messages (id, chat_jid, sender, content, timestamp, is_from_me) VALUES (?, ?, ?, ?, ?, ?)",
            (f"m{i}", jid, jid.split("@")[0], f"hello from {name}", ts, False),
        )
    conn.commit()
    conn.close()
    monkeypatch.setattr(whatsapp, "MESSAGES_DB_PATH", str(path))
    monkeypatch.delenv("WHATSAPP_ALLOWED_CHATS", raising=False)
    return path


def message_context(message_id):
    try:
        return whatsapp.get_message_context(message_id)
    except ValueError as e:  # raised for unknown (or hidden) messages
        return str(e)


def visible_text():
    """Everything the read tools can return, as one string."""
    parts = [
        whatsapp.list_chats(limit=50),
        whatsapp.list_messages(limit=50, include_context=False),
        whatsapp.list_messages(query="hello", limit=50),
        whatsapp.search_contacts(""),
        whatsapp.get_contact_chats(HIDDEN),
        whatsapp.get_last_interaction(HIDDEN),
        whatsapp.get_chat(HIDDEN),
        whatsapp.get_direct_chat_by_contact("32400000002"),
        message_context("m1"),
        whatsapp.get_sender_name(HIDDEN),
        whatsapp.list_chats(limit=50, include_last_message=False),
        whatsapp.get_chat(GROUP),
    ]
    return repr(parts)


def test_everything_visible_without_allowlist(db):
    text = visible_text()
    for name in ("Alice", "Bob", "Team"):
        assert f"hello from {name}" in text


@pytest.mark.parametrize("value", [ALLOWED, "32400000001", "+32 400 000 001", f"32400000001:12@s.whatsapp.net, "])
def test_only_allowed_chats_visible(db, monkeypatch, value):
    monkeypatch.setenv("WHATSAPP_ALLOWED_CHATS", value)
    text = visible_text()
    assert "hello from Alice" in text
    assert "Bob" not in text and "hello from Bob" not in text
    assert "Team" not in text and GROUP not in text
    assert whatsapp.get_chat(HIDDEN) is None
    assert whatsapp.get_chat(ALLOWED).name == "Alice"


def test_allowlist_parsing(monkeypatch):
    monkeypatch.setenv("WHATSAPP_ALLOWED_CHATS", f" +15550000001 , {GROUP},15550000002@s.whatsapp.net,,not a jid@@,@x")
    assert whatsapp.allowed_chat_jids() == ["15550000001@s.whatsapp.net", GROUP, "15550000002@s.whatsapp.net"]
    monkeypatch.setenv("WHATSAPP_ALLOWED_CHATS", " , ")
    assert whatsapp.allowed_chat_jids() is None


def test_database_file_is_not_modified(db, monkeypatch):
    monkeypatch.setenv("WHATSAPP_ALLOWED_CHATS", ALLOWED)
    whatsapp.list_chats()
    conn = sqlite3.connect(db)
    names = {row[0] for row in conn.execute("SELECT name FROM sqlite_master")}
    conn.close()
    assert "allowed_chat_jids" not in names
