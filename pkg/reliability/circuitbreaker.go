package reliability

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreakerConfig defines circuit breaker configuration
type CircuitBreakerConfig struct {
	MaxRequests      uint32        // Max requests allowed in half-open state
	Interval         time.Duration // Statistical window duration
	Timeout          time.Duration // Time to wait in open state before half-open
	FailureThreshold uint32        // Number of failures to open the circuit
	SuccessThreshold uint32        // Number of successes to close from half-open
	OnStateChange    func(name string, from gobreaker.State, to gobreaker.State)
}

// DefaultCircuitBreakerConfig returns sensible defaults
func DefaultCircuitBreakerConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		MaxRequests:      3,
		Interval:         60 * time.Second,
		Timeout:          30 * time.Second,
		FailureThreshold: 5,
		SuccessThreshold: 2,
	}
}

// CircuitBreaker wraps gobreaker with additional features
type CircuitBreaker struct {
	name    string
	breaker *gobreaker.CircuitBreaker
	config  *CircuitBreakerConfig
	mu      sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, config *CircuitBreakerConfig) *CircuitBreaker {
	if config == nil {
		config = DefaultCircuitBreakerConfig()
	}

	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= config.FailureThreshold
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			if config.OnStateChange != nil {
				config.OnStateChange(name, from, to)
			}
		},
	}

	return &CircuitBreaker{
		name:    name,
		breaker: gobreaker.NewCircuitBreaker(settings),
		config:  config,
	}
}

// Execute runs a function through the circuit breaker
func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	return cb.breaker.Execute(fn)
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() gobreaker.State {
	return cb.breaker.State()
}

// Counts returns current counts
func (cb *CircuitBreaker) Counts() gobreaker.Counts {
	return cb.breaker.Counts()
}

// Name returns the circuit breaker name
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

// CircuitBreakerManager manages multiple circuit breakers
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager creates a new manager
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// Get retrieves or creates a circuit breaker
func (m *CircuitBreakerManager) Get(name string, config *CircuitBreakerConfig) *CircuitBreaker {
	m.mu.RLock()
	breaker, exists := m.breakers[name]
	m.mu.RUnlock()

	if exists {
		return breaker
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if breaker, exists := m.breakers[name]; exists {
		return breaker
	}

	breaker = NewCircuitBreaker(name, config)
	m.breakers[name] = breaker
	return breaker
}

// GetAll returns all registered circuit breakers
func (m *CircuitBreakerManager) GetAll() map[string]*CircuitBreaker {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*CircuitBreaker, len(m.breakers))
	for k, v := range m.breakers {
		result[k] = v
	}
	return result
}

// FallbackHandler defines a fallback function when circuit is open
type FallbackHandler func() (interface{}, error)

// ExecuteWithFallback executes function with fallback on circuit open
func (cb *CircuitBreaker) ExecuteWithFallback(fn func() (interface{}, error), fallback FallbackHandler) (interface{}, error) {
	result, err := cb.Execute(fn)

	if err != nil {
		// Check if error is due to open circuit
		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			if fallback != nil {
				return fallback()
			}
		}
	}

	return result, err
}

// ErrorClassifier determines if an error should count as a failure
type ErrorClassifier func(err error) bool

// DefaultErrorClassifier treats all errors as failures
func DefaultErrorClassifier(err error) bool {
	return err != nil
}

// HTTPErrorClassifier only treats 5xx errors as failures
func HTTPErrorClassifier(err error) bool {
	if err == nil {
		return false
	}
	// In real implementation, check if error is HTTP 5xx
	// This is a simplified version
	return true
}

// CircuitBreakerStats holds statistics about a circuit breaker
type CircuitBreakerStats struct {
	Name               string
	State              string
	TotalRequests      uint32
	TotalSuccesses     uint32
	TotalFailures      uint32
	ConsecutiveSuccess uint32
	ConsecutiveFailure uint32
}

// GetStats returns current statistics
func (cb *CircuitBreaker) GetStats() CircuitBreakerStats {
	counts := cb.Counts()
	return CircuitBreakerStats{
		Name:               cb.name,
		State:              cb.State().String(),
		TotalRequests:      counts.Requests,
		TotalSuccesses:     counts.TotalSuccesses,
		TotalFailures:      counts.TotalFailures,
		ConsecutiveSuccess: counts.ConsecutiveSuccesses,
		ConsecutiveFailure: counts.ConsecutiveFailures,
	}
}

// Reset manually resets the circuit breaker
func (cb *CircuitBreaker) Reset() {
	// Note: gobreaker doesn't provide a direct reset method
	// In production, you might need to recreate the breaker
	cb.mu.Lock()
	defer cb.mu.Unlock()

	settings := gobreaker.Settings{
		Name:        cb.name,
		MaxRequests: cb.config.MaxRequests,
		Interval:    cb.config.Interval,
		Timeout:     cb.config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cb.config.FailureThreshold
		},
		OnStateChange: cb.config.OnStateChange,
	}

	cb.breaker = gobreaker.NewCircuitBreaker(settings)
}

// MonitorCircuitBreakers logs circuit breaker states periodically
func MonitorCircuitBreakers(manager *CircuitBreakerManager, interval time.Duration, logger func(string)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		breakers := manager.GetAll()
		for _, breaker := range breakers {
			stats := breaker.GetStats()
			msg := fmt.Sprintf("CircuitBreaker[%s]: State=%s, Requests=%d, Successes=%d, Failures=%d",
				stats.Name, stats.State, stats.TotalRequests, stats.TotalSuccesses, stats.TotalFailures)
			logger(msg)
		}
	}
}
