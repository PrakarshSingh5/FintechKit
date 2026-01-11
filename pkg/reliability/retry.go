package reliability

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryPolicy defines retry behavior
type RetryPolicy struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	RandomizeJitter bool
	RetryableErrors []error
}

// DefaultRetryPolicy returns sensible defaults for retries
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:      3,
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		RandomizeJitter: true,
	}
}

// IsRetryable checks if an error should trigger a retry
func (p *RetryPolicy) IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// If no specific retryable errors defined, retry on all errors
	if len(p.RetryableErrors) == 0 {
		return true
	}

	for _, retryableErr := range p.RetryableErrors {
		if errors.Is(err, retryableErr) {
			return true
		}
	}

	return false
}

// CalculateBackoff calculates the backoff duration for a given attempt
func (p *RetryPolicy) CalculateBackoff(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}

	// Calculate exponential backoff
	backoff := float64(p.InitialInterval) * math.Pow(p.Multiplier, float64(attempt-1))

	// Cap at max interval
	if backoff > float64(p.MaxInterval) {
		backoff = float64(p.MaxInterval)
	}

	// Add jitter to prevent thundering herd
	if p.RandomizeJitter {
		jitter := rand.Float64() * 0.3 * backoff // 0-30% jitter
		backoff = backoff + jitter
	}

	return time.Duration(backoff)
}

// WithRetry executes a function with retry logic
func WithRetry(ctx context.Context, policy *RetryPolicy, fn func() error) error {
	if policy == nil {
		policy = DefaultRetryPolicy()
	}

	var lastErr error
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		// Execute the function
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if we should retry
		if !policy.IsRetryable(err) {
			return err
		}

		// Don't wait after the last attempt
		if attempt == policy.MaxRetries {
			break
		}

		// Calculate backoff
		backoff := policy.CalculateBackoff(attempt + 1)

		// Wait with context awareness
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		case <-time.After(backoff):
			// Continue to next attempt
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// RetryableFunc represents a function that can be retried and returns a value
type RetryableFunc[T any] func() (T, error)

// WithRetryTyped executes a function with retry logic and returns a typed value
func WithRetryTyped[T any](ctx context.Context, policy *RetryPolicy, fn RetryableFunc[T]) (T, error) {
	var zero T
	var result T
	var lastErr error

	if policy == nil {
		policy = DefaultRetryPolicy()
	}

	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		var err error
		result, err = fn()
		if err == nil {
			return result, nil
		}

		lastErr = err

		if !policy.IsRetryable(err) {
			return zero, err
		}

		if attempt == policy.MaxRetries {
			break
		}

		backoff := policy.CalculateBackoff(attempt + 1)

		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("retry cancelled: %w", ctx.Err())
		case <-time.After(backoff):
			// Continue to next attempt
		}
	}

	return zero, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// Common retryable errors
var (
	ErrTimeout            = errors.New("operation timeout")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrRateLimited        = errors.New("rate limited")
	ErrNetworkError       = errors.New("network error")
)

// StripeRetryPolicy returns a retry policy suitable for Stripe API
func StripeRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:      3,
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		RandomizeJitter: true,
		RetryableErrors: []error{ErrTimeout, ErrServiceUnavailable, ErrRateLimited},
	}
}

// PlaidRetryPolicy returns a retry policy suitable for Plaid API
func PlaidRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:      2,
		InitialInterval: 1 * time.Second,
		MaxInterval:     5 * time.Second,
		Multiplier:      1.5,
		RandomizeJitter: true,
		RetryableErrors: []error{ErrTimeout, ErrServiceUnavailable},
	}
}
