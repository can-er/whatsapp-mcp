package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	// Default address of the REST API: loopback only, so the API is not reachable from the network
	defaultListenAddr = "127.0.0.1:8080"

	// Maximum accepted size of a request body
	maxRequestBodyBytes = 1 << 20
)

// listenAddr returns the address the REST API listens on (WHATSAPP_BRIDGE_ADDR, loopback by default)
func listenAddr() string {
	if addr := strings.TrimSpace(os.Getenv("WHATSAPP_BRIDGE_ADDR")); addr != "" {
		return addr
	}
	return defaultListenAddr
}

// loadOrCreateBridgeToken returns the shared secret required by the REST API.
// WHATSAPP_BRIDGE_TOKEN takes precedence; otherwise the token is read from tokenPath,
// and generated (and stored with owner-only permissions) on first start.
func loadOrCreateBridgeToken(tokenPath string) (string, error) {
	if token := strings.TrimSpace(os.Getenv("WHATSAPP_BRIDGE_TOKEN")); token != "" {
		return token, nil
	}

	data, err := os.ReadFile(tokenPath)
	if err == nil {
		if token := strings.TrimSpace(string(data)); token != "" {
			return token, nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("failed to read API token: %v", err)
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate API token: %v", err)
	}
	token := hex.EncodeToString(buf)

	if err := os.WriteFile(tokenPath, []byte(token+"\n"), 0600); err != nil {
		return "", fmt.Errorf("failed to write API token: %v", err)
	}
	// WriteFile keeps the mode of an existing file, so enforce it explicitly
	if err := os.Chmod(tokenPath, 0600); err != nil {
		return "", fmt.Errorf("failed to restrict API token permissions: %v", err)
	}
	return token, nil
}

// requireAPIAuth rejects browser requests, unauthenticated requests and non-JSON bodies
// before they reach the API handlers, and caps the request body size.
func requireAPIAuth(tokens apiTokens, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Browsers always send Origin on cross-site POSTs; the API is never meant to be called from a web page
		if r.Header.Get("Origin") != "" {
			http.Error(w, "Cross-origin requests are not allowed", http.StatusForbidden)
			return
		}

		// /api/send requires the dedicated send token when one is configured
		token := tokens.bridge
		if tokens.send != "" && path.Clean(r.URL.Path) == "/api/send" {
			token = tokens.send
		}
		if subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte("Bearer "+token)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method == http.MethodPost {
			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil || mediaType != "application/json" {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
		next.ServeHTTP(w, r)
	})
}

// decodeJSONBody decodes the request body into v, writing the error response if it fails
func decodeJSONBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
		}
		return false
	}
	return true
}
