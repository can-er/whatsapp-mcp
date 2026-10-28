import os
import sys

import pytest

sys.path.insert(0, os.path.join(os.path.dirname(os.path.abspath(__file__)), ".."))

import whatsapp  # noqa: E402


@pytest.fixture
def layout(tmp_path, monkeypatch):
    media_dir = tmp_path / "outbox"
    (media_dir / "sub").mkdir(parents=True)
    (media_dir / "photo.jpg").write_bytes(b"x")
    (media_dir / "sub" / "doc.pdf").write_bytes(b"x")
    (tmp_path / "secret.txt").write_bytes(b"x")
    (tmp_path / "outbox-evil").mkdir()
    (tmp_path / "outbox-evil" / "x.jpg").write_bytes(b"x")
    monkeypatch.setenv("WHATSAPP_MEDIA_DIR", str(media_dir))
    return tmp_path, media_dir


@pytest.mark.parametrize("rel", ["photo.jpg", os.path.join("sub", "doc.pdf")])
def test_accepts_files_inside_media_dir(layout, rel):
    _, media_dir = layout
    for candidate in (str(media_dir / rel), rel):
        assert os.path.samefile(whatsapp.resolve_media_path(candidate), media_dir / rel)


@pytest.mark.parametrize("make_path", [
    lambda base, media: str(base / "secret.txt"),
    lambda base, media: os.path.join("..", "secret.txt"),
    lambda base, media: str(media / ".." / "secret.txt"),
    lambda base, media: str(media),
    lambda base, media: str(media / "sub"),
    lambda base, media: str(media / "missing.jpg"),
    lambda base, media: str(base / "outbox-evil" / "x.jpg"),
], ids=["outside", "relative-traversal", "traversal-from-inside", "media-dir", "directory", "missing", "common-prefix-sibling"])
def test_rejects_paths_outside_media_dir(layout, make_path):
    base, media_dir = layout
    with pytest.raises(ValueError):
        whatsapp.resolve_media_path(make_path(base, media_dir))


def test_rejects_symlink_pointing_outside(layout):
    base, media_dir = layout
    link = media_dir / "link.txt"
    try:
        os.symlink(base / "secret.txt", link)
    except (OSError, NotImplementedError):
        pytest.skip("symlinks unavailable")
    with pytest.raises(ValueError):
        whatsapp.resolve_media_path(str(link))


def test_sending_files_disabled_without_media_dir(layout, monkeypatch):
    _, media_dir = layout
    monkeypatch.delenv("WHATSAPP_MEDIA_DIR")
    with pytest.raises(ValueError, match="disabled"):
        whatsapp.resolve_media_path(str(media_dir / "photo.jpg"))


def test_send_tools_never_call_bridge_for_outside_file(layout, monkeypatch):
    base, _ = layout

    def fail(*args, **kwargs):
        raise AssertionError("the bridge must not be called")

    monkeypatch.setattr(whatsapp.requests, "post", fail)
    for send in (whatsapp.send_file, whatsapp.send_audio_message):
        success, message = send("123", str(base / "secret.txt"))
        assert not success
        assert "WHATSAPP_MEDIA_DIR" in message
