package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newGroupsHandler(source func(ctx context.Context) ([]groupSummary, error)) http.Handler {
	tokens := apiTokens{bridge: testToken, send: testSendToken}
	return requireAPIAuth(tokens, newAPIMux(nil, nil, parseChatAllowlist("15550000001"), apiOptions{groups: source}))
}

func getGroups(handler http.Handler, method, token string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	headers := map[string]string{"Content-Type": "application/json"}
	if token != "" {
		headers["Authorization"] = "Bearer " + token
	}
	handler.ServeHTTP(rec, newTestRequest(method, "/api/groups", "{}", headers))
	return rec
}

func TestGroupsAPI(t *testing.T) {
	handler := newGroupsHandler(func(ctx context.Context) ([]groupSummary, error) {
		return []groupSummary{
			{JID: "120363000000000002@g.us", Name: "Test Support", ParticipantsCount: 3},
			{JID: "120363000000000001@g.us", Name: "another group", ParticipantsCount: 12},
		}, nil
	})

	rec := getGroups(handler, http.MethodGet, testToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body)
	}
	var raw map[string][]map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	groups := raw["groups"]
	if len(groups) != 2 || groups[0]["name"] != "another group" || groups[1]["jid"] != "120363000000000002@g.us" ||
		groups[1]["participants_count"] != 3.0 {
		t.Errorf("unexpected groups (sorted by name expected): %v", groups)
	}
	for _, g := range groups {
		if len(g) != 3 {
			t.Errorf("only jid, name and participants_count may be returned: %v", g)
		}
	}
	if strings.Contains(rec.Body.String(), "s.whatsapp.net") || strings.Contains(rec.Body.String(), "@lid") {
		t.Error("the response must not contain participant JIDs")
	}

	for name, tc := range map[string]struct {
		method, token string
		want          int
	}{
		"no token":   {http.MethodGet, "", http.StatusUnauthorized},
		"send token": {http.MethodGet, testSendToken, http.StatusUnauthorized},
		"POST":       {http.MethodPost, testToken, http.StatusMethodNotAllowed},
	} {
		if rec := getGroups(handler, tc.method, tc.token); rec.Code != tc.want {
			t.Errorf("%s: status %d, want %d", name, rec.Code, tc.want)
		}
	}
}

func TestGroupsAPIErrors(t *testing.T) {
	for err, want := range map[error]int{
		errNotConnected:          http.StatusServiceUnavailable,
		errors.New("iq timeout"): http.StatusBadGateway,
	} {
		handler := newGroupsHandler(func(ctx context.Context) ([]groupSummary, error) { return nil, err })
		if rec := getGroups(handler, http.MethodGet, testToken); rec.Code != want {
			t.Errorf("%v: status %d, want %d", err, rec.Code, want)
		}
	}
	if rec := getGroups(newGroupsHandler(nil), http.MethodGet, testToken); rec.Code != http.StatusServiceUnavailable {
		t.Errorf("no source: status %d, want 503", rec.Code)
	}
	// Without a client, the production source reports "not connected"
	if _, err := clientGroups(nil)(context.Background()); !errors.Is(err, errNotConnected) {
		t.Errorf("clientGroups(nil): %v", err)
	}
}
