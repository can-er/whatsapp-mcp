package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxFilenameLength = 200

// Device names Windows refuses as file names, with or without an extension
var windowsReservedNames = map[string]bool{
	"CON": true, "PRN": true, "AUX": true, "NUL": true,
	"COM1": true, "COM2": true, "COM3": true, "COM4": true, "COM5": true, "COM6": true, "COM7": true, "COM8": true, "COM9": true,
	"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true, "LPT5": true, "LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
}

// sanitizeFilename reduces an untrusted name (e.g. a document name chosen by the sender) to a
// single safe path component. It returns "" when nothing usable remains.
func sanitizeFilename(name string) string {
	// Keep only the last component, whatever the separator
	if i := strings.LastIndexAny(name, `/\`); i >= 0 {
		name = name[i+1:]
	}

	var b strings.Builder
	for _, r := range name {
		switch {
		case unicode.IsControl(r) || unicode.Is(unicode.Cf, r):
			// Drop control and invisible formatting characters (e.g. right-to-left override)
		case strings.ContainsRune(`<>:"|?*`, r):
			b.WriteRune('_')
		default:
			b.WriteRune(r)
		}
	}

	// No hidden files or "..", and Windows silently strips trailing dots and spaces
	name = strings.TrimLeft(strings.TrimSpace(b.String()), ".")
	name = strings.TrimRight(name, ". ")
	if name == "" {
		return ""
	}

	stem := strings.ToUpper(strings.TrimSpace(strings.SplitN(name, ".", 2)[0]))
	if windowsReservedNames[stem] {
		name = "_" + name
	}

	if len(name) > maxFilenameLength {
		ext := filepath.Ext(name)
		if len(ext) > 20 {
			ext = ""
		}
		stemPart := name[:maxFilenameLength-len(ext)]
		for !utf8.ValidString(stemPart) {
			stemPart = stemPart[:len(stemPart)-1]
		}
		name = stemPart + ext
	}
	return name
}

// localMediaPath returns where the media of a stored message is (or will be) downloaded.
// Rows stored by older versions may hold a sender-controlled name, so it is sanitized again, with a
// name generated from the message ID as fallback.
func localMediaPath(chatJID, mediaType, messageID, filename string) (string, string, error) {
	filename = sanitizeFilename(filename)
	if filename == "" {
		filename = mediaType + "_" + sanitizeFilename(messageID)
	}
	return mediaFilePath(storePath(), chatJID, filename)
}

// mediaFilePath returns the directory for a chat's media under storeDir and the path of the
// given file in it. Both the chat JID and the file name are sanitized, and the result is
// guaranteed to be a direct child of the chat directory.
func mediaFilePath(storeDir, chatJID, filename string) (string, string, error) {
	chatComponent := sanitizeFilename(chatJID)
	if chatComponent == "" {
		return "", "", fmt.Errorf("invalid chat JID")
	}
	safeName := sanitizeFilename(filename)
	if safeName == "" {
		return "", "", fmt.Errorf("invalid media file name")
	}

	chatDir := filepath.Join(storeDir, chatComponent)
	localPath := filepath.Join(chatDir, safeName)
	if filepath.Dir(localPath) != chatDir {
		return "", "", fmt.Errorf("invalid media file name")
	}
	return chatDir, localPath, nil
}
