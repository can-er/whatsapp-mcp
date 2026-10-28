package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// resolveMediaPath returns the real path of a file the API may send, or an error.
// Only regular files located (after resolving symlinks) under mediaDir are allowed;
// sending files is disabled when mediaDir is empty. Relative paths are relative to mediaDir.
func resolveMediaPath(mediaDir, mediaPath string) (string, error) {
	if strings.TrimSpace(mediaDir) == "" {
		return "", fmt.Errorf("sending files is disabled: set WHATSAPP_MEDIA_DIR to the folder files may be sent from")
	}

	root, err := filepath.Abs(mediaDir)
	if err != nil {
		return "", fmt.Errorf("invalid WHATSAPP_MEDIA_DIR: %v", err)
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("invalid WHATSAPP_MEDIA_DIR: %v", err)
	}

	if !filepath.IsAbs(mediaPath) {
		mediaPath = filepath.Join(root, mediaPath)
	}
	resolved, err := filepath.EvalSymlinks(filepath.Clean(mediaPath))
	if err != nil {
		return "", fmt.Errorf("media file not found: %s", mediaPath)
	}

	rel, err := filepath.Rel(root, resolved)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("media file must be inside WHATSAPP_MEDIA_DIR: %s", mediaPath)
	}

	info, err := os.Stat(resolved)
	if err != nil || !info.Mode().IsRegular() {
		return "", fmt.Errorf("media path is not a regular file: %s", mediaPath)
	}
	return resolved, nil
}
