package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Pairing by code (fork-specific, ADR-0006): with WHATSAPP_PAIR_PHONE, an unpaired bridge requests a
// pairing code instead of showing a QR code, and delivers it as a "pairing_code" event (inbox + webhook)
// so it can be relayed to the account owner. A "paired" event follows once the phone linked the device.

const (
	// Minimum delay between two pairing codes, so the owner is not flooded with codes
	pairingCodeMinInterval = 60 * time.Second
	// WhatsApp closes the login websocket (and invalidates the code) after about 160 seconds
	pairingSessionLifetime = 160 * time.Second
)

// PairingEvent is the JSON document delivered for pairing codes and successful pairing
type PairingEvent struct {
	Version   int        `json:"version"`
	Type      string     `json:"type"`
	ID        string     `json:"id"`
	Code      string     `json:"code,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	PhoneHint string     `json:"phone_hint,omitempty"`
}

type pairingNotifier struct {
	sink        *eventSink // nil: codes are only printed in the terminal
	phone       string
	minInterval time.Duration
	lastCode    time.Time
	now         func() time.Time
}

func newPairingNotifier(sink *eventSink, phone string) *pairingNotifier {
	return &pairingNotifier{sink: sink, phone: phone, minInterval: pairingCodeMinInterval, now: time.Now}
}

var nonDigits = regexp.MustCompile(`\D`)

// phoneHint keeps only the last two digits of the account number
func phoneHint(phone string) string {
	digits := nonDigits.ReplaceAllString(phone, "")
	if len(digits) < 2 {
		return ""
	}
	return "…" + digits[len(digits)-2:]
}

// ShouldRequestCode reports whether a new pairing code may be requested now
func (n *pairingNotifier) ShouldRequestCode() bool {
	return n.lastCode.IsZero() || n.now().Sub(n.lastCode) >= n.minInterval
}

// CodeReady delivers a new pairing code
func (n *pairingNotifier) CodeReady(code string) error {
	now := n.now().UTC()
	n.lastCode = now
	expires := now.Add(pairingSessionLifetime)
	fmt.Printf("\nPairing code: %s (valid until %s). Enter it in WhatsApp > Linked devices > Link with phone number.\n",
		code, expires.Format(time.RFC3339))
	return n.emit(PairingEvent{Type: "pairing_code", ID: "pair-" + now.Format("20060102T150405Z"), Code: code,
		ExpiresAt: &expires, PhoneHint: phoneHint(n.phone)})
}

// Paired announces that the device was linked
func (n *pairingNotifier) Paired() error {
	now := n.now().UTC()
	return n.emit(PairingEvent{Type: "paired", ID: "paired-" + now.Format("20060102T150405Z")})
}

func (n *pairingNotifier) emit(ev PairingEvent) error {
	if n.sink == nil {
		return nil
	}
	ev.Version = eventSchemaVersion
	body, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	// Same naming as message events, so files sort chronologically
	return n.sink.deliver(n.now().UTC().Format("20060102T150405Z")+"_"+ev.ID+".json", body)
}

// pairTimeout reads WHATSAPP_PAIR_TIMEOUT_MINUTES (default 30)
func pairTimeout() time.Duration {
	if minutes, err := strconv.Atoi(strings.TrimSpace(os.Getenv("WHATSAPP_PAIR_TIMEOUT_MINUTES"))); err == nil && minutes > 0 {
		return time.Duration(minutes) * time.Minute
	}
	return 30 * time.Minute
}

// pairDevice links a new device. Without a notifier it shows QR codes in the terminal (one login
// session, as before). With a notifier it requests one pairing code per login session and opens new
// sessions until the device is linked or the timeout expires.
func pairDevice(ctx context.Context, client *whatsmeow.Client, notifier *pairingNotifier, timeout time.Duration, logger waLog.Logger) error {
	deadline := time.Now().Add(timeout)
	for {
		qrChan, err := client.GetQRChannel(ctx)
		if err != nil {
			return err
		}
		if err := client.Connect(); err != nil {
			return err
		}

		requested := false
		for evt := range qrChan {
			switch evt.Event {
			case whatsmeow.QRChannelEventCode:
				if notifier == nil {
					fmt.Println("\nScan this QR code with your WhatsApp app:")
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				} else if !requested && notifier.ShouldRequestCode() {
					// The first QR event means the login session is ready for a pairing code request
					code, err := client.PairPhone(ctx, notifier.phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
					if err != nil {
						logger.Errorf("Failed to request pairing code: %v", err)
						continue
					}
					requested = true
					if err := notifier.CodeReady(code); err != nil {
						logger.Warnf("Failed to deliver pairing code event: %v", err)
					}
				}
			case whatsmeow.QRChannelSuccess.Event:
				if notifier != nil {
					if err := notifier.Paired(); err != nil {
						logger.Warnf("Failed to deliver paired event: %v", err)
					}
				}
				return nil
			case "error":
				logger.Errorf("Pairing error: %v", evt.Error)
			}
		}

		// The login session ended without pairing
		if notifier == nil {
			return fmt.Errorf("timeout waiting for QR code scan")
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("device not linked within %s", timeout)
		}
		client.Disconnect()
		logger.Infof("Pairing session expired, opening a new one")
	}
}
