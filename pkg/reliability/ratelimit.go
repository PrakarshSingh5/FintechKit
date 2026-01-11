package reliability

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitConfig defines rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond float64
	Burst             int
	WaitTimeout       time.Duration
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	limiter *rate.Limiter
	config  *RateLimitConfig
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config *RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(config.RequestsPerSecond), config.Burst),
		config:  config,
	}
}

// Wait waits for permission to make a request
func (rl *RateLimiter) Wait(ctx context.Context) error {
	if rl.config.WaitTimeout > 0 {
		waitCtx, cancel := context.WithTimeout(ctx, rl.config.WaitTimeout)
		defer cancel()
		return rl.limiter.Wait(waitCtx)
	}
	return rl.limiter.Wait(ctx)
}

// Allow checks if a request is allowed immediately
func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}

// Reserve reserves a token and returns a reservation
func (rl *RateLimiter) Reserve() *rate.Reservation {
	return rl.limiter.Reserve()
}

// MultiTierRateLimiter manages multiple rate limiters for different tiers
type MultiTierRateLimiter struct {
	limiters map[string]*RateLimiter
	mu       sync.RWMutex
}

// NewMultiTierRateLimiter creates a new multi-tier rate limiter
func NewMultiTierRateLimiter() *MultiTierRateLimiter {
	return &MultiTierRateLimiter{
		limiters: make(map[string]*RateLimiter),
	}
}

// AddTier adds a rate limiter for a specific tier
func (m *MultiTierRateLimiter) AddTier(tier string, config *RateLimitConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limiters[tier] = NewRateLimiter(config)
}

// Wait waits for permission for a specific tier
func (m *MultiTierRateLimiter) Wait(ctx context.Context, tier string) error {
	m.mu.RLock()
	limiter, ok := m.limiters[tier]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("tier %s not found", tier)
	}

	return limiter.Wait(ctx)
}

// DistributedRateLimiter interface for distributed rate limiting (e.g., with Redis)
type DistributedRateLimiter interface {
	// AllowN checks if n requests are allowed for a key
	AllowN(ctx context.Context, key string, n int) (bool, error)

	// Wait waits for permission for a request
	Wait(ctx context.Context, key string) error
}

// ProviderRateLimits holds preconfigured rate limits for different providers
var ProviderRateLimits = map[string]*RateLimitConfig{
	"stripe": {
		RequestsPerSecond: 100.0 / 1.0, // 100 requests per second
		Burst:             25,
		WaitTimeout:       5 * time.Second,
	},
	"plaid": {
		RequestsPerSecond: 10.0 / 1.0, // 10 requests per second (conservative)
		Burst:             5,
		WaitTimeout:       10 * time.Second,
	},
	"truelayer": {
		RequestsPerSecond: 30.0 / 1.0, // 30 requests per second
		Burst:             10,
		WaitTimeout:       5 * time.Second,
	},
	"coingecko": {
		RequestsPerSecond: 10.0 / 60.0, // 10 requests per minute (free tier)
		Burst:             5,
		WaitTimeout:       15 * time.Second,
	},
}

// AdaptiveRateLimiter adjusts rate limits based on 429 responses
type AdaptiveRateLimiter struct {
	limiter      *RateLimiter
	baseConfig   *RateLimitConfig
	mu           sync.RWMutex
	backoffUntil time.Time
}

// NewAdaptiveRateLimiter creates a new adaptive rate limiter
func NewAdaptiveRateLimiter(config *RateLimitConfig) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		limiter:    NewRateLimiter(config),
		baseConfig: config,
	}
}

// Wait waits for permission, respecting backoff
func (a *AdaptiveRateLimiter) Wait(ctx context.Context) error {
	a.mu.RLock()
	backoffUntil := a.backoffUntil
	a.mu.RUnlock()

	// Wait for backoff period if needed
	if time.Now().Before(backoffUntil) {
		waitTime := time.Until(backoffUntil)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			// Backoff complete
		}
	}

	return a.limiter.Wait(ctx)
}

// OnRateLimitError should be called when a 429 response is received
func (a *AdaptiveRateLimiter) OnRateLimitError(retryAfter time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Set backoff period
	a.backoffUntil = time.Now().Add(retryAfter)

	// Optionally reduce rate limit temporarily
	// This is a simple implementation; production might be more sophisticated
}

// Reset resets the rate limiter to base configuration
func (a *AdaptiveRateLimiter) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.limiter = NewRateLimiter(a.baseConfig)
	a.backoffUntil = time.Time{}
}

// ErrRateLimitExceeded is returned when rate limit is exceeded
var ErrRateLimitExceeded = errors.New("rate limit exceeded")

// WaitOrError waits for permission or returns error if timeout
func (rl *RateLimiter) WaitOrError(ctx context.Context) error {
	if !rl.Allow() {
		return ErrRateLimitExceeded
	}
	return nil
}
