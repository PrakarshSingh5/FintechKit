package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

// CredentialType represents the type of credentials being managed
type CredentialType string

const (
	CredentialTypeAPIKey CredentialType = "api_key"
	CredentialTypeOAuth  CredentialType = "oauth"
	CredentialTypeBearer CredentialType = "bearer"
)

// Credentials represents authentication credentials
type Credentials struct {
	Type         CredentialType
	APIKey       string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	Metadata     map[string]string
}

// CredentialStore interface for storing and retrieving credentials
type CredentialStore interface {
	Get(ctx context.Context, providerID string) (*Credentials, error)
	Set(ctx context.Context, providerID string, creds *Credentials) error
	Delete(ctx context.Context, providerID string) error
}

// InMemoryStore is a simple in-memory credential store (not for production)
type InMemoryStore struct {
	mu    sync.RWMutex
	store map[string]*Credentials
}

// NewInMemoryStore creates a new in-memory credential store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string]*Credentials),
	}
}

// Get retrieves credentials for a provider
func (s *InMemoryStore) Get(ctx context.Context, providerID string) (*Credentials, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	creds, ok := s.store[providerID]
	if !ok {
		return nil, errors.New("credentials not found")
	}
	return creds, nil
}

// Set stores credentials for a provider
func (s *InMemoryStore) Set(ctx context.Context, providerID string, creds *Credentials) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[providerID] = creds
	return nil
}

// Delete removes credentials for a provider
func (s *InMemoryStore) Delete(ctx context.Context, providerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, providerID)
	return nil
}

// Manager handles credential management for multiple providers
type Manager struct {
	store           CredentialStore
	rotationEnabled bool
	rotationHandlers map[string]RotationHandler
	mu              sync.RWMutex
}

// RotationHandler is called when credentials are rotated
type RotationHandler func(ctx context.Context, providerID string, oldCreds, newCreds *Credentials) error

// NewManager creates a new credential manager
func NewManager(store CredentialStore) *Manager {
	return &Manager{
		store:            store,
		rotationEnabled:  false,
		rotationHandlers: make(map[string]RotationHandler),
	}
}

// GetCredentials retrieves credentials for a provider
func (m *Manager) GetCredentials(ctx context.Context, providerID string) (*Credentials, error) {
	creds, err := m.store.Get(ctx, providerID)
	if err != nil {
		return nil, err
	}

	// Check if OAuth token needs refresh
	if creds.Type == CredentialTypeOAuth && !creds.ExpiresAt.IsZero() && time.Now().After(creds.ExpiresAt) {
		return nil, errors.New("access token expired, refresh required")
	}

	return creds, nil
}

// SetCredentials stores credentials for a provider
func (m *Manager) SetCredentials(ctx context.Context, providerID string, creds *Credentials) error {
	return m.store.Set(ctx, providerID, creds)
}

// RegisterRotationHandler registers a handler for credential rotation
func (m *Manager) RegisterRotationHandler(providerID string, handler RotationHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rotationHandlers[providerID] = handler
}

// RotateAPIKey generates a new API key and calls the rotation handler
func (m *Manager) RotateAPIKey(ctx context.Context, providerID string) error {
	oldCreds, err := m.store.Get(ctx, providerID)
	if err != nil {
		return err
	}

	if oldCreds.Type != CredentialTypeAPIKey {
		return errors.New("can only rotate API key credentials")
	}

	// Generate new API key
	newKey, err := generateAPIKey()
	if err != nil {
		return err
	}

	newCreds := &Credentials{
		Type:     CredentialTypeAPIKey,
		APIKey:   newKey,
		Metadata: oldCreds.Metadata,
	}

	// Call rotation handler if registered
	m.mu.RLock()
	handler, exists := m.rotationHandlers[providerID]
	m.mu.RUnlock()

	if exists {
		if err := handler(ctx, providerID, oldCreds, newCreds); err != nil {
			return err
		}
	}

	// Store new credentials
	return m.store.Set(ctx, providerID, newCreds)
}

// generateAPIKey generates a random API key
func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// IsExpired checks if credentials are expired
func (c *Credentials) IsExpired() bool {
	if c.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(c.ExpiresAt)
}

// NeedsRefresh checks if credentials need to be refreshed (5 min buffer)
func (c *Credentials) NeedsRefresh() bool {
	if c.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().Add(5 * time.Minute).After(c.ExpiresAt)
}
