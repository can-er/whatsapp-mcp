package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const testToken = "0123456789abcdef0123456789abcdef"

func newTestRequest(method, path, body string, headers map[string]string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}

func validHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + testToken,
		"Content-Type":  "application/json",
	}
}

func TestListenAddrDefaultsToLoopback(t *testing.T) {
	t.Setenv("WHATSAPP_BRIDGE_ADDR", "")
	if got := listenAddr(); got != "127.0.0.1:8080" {
		t.Fatalf("listenAddr() = %q, want 127.0.0.1:8080", got)
	}

	t.Setenv("WHATSAPP_BRIDGE_ADDR", "127.0.0.1:9090")
	if got := listenAddr(); got != "127.0.0.1:9090" {
		t.Fatalf("listenAddr() = %q, want override", got)
	}
}

func TestBridgeTokenCreatedOnceAndReused(t *testing.T) {
	t.Setenv("WHATSAPP_BRIDGE_TOKEN", "")
	path := filepath.Join(t.TempDir(), "bridge_token")

	first, err := loadOrCreateBridgeToken(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 64 {
		t.Fatalf("token length = %d, want 64 hex chars", len(first))
	}

	if runtime.GOOS != "windows" {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if perm := info.Mode().Perm(); perm != 0600 {
			t.Fatalf("token file mode = %o, want 600", perm)
		}
	}

	second, err := loadOrCreateBridgeToken(path)
	if err != nil {
		t.Fatal(err)
	}
	if second != first {
		t.Fatal("token was regenerated instead of reused")
	}
}

func TestBridgeTokenEnvOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bridge_token")
	t.Setenv("WHATSAPP_BRIDGE_TOKEN", "from-env")

	got, err := loadOrCreateBridgeToken(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "from-env" {
		t.Fatalf("token = %q, want value from WHATSAPP_BRIDGE_TOKEN", got)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("token file must not be created when the env variable is set")
	}
}

// Requests rejected by the API must never reach the WhatsApp client: the real routes are
// served with a nil client and nil store, so any rejected request reaching them would panic.
func TestAPIRejectsUnsafeRequests(t *testing.T) {
	handler := requireAPIAuth(apiTokens{bridge: testToken}, newAPIMux(nil, nil, nil, apiOptions{}))
	validBody := `{"recipient":"123","message":"hi"}`

	tests := []struct {
		name    string
		method  string
		body    string
		headers map[string]string
		want    int
	}{
		{"missing token", http.MethodPost, validBody, map[string]string{"Content-Type": "application/json"}, http.StatusUnauthorized},
		{"wrong token", http.MethodPost, validBody, map[string]string{"Authorization": "Bearer nope", "Content-Type": "application/json"}, http.StatusUnauthorized},
		{"token without scheme", http.MethodPost, validBody, map[string]string{"Authorization": testToken, "Content-Type": "application/json"}, http.StatusUnauthorized},
		{"text/plain body (simple CORS request)", http.MethodPost, validBody, map[string]string{"Authorization": "Bearer " + testToken, "Content-Type": "text/plain"}, http.StatusUnsupportedMediaType},
		{"missing content type", http.MethodPost, validBody, map[string]string{"Authorization": "Bearer " + testToken}, http.StatusUnsupportedMediaType},
		{"browser origin", http.MethodPost, validBody, map[string]string{"Authorization": "Bearer " + testToken, "Content-Type": "application/json", "Origin": "https://example.com"}, http.StatusForbidden},
		{"browser origin without token", http.MethodPost, validBody, map[string]string{"Content-Type": "text/plain", "Origin": "https://example.com"}, http.StatusForbidden},
		{"oversized body", http.MethodPost, `{"recipient":"` + strings.Repeat("a", maxRequestBodyBytes) + `"}`, validHeaders(), http.StatusRequestEntityTooLarge},
		{"invalid json", http.MethodPost, `{`, validHeaders(), http.StatusBadRequest},
		{"GET with token", http.MethodGet, "", map[string]string{"Authorization": "Bearer " + testToken}, http.StatusMethodNotAllowed},
	}

	for _, path := range []string{"/api/send", "/api/download"} {
		for _, tc := range tests {
			t.Run(path+"/"+tc.name, func(t *testing.T) {
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, newTestRequest(tc.method, path, tc.body, tc.headers))
				if rec.Code != tc.want {
					t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.want, rec.Body.String())
				}
			})
		}
	}
}

func TestAPIAcceptsAuthenticatedJSON(t *testing.T) {
	reached := false
	handler := requireAPIAuth(apiTokens{bridge: testToken}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		w.WriteHeader(http.StatusNoContent)
	}))

	headers := validHeaders()
	headers["Content-Type"] = "application/json; charset=utf-8"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, newTestRequest(http.MethodPost, "/api/send", `{}`, headers))

	if !reached || rec.Code != http.StatusNoContent {
		t.Fatalf("authenticated JSON request was rejected: status %d", rec.Code)
	}
}
