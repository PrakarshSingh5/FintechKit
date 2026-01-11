package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
)

// OAuthConfig holds OAuth 2.0 configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	AuthURL      string
	TokenURL     string
	UsePKCE      bool // Use PKCE for enhanced security
}

// OAuthManager handles OAuth 2.0 flows
type OAuthManager struct {
	config *OAuthConfig
	oauth  *oauth2.Config
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager(config *OAuthConfig) *OAuthManager {
	oauth := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
	}

	return &OAuthManager{
		config: config,
		oauth:  oauth,
	}
}

// PKCEParams holds PKCE parameters
type PKCEParams struct {
	CodeVerifier  string
	CodeChallenge string
}

// GeneratePKCE generates PKCE parameters
func GeneratePKCE() (*PKCEParams, error) {
	// Generate random code verifier (43-128 characters)
	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}
	verifier := base64.RawURLEncoding.EncodeToString(verifierBytes)

	// For simplicity, using plain challenge (production should use S256)
	return &PKCEParams{
		CodeVerifier:  verifier,
		CodeChallenge: verifier,
	}, nil
}

// GetAuthorizationURL generates the OAuth authorization URL
func (m *OAuthManager) GetAuthorizationURL(state string, pkce *PKCEParams) string {
	opts := []oauth2.AuthCodeOption{}

	if m.config.UsePKCE && pkce != nil {
		opts = append(opts,
			oauth2.SetAuthURLParam("code_challenge", pkce.CodeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", "plain"),
		)
	}

	return m.oauth.AuthCodeURL(state, opts...)
}

// ExchangeCode exchanges authorization code for access token
func (m *OAuthManager) ExchangeCode(ctx context.Context, code string, pkce *PKCEParams) (*Credentials, error) {
	opts := []oauth2.AuthCodeOption{}

	if m.config.UsePKCE && pkce != nil {
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", pkce.CodeVerifier))
	}

	token, err := m.oauth.Exchange(ctx, code, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return &Credentials{
		Type:         CredentialTypeOAuth,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
		Metadata:     map[string]string{"token_type": token.TokenType},
	}, nil
}

// RefreshToken refreshes an OAuth access token
func (m *OAuthManager) RefreshToken(ctx context.Context, refreshToken string) (*Credentials, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := m.oauth.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &Credentials{
		Type:         CredentialTypeOAuth,
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiresAt:    newToken.Expiry,
		Metadata:     map[string]string{"token_type": newToken.TokenType},
	}, nil
}

// AutoRefreshMiddleware automatically refreshes tokens when needed
type AutoRefreshMiddleware struct {
	manager     *OAuthManager
	credManager *Manager
	providerID  string
}

// NewAutoRefreshMiddleware creates a middleware that auto-refreshes tokens
func NewAutoRefreshMiddleware(manager *OAuthManager, credManager *Manager, providerID string) *AutoRefreshMiddleware {
	return &AutoRefreshMiddleware{
		manager:     manager,
		credManager: credManager,
		providerID:  providerID,
	}
}

// GetValidToken ensures a valid access token, refreshing if needed
func (m *AutoRefreshMiddleware) GetValidToken(ctx context.Context) (string, error) {
	creds, err := m.credManager.GetCredentials(ctx, m.providerID)
	if err != nil {
		return "", err
	}

	// Check if refresh is needed
	if creds.NeedsRefresh() && creds.RefreshToken != "" {
		newCreds, err := m.manager.RefreshToken(ctx, creds.RefreshToken)
		if err != nil {
			return "", fmt.Errorf("failed to refresh token: %w", err)
		}

		// Update stored credentials
		if err := m.credManager.SetCredentials(ctx, m.providerID, newCreds); err != nil {
			return "", fmt.Errorf("failed to store refreshed credentials: %w", err)
		}

		return newCreds.AccessToken, nil
	}

	return creds.AccessToken, nil
}
