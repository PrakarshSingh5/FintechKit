package auth

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RotationPolicy defines when and how credentials should be rotated
type RotationPolicy struct {
	Enabled         bool
	RotationPeriod  time.Duration
	OverlapDuration time.Duration // Time both old and new credentials are valid
}

// RotationScheduler handles scheduled credential rotation
type RotationScheduler struct {
	manager  *Manager
	policies map[string]*RotationPolicy
	stopChan chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

// NewRotationScheduler creates a new rotation scheduler
func NewRotationScheduler(manager *Manager) *RotationScheduler {
	return &RotationScheduler{
		manager:  manager,
		policies: make(map[string]*RotationPolicy),
		stopChan: make(chan struct{}),
	}
}

// AddPolicy adds a rotation policy for a provider
func (s *RotationScheduler) AddPolicy(providerID string, policy *RotationPolicy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.policies[providerID] = policy
}

// Start starts the rotation scheduler
func (s *RotationScheduler) Start(ctx context.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for providerID, policy := range s.policies {
		if !policy.Enabled {
			continue
		}

		s.wg.Add(1)
		go s.rotationWorker(ctx, providerID, policy)
	}
}

// Stop stops the rotation scheduler
func (s *RotationScheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}

// rotationWorker runs the rotation task for a specific provider
func (s *RotationScheduler) rotationWorker(ctx context.Context, providerID string, policy *RotationPolicy) {
	defer s.wg.Done()

	ticker := time.NewTicker(policy.RotationPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			if err := s.rotateCredentials(ctx, providerID); err != nil {
				// In production, this should be logged to an error tracking system
				fmt.Printf("Failed to rotate credentials for provider %s: %v\n", providerID, err)
			}
		}
	}
}

// rotateCredentials performs the credential rotation
func (s *RotationScheduler) rotateCredentials(ctx context.Context, providerID string) error {
	return s.manager.RotateAPIKey(ctx, providerID)
}

// GracefulRotationHandler provides graceful credential rotation with overlap
type GracefulRotationHandler struct {
	overlapDuration time.Duration
	onRotate        func(providerID string, oldKey, newKey string) error
}

// NewGracefulRotationHandler creates a handler that manages credential overlap
func NewGracefulRotationHandler(overlapDuration time.Duration, onRotate func(providerID string, oldKey, newKey string) error) *GracefulRotationHandler {
	return &GracefulRotationHandler{
		overlapDuration: overlapDuration,
		onRotate:        onRotate,
	}
}

// Handle implements the RotationHandler interface
func (h *GracefulRotationHandler) Handle(ctx context.Context, providerID string, oldCreds, newCreds *Credentials) error {
	// Notify external system of new credentials
	if h.onRotate != nil {
		if err := h.onRotate(providerID, oldCreds.APIKey, newCreds.APIKey); err != nil {
			return fmt.Errorf("rotation callback failed: %w", err)
		}
	}

	// In a production system, you might:
	// 1. Store both old and new credentials temporarily
	// 2. Schedule deletion of old credentials after overlap period
	// 3. Update all services to use new credentials
	// 4. Monitor for any services still using old credentials

	return nil
}

// RotationCallback is a convenience type for rotation callbacks
type RotationCallback func(ctx context.Context, providerID string, oldCreds, newCreds *Credentials) error

// WrapRotationCallback wraps a simple callback into a RotationHandler
func WrapRotationCallback(callback RotationCallback) RotationHandler {
	return RotationHandler(callback)
}

// ChainRotationHandlers chains multiple rotation handlers
func ChainRotationHandlers(handlers ...RotationHandler) RotationHandler {
	return func(ctx context.Context, providerID string, oldCreds, newCreds *Credentials) error {
		for _, handler := range handlers {
			if err := handler(ctx, providerID, oldCreds, newCreds); err != nil {
				return err
			}
		}
		return nil
	}
}
