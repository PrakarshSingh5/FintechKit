package client

import (
	"context"
	"fmt"

	"github.com/PrakarshSingh5/fintechkit/pkg/auth"
	"github.com/PrakarshSingh5/fintechkit/pkg/reliability"
)

// ProviderConfig holds configuration for creating a provider
type ProviderConfig struct {
	Name        string
	Credentials *auth.Credentials

	// Reliability settings
	RetryPolicy     *reliability.RetryPolicy
	RateLimitConfig *reliability.RateLimitConfig
	CircuitBreaker  *reliability.CircuitBreakerConfig
}

// Factory creates provider instances with built-in reliability features
type Factory struct {
	authManager *auth.Manager
	providers   map[string]ProviderConstructor
}

// ProviderConstructor is a function that creates a new provider instance
type ProviderConstructor func(config *ProviderConfig) (Provider, error)

// NewFactory creates a new provider factory
func NewFactory(authManager *auth.Manager) *Factory {
	return &Factory{
		authManager: authManager,
		providers:   make(map[string]ProviderConstructor),
	}
}

// Register registers a provider constructor
func (f *Factory) Register(name string, constructor ProviderConstructor) {
	f.providers[name] = constructor
}

// Create creates a provider instance with reliability features
func (f *Factory) Create(ctx context.Context, config *ProviderConfig) (Provider, error) {
	constructor, ok := f.providers[config.Name]
	if !ok {
		return nil, fmt.Errorf("provider %s not registered", config.Name)
	}

	// Retrieve credentials if not provided
	if config.Credentials == nil {
		creds, err := f.authManager.GetCredentials(ctx, config.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get credentials for %s: %w", config.Name, err)
		}
		config.Credentials = creds
	}

	// Create the provider
	provider, err := constructor(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider %s: %w", config.Name, err)
	}

	// Wrap with reliability features if configured
	wrapped := provider

	if config.RetryPolicy != nil {
		wrapped = &RetryWrapper{
			provider: wrapped,
			policy:   config.RetryPolicy,
		}
	}

	if config.RateLimitConfig != nil {
		limiter := reliability.NewRateLimiter(config.RateLimitConfig)
		wrapped = &RateLimitWrapper{
			provider: wrapped,
			limiter:  limiter,
		}
	}

	if config.CircuitBreaker != nil {
		breaker := reliability.NewCircuitBreaker(config.Name, config.CircuitBreaker)
		wrapped = &CircuitBreakerWrapper{
			provider: wrapped,
			breaker:  breaker,
		}
	}

	return wrapped, nil
}

// RetryWrapper wraps a provider with retry logic
type RetryWrapper struct {
	provider Provider
	policy   *reliability.RetryPolicy
}

func (w *RetryWrapper) Name() string {
	return w.provider.Name()
}

func (w *RetryWrapper) Authenticate(ctx context.Context) error {
	return reliability.WithRetry(ctx, w.policy, func() error {
		return w.provider.Authenticate(ctx)
	})
}

func (w *RetryWrapper) HealthCheck(ctx context.Context) error {
	return reliability.WithRetry(ctx, w.policy, func() error {
		return w.provider.HealthCheck(ctx)
	})
}

// RateLimitWrapper wraps a provider with rate limiting
type RateLimitWrapper struct {
	provider Provider
	limiter  *reliability.RateLimiter
}

func (w *RateLimitWrapper) Name() string {
	return w.provider.Name()
}

func (w *RateLimitWrapper) Authenticate(ctx context.Context) error {
	if err := w.limiter.Wait(ctx); err != nil {
		return err
	}
	return w.provider.Authenticate(ctx)
}

func (w *RateLimitWrapper) HealthCheck(ctx context.Context) error {
	if err := w.limiter.Wait(ctx); err != nil {
		return err
	}
	return w.provider.HealthCheck(ctx)
}

// CircuitBreakerWrapper wraps a provider with circuit breaker
type CircuitBreakerWrapper struct {
	provider Provider
	breaker  *reliability.CircuitBreaker
}

func (w *CircuitBreakerWrapper) Name() string {
	return w.provider.Name()
}

func (w *CircuitBreakerWrapper) Authenticate(ctx context.Context) error {
	result, err := w.breaker.Execute(func() (interface{}, error) {
		return nil, w.provider.Authenticate(ctx)
	})
	_ = result
	return err
}

func (w *CircuitBreakerWrapper) HealthCheck(ctx context.Context) error {
	result, err := w.breaker.Execute(func() (interface{}, error) {
		return nil, w.provider.HealthCheck(ctx)
	})
	_ = result
	return err
}
