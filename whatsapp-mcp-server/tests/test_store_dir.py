import importlib
import os
import sys

sys.path.insert(0, os.path.join(os.path.dirname(os.path.abspath(__file__)), ".."))


def load_whatsapp(monkeypatch, value):
    if value is None:
        monkeypatch.delenv("WHATSAPP_STORE_DIR", raising=False)
    else:
        monkeypatch.setenv("WHATSAPP_STORE_DIR", value)
    import whatsapp
    return importlib.reload(whatsapp)


def test_paths_follow_store_dir(monkeypatch, tmp_path):
    w = load_whatsapp(monkeypatch, str(tmp_path))
    assert w.MESSAGES_DB_PATH == os.path.join(str(tmp_path), "messages.db")
    assert w.BRIDGE_TOKEN_PATH == os.path.join(str(tmp_path), "bridge_token")

    (tmp_path / "bridge_token").write_text("token-from-store-dir\n")
    monkeypatch.delenv("WHATSAPP_BRIDGE_TOKEN", raising=False)
    monkeypatch.delenv("WHATSAPP_SEND_TOKEN", raising=False)
    assert w.bridge_headers() == {"Authorization": "Bearer token-from-store-dir"}


def test_default_store_dir(monkeypatch):
    w = load_whatsapp(monkeypatch, None)
    expected = os.path.join("whatsapp-bridge", "store", "messages.db")
    assert w.MESSAGES_DB_PATH.endswith(expected)
