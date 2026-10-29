package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"report.pdf", "report.pdf"},
		{"../../evil.sh", "evil.sh"},
		{`..\..\evil.bat`, "evil.bat"},
		{"/etc/passwd", "passwd"},
		{`C:\Windows\win.ini`, "win.ini"},
		{"C:evil.txt", "C_evil.txt"},
		{"..", ""},
		{"...", ""},
		{"", ""},
		{"../", ""},
		{".bashrc", "bashrc"},
		{"CON", "_CON"},
		{"nul.txt", "_nul.txt"},
		{"com1.tar.gz", "_com1.tar.gz"},
		{"a\x00b.txt", "ab.txt"},
		{"line\nbreak.txt", "linebreak.txt"},
		{"invoice\u202Efdp.exe", "invoicefdp.exe"},
		{`what?<>|*".txt`, "what______.txt"},
		{"trailing. . ", "trailing"},
		{"123456789@s.whatsapp.net", "123456789@s.whatsapp.net"},
		{"123:4@s.whatsapp.net", "123_4@s.whatsapp.net"},
	}
	for _, tc := range tests {
		if got := sanitizeFilename(tc.in); got != tc.want {
			t.Errorf("sanitizeFilename(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}

	long := sanitizeFilename(strings.Repeat("é", 300) + ".pdf")
	if len(long) > maxFilenameLength || !strings.HasSuffix(long, ".pdf") {
		t.Errorf("long name not truncated correctly: %d bytes, %q", len(long), long[len(long)-8:])
	}
}

func TestMediaFilePathStaysInChatDir(t *testing.T) {
	store := filepath.Join(t.TempDir(), "store")
	names := []string{
		"photo.jpg", "../../evil.sh", `..\..\evil.bat`, "/etc/passwd", `C:\x\y.txt`,
		"CON", "a\x00b", strings.Repeat("x", 1000),
	}
	chats := []string{"123@s.whatsapp.net", "123:4@s.whatsapp.net", "../../outside", `..\..\outside`}

	for _, chat := range chats {
		for _, name := range names {
			chatDir, path, err := mediaFilePath(store, chat, name)
			if err != nil {
				t.Fatalf("mediaFilePath(%q, %q): %v", chat, name, err)
			}
			if filepath.Dir(chatDir) != store {
				t.Errorf("chat dir %q escapes store %q", chatDir, store)
			}
			if filepath.Dir(path) != chatDir || filepath.Base(path) == ".." {
				t.Errorf("path %q is not a direct child of %q", path, chatDir)
			}
		}
	}

	for _, name := range []string{"", "..", "../"} {
		if _, _, err := mediaFilePath(store, "123@s.whatsapp.net", name); err == nil {
			t.Errorf("mediaFilePath accepted unusable name %q", name)
		}
	}
}
