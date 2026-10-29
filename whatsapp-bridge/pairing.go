package main

import (
	"context"
	"encoding/json"
	"errors"
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
// Requesting codes for hours gets the number rate-limited by WhatsApp, so the bridge gives up after
// WHATSAPP_PAIR_MAX_CODES requests (or the timeout), emits "pairing_stopped" and exits with status 3.

const (
	// Minimum delay between two pairing codes, so the owner is not flooded with codes
	pairingCodeMinInterval = 60 * time.Second
	// WhatsApp closes the login websocket (and invalidates the code) after about 160 seconds
	pairingSessionLifetime = 160 * time.Second
	// Delay before the next request after a failed one, doubled after each rate-limit answer (429)
	pairingRetryBackoff    = 5 * time.Minute
	pairingMaxRetryBackoff = time.Hour
	// Exit status when pairing is given up, so a service manager does not restart the bridge
	// (systemd: RestartPreventExitStatus=3); pairing is then started again on demand
	exitPairingStopped = 3
)

// errPairingStopped is returned (wrapped) by pairDevice when pairing is given up
var errPairingStopped = errors.New("pairing stopped")

// PairingEvent is the JSON document delivered for pairing codes, successful pairing and abandon
type PairingEvent struct {
	Version   int        `json:"version"`
	Type      string     `json:"type"`
	ID        string     `json:"id"`
	Code      string     `json:"code,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	PhoneHint string     `json:"phone_hint,omitempty"`
	Reason    string     `json:"reason,omitempty"` // pairing_stopped: max_codes, timeout or rate_limited
}

type pairingNotifier struct {
	sink        *eventSink // nil: codes are only printed in the terminal
	phone       string
	minInterval time.Duration
	maxRequests int
	requests    int // pairing code requests, successful or not
	lastCode    time.Time
	backoff     time.Duration
	notBefore   time.Time
	now         func() time.Time
}

func newPairingNotifier(sink *eventSink, phone string, maxRequests int) *pairingNotifier {
	return &pairingNotifier{sink: sink, phone: phone, minInterval: pairingCodeMinInterval, maxRequests: maxRequests, now: time.Now}
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
	if n.LimitReached() || n.now().Before(n.notBefore) {
		return false
	}
	return n.lastCode.IsZero() || n.now().Sub(n.lastCode) >= n.minInterval
}

// LimitReached reports whether WHATSAPP_PAIR_MAX_CODES requests were made
func (n *pairingNotifier) LimitReached() bool {
	return n.maxRequests > 0 && n.requests >= n.maxRequests
}

// Wait returns how long to wait before opening a new login session (after a failed request)
func (n *pairingNotifier) Wait() time.Duration {
	if d := n.notBefore.Sub(n.now()); d > 0 {
		return d
	}
	return 0
}

// RequestFailed records a failed code request and returns the delay before the next one. Rate-limit
// answers double the delay (up to one hour); other errors wait the base delay.
func (n *pairingNotifier) RequestFailed(err error) time.Duration {
	n.requests++
	switch {
	case !isRateLimited(err), n.backoff == 0:
		n.backoff = pairingRetryBackoff
	case 2*n.backoff > pairingMaxRetryBackoff:
		n.backoff = pairingMaxRetryBackoff
	default:
		n.backoff *= 2
	}
	n.notBefore = n.now().Add(n.backoff)
	return n.backoff
}

// isRateLimited reports whether WhatsApp refused the request because of too many attempts
func isRateLimited(err error) bool {
	var iqErr *whatsmeow.IQError
	return errors.Is(err, whatsmeow.ErrIQRateOverLimit) || (errors.As(err, &iqErr) && iqErr.Code == 429)
}

// CodeReady delivers a new pairing code
func (n *pairingNotifier) CodeReady(code string) error {
	now := n.now().UTC()
	n.requests++
	n.lastCode = now
	n.backoff = 0
	expires := now.Add(pairingSessionLifetime)
	fmt.Printf("\nPairing code %d/%d: %s (valid until %s). Enter it in WhatsApp > Linked devices > Link with phone number.\n",
		n.requests, n.maxRequests, code, expires.Format(time.RFC3339))
	return n.emit(PairingEvent{Type: "pairing_code", ID: "pair-" + now.Format("20060102T150405Z"), Code: code,
		ExpiresAt: &expires, PhoneHint: phoneHint(n.phone)})
}

// Paired announces that the device was linked
func (n *pairingNotifier) Paired() error {
	now := n.now().UTC()
	return n.emit(PairingEvent{Type: "paired", ID: "paired-" + now.Format("20060102T150405Z")})
}

// Stopped announces that pairing was given up and returns the error ending pairDevice
func (n *pairingNotifier) Stopped(reason, detail string, logger waLog.Logger) error {
	now := n.now().UTC()
	if err := n.emit(PairingEvent{Type: "pairing_stopped", ID: "pairing-stopped-" + now.Format("20060102T150405Z"),
		Reason: reason, PhoneHint: phoneHint(n.phone)}); err != nil {
		logger.Errorf("Failed to deliver pairing_stopped event: %v", err)
	}
	return fmt.Errorf("%w (%s): %s", errPairingStopped, reason, detail)
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

// pairMaxCodes reads WHATSAPP_PAIR_MAX_CODES (default 5, about 13 minutes of codes)
func pairMaxCodes() int {
	if n, err := strconv.Atoi(strings.TrimSpace(os.Getenv("WHATSAPP_PAIR_MAX_CODES"))); err == nil && n > 0 {
		return n
	}
	return 5
}

// pairDevice links a new device. Without a notifier it shows QR codes in the terminal (one login
// session, as before). With a notifier it requests one pairing code per login session and opens new
// sessions until the device is linked, or gives up (errPairingStopped) after the maximum number of
// requests or the timeout.
func pairDevice(ctx context.Context, client *whatsmeow.Client, notifier *pairingNotifier, timeout time.Duration, logger waLog.Logger) error {
	deadline := time.Now().Add(timeout)
	for {
		if notifier != nil {
			if notifier.LimitReached() {
				return notifier.Stopped("max_codes", fmt.Sprintf("device not linked after %d pairing code requests", notifier.requests), logger)
			}
			if wait := notifier.Wait(); wait > 0 {
				if time.Now().Add(wait).After(deadline) {
					return notifier.Stopped("rate_limited", fmt.Sprintf("next pairing code request allowed in %s, after the %s timeout", wait, timeout), logger)
				}
				logger.Warnf("Waiting %s before requesting a new pairing code", wait)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(wait):
				}
			}
		}

		// One context per login session, cancelled to close the session early
		sessionCtx, cancel := context.WithCancel(ctx)
		qrChan, err := client.GetQRChannel(sessionCtx)
		if err != nil {
			cancel()
			return err
		}
		if err := client.Connect(); err != nil {
			cancel()
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
					// The first QR event means the login session is ready for a pairing code request.
					// At most one request per session, successful or not.
					requested = true
					code, err := client.PairPhone(sessionCtx, notifier.phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
					if err != nil {
						wait := notifier.RequestFailed(err)
						logger.Errorf("Failed to request pairing code (%d/%d), next attempt in %s: %v", notifier.requests, notifier.maxRequests, wait, err)
						// Close this login session; the next one opens after the delay
						cancel()
						continue
					}
					if err := notifier.CodeReady(code); err != nil {
						logger.Errorf("Failed to deliver pairing code event: %v", err)
					}
				}
			case whatsmeow.QRChannelSuccess.Event:
				cancel()
				if notifier != nil {
					if err := notifier.Paired(); err != nil {
						logger.Errorf("Failed to deliver paired event: %v", err)
					}
				}
				return nil
			case "error":
				logger.Errorf("Pairing error: %v", evt.Error)
			}
		}
		cancel()

		// The login session ended without pairing
		if notifier == nil {
			return fmt.Errorf("timeout waiting for QR code scan")
		}
		if time.Now().After(deadline) {
			return notifier.Stopped("timeout", fmt.Sprintf("device not linked within %s", timeout), logger)
		}
		client.Disconnect()
		logger.Infof("Pairing session ended, opening a new one")
	}
}
