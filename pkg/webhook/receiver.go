package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Event represents a webhook event
type Event struct {
	ID        string
	Type      string
	Provider  string
	Timestamp time.Time
	Data      json.RawMessage
	Signature string
}

// Handler is a function that processes a webhook event
type Handler func(ctx context.Context, event *Event) error

// Receiver manages webhook reception and verification
type Receiver struct {
	handlers  map[string][]Handler
	verifiers map[string]SignatureVerifier
}

// SignatureVerifier verifies webhook signatures
type SignatureVerifier interface {
	Verify(payload []byte, signature string) error
}

// NewReceiver creates a new webhook receiver
func NewReceiver() *Receiver {
	return &Receiver{
		handlers:  make(map[string][]Handler),
		verifiers: make(map[string]SignatureVerifier),
	}
}

// RegisterHandler registers a handler for a specific event type
func (r *Receiver) RegisterHandler(eventType string, handler Handler) {
	r.handlers[eventType] = append(r.handlers[eventType], handler)
}

// RegisterVerifier registers a signature verifier for a provider
func (r *Receiver) RegisterVerifier(provider string, verifier SignatureVerifier) {
	r.verifiers[provider] = verifier
}

// ProcessEvent processes an incoming webhook event
func (r *Receiver) ProcessEvent(ctx context.Context, provider string, payload []byte, signature string) error {
	// Verify signature
	verifier, ok := r.verifiers[provider]
	if ok {
		if err := verifier.Verify(payload, signature); err != nil {
			return fmt.Errorf("signature verification failed: %w", err)
		}
	}

	// Parse event
	var event Event
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	event.Provider = provider
	event.Signature = signature

	// Get handlers for this event type
	handlers, ok := r.handlers[event.Type]
	if !ok {
		// No handlers registered, but not an error
		return nil
	}

	// Execute all handlers
	for _, handler := range handlers {
		if err := handler(ctx, &event); err != nil {
			return fmt.Errorf("handler error for event %s: %w", event.Type, err)
		}
	}

	return nil
}

// HMACVerifier implements HMAC signature verification
type HMACVerifier struct {
	secret string
}

// NewHMACVerifier creates a new HMAC verifier
func NewHMACVerifier(secret string) *HMACVerifier {
	return &HMACVerifier{secret: secret}
}

// Verify verifies an HMAC signature
func (v *HMACVerifier) Verify(payload []byte, signature string) error {
	mac := hmac.New(sha256.New, []byte(v.secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedMAC)) {
		return errors.New("invalid signature")
	}

	return nil
}

// StripeVerifier implements Stripe webhook signature verification
type StripeVerifier struct {
	secret string
}

// NewStripeVerifier creates a new Stripe verifier
func NewStripeVerifier(secret string) *StripeVerifier {
	return &StripeVerifier{secret: secret}
}

// Verify verifies a Stripe webhook signature
func (v *StripeVerifier) Verify(payload []byte, signature string) error {
	// Real implementation would:
	// 1. Extract timestamp and signatures from header
	// 2. Construct signed payload
	// 3. Compute expected signature
	// 4. Compare signatures
	// 5. Check timestamp tolerance

	// Simplified for demonstration
	return nil
}

// IdempotencyTracker tracks processed events to prevent duplicates
type IdempotencyTracker struct {
	processed map[string]time.Time
}

// NewIdempotencyTracker creates a new idempotency tracker
func NewIdempotencyTracker() *IdempotencyTracker {
	return &IdempotencyTracker{
		processed: make(map[string]time.Time),
	}
}

// IsProcessed checks if an event has been processed
func (t *IdempotencyTracker) IsProcessed(eventID string) bool {
	_, exists := t.processed[eventID]
	return exists
}

// MarkProcessed marks an event as processed
func (t *IdempotencyTracker) MarkProcessed(eventID string) {
	t.processed[eventID] = time.Now()
}

// Cleanup removes old processed events (call periodically)
func (t *IdempotencyTracker) Cleanup(maxAge time.Duration) {
	cutoff := time.Now().Add(-maxAge)
	for id, processedAt := range t.processed {
		if processedAt.Before(cutoff) {
			delete(t.processed, id)
		}
	}
}

// IdempotentReceiver wraps a receiver with idempotency checks
type IdempotentReceiver struct {
	receiver *Receiver
	tracker  *IdempotencyTracker
}

// NewIdempotentReceiver creates a receiver with idempotency
func NewIdempotentReceiver(receiver *Receiver) *IdempotentReceiver {
	return &IdempotentReceiver{
		receiver: receiver,
		tracker:  NewIdempotencyTracker(),
	}
}

// ProcessEvent processes an event with idempotency checks
func (r *IdempotentReceiver) ProcessEvent(ctx context.Context, provider string, payload []byte, signature string) error {
	// Parse event to get ID
	var event Event
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	// Check idempotency
	if r.tracker.IsProcessed(event.ID) {
		return nil // Already processed, skip
	}

	// Process event
	if err := r.receiver.ProcessEvent(ctx, provider, payload, signature); err != nil {
		return err
	}

	// Mark as processed
	r.tracker.MarkProcessed(event.ID)
	return nil
}
