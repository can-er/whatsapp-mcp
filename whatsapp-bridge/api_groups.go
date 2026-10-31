package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
)

// groupSummary is what GET /api/groups returns for each joined group: enough to find the JID to put
// in WHATSAPP_ALLOWED_CHATS, without messages or participant numbers
type groupSummary struct {
	JID               string `json:"jid"`
	Name              string `json:"name"`
	ParticipantsCount int    `json:"participants_count"`
}

// errNotConnected is returned by the groups source when the bridge is not connected to WhatsApp
var errNotConnected = errors.New("not connected to WhatsApp")

// clientGroups lists the groups the account has joined
func clientGroups(client *whatsmeow.Client) func(ctx context.Context) ([]groupSummary, error) {
	return func(ctx context.Context) ([]groupSummary, error) {
		if client == nil || !client.IsConnected() {
			return nil, errNotConnected
		}
		groups, err := client.GetJoinedGroups(ctx)
		if err != nil {
			return nil, err
		}
		result := make([]groupSummary, 0, len(groups))
		for _, g := range groups {
			count := g.ParticipantCount
			if count == 0 {
				count = len(g.Participants)
			}
			result = append(result, groupSummary{JID: g.JID.ToNonAD().String(), Name: g.Name, ParticipantsCount: count})
		}
		return result, nil
	}
}

// groupsHandler serves GET /api/groups
func groupsHandler(source func(ctx context.Context) ([]groupSummary, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if source == nil {
			http.Error(w, "Groups unavailable", http.StatusServiceUnavailable)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		groups, err := source(ctx)
		if errors.Is(err, errNotConnected) {
			http.Error(w, "Not connected to WhatsApp", http.StatusServiceUnavailable)
			return
		}
		if err != nil {
			http.Error(w, "Failed to list groups", http.StatusBadGateway)
			return
		}
		sort.SliceStable(groups, func(i, j int) bool {
			return strings.ToLower(groups[i].Name) < strings.ToLower(groups[j].Name)
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]groupSummary{"groups": groups})
	}
}
