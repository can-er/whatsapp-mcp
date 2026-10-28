package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveMediaPath(t *testing.T) {
	base := t.TempDir()
	mediaDir := filepath.Join(base, "outbox")
	outside := filepath.Join(base, "secret.txt")
	allowed := filepath.Join(mediaDir, "photo.jpg")
	nested := filepath.Join(mediaDir, "sub", "doc.pdf")

	for _, dir := range []string{mediaDir, filepath.Dir(nested)} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{outside, allowed, nested} {
		if err := os.WriteFile(f, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	accepted := map[string]string{
		"absolute path inside":  allowed,
		"nested path inside":    nested,
		"relative to media dir": "photo.jpg",
	}
	for name, p := range accepted {
		t.Run("accept/"+name, func(t *testing.T) {
			got, err := resolveMediaPath(mediaDir, p)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if filepath.Base(got) != filepath.Base(p) {
				t.Fatalf("resolved to %q", got)
			}
		})
	}

	rejected := map[string]string{
		"absolute path outside":      outside,
		"relative traversal":         filepath.Join("..", "secret.txt"),
		"traversal from inside":      filepath.Join(mediaDir, "..", "secret.txt"),
		"media dir itself":           mediaDir,
		"directory inside":           filepath.Dir(nested),
		"missing file":               filepath.Join(mediaDir, "missing.jpg"),
		"sibling with common prefix": filepath.Join(base, "outbox-evil", "x.jpg"),
	}
	if err := os.MkdirAll(filepath.Join(base, "outbox-evil"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "outbox-evil", "x.jpg"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}

	// A symlink inside the media dir pointing outside must not be followed out of it
	link := filepath.Join(mediaDir, "link.txt")
	if err := os.Symlink(outside, link); err == nil {
		rejected["symlink pointing outside"] = link
	} else {
		t.Logf("symlinks unavailable, skipping symlink case: %v", err)
	}

	for name, p := range rejected {
		t.Run("reject/"+name, func(t *testing.T) {
			if got, err := resolveMediaPath(mediaDir, p); err == nil {
				t.Fatalf("expected rejection, got %q", got)
			}
		})
	}

	t.Run("reject/media dir not configured", func(t *testing.T) {
		if _, err := resolveMediaPath("", allowed); err == nil {
			t.Fatal("sending files must be disabled when WHATSAPP_MEDIA_DIR is empty")
		}
	})
}
