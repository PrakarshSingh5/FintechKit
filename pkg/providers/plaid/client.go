package plaid

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/fintechkit/pkg/client"
)

// Client implements the Plaid banking provider
type Client struct {
	clientID string
	secret   string
	baseURL  string
	env      string
}

// Config holds Plaid configuration
type Config struct {
	ClientID string
	Secret   string
	Env      string // sandbox, development, production
}

// NewClient creates a new Plaid client
func NewClient(config *Config) (*Client, error) {
	if config.ClientID == "" || config.Secret == "" {
		return nil, fmt.Errorf("client ID and secret are required")
	}

	env := config.Env
	if env == "" {
		env = "sandbox"
	}

	baseURL := getBaseURL(env)

	return &Client{
		clientID: config.ClientID,
		secret:   config.Secret,
		baseURL:  baseURL,
		env:      env,
	}, nil
}

// getBaseURL returns the appropriate base URL for the environment
func getBaseURL(env string) string {
	switch env {
	case "production":
		return "https://production.plaid.com"
	case "development":
		return "https://development.plaid.com"
	default:
		return "https://sandbox.plaid.com"
	}
}

// Name returns the provider name
func (c *Client) Name() string {
	return "plaid"
}

// Authenticate verifies credentials are valid
func (c *Client) Authenticate(ctx context.Context) error {
	if c.clientID == "" || c.secret == "" {
		return fmt.Errorf("credentials not set")
	}
	return nil
}

// HealthCheck verifies Plaid API is accessible
func (c *Client) HealthCheck(ctx context.Context) error {
	// Real implementation would call a test endpoint
	return nil
}

// CreateLinkToken creates a link token for Plaid Link
func (c *Client) CreateLinkToken(ctx context.Context, userID string, products []string) (string, error) {
	// Real implementation would:
	// 1. Make HTTP POST to /link/token/create
	// 2. Return the link token

	return "link-sandbox-" + generateID(), nil
}

// ExchangePublicToken exchanges a public token for an access token
func (c *Client) ExchangePublicToken(ctx context.Context, publicToken string) (string, string, error) {
	// Real implementation would:
	// 1. Make HTTP POST to /item/public_token/exchange
	// 2. Return access_token and item_id

	accessToken := "access-sandbox-" + generateID()
	itemID := "item-" + generateID()
	return accessToken, itemID, nil
}

// GetAccounts retrieves all linked accounts
func (c *Client) GetAccounts(ctx context.Context) ([]*client.Account, error) {
	// Real implementation would:
	// 1. Make HTTP POST to /accounts/get
	// 2. Parse response

	return []*client.Account{
		{
			ID:            "acc_" + generateID(),
			Name:          "Checking Account",
			Type:          "depository",
			Balance:       150000, // $1,500.00
			Currency:      "USD",
			AccountNumber: "****1234",
			Institution:   "Chase",
		},
	}, nil
}

// GetAccount retrieves a specific account
func (c *Client) GetAccount(ctx context.Context, accountID string) (*client.Account, error) {
	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return nil, fmt.Errorf("account not found")
}

// GetTransactions retrieves transactions for an account
func (c *Client) GetTransactions(ctx context.Context, accountID string, startDate, endDate time.Time) ([]*client.Transaction, error) {
	// Real implementation would:
	// 1. Make HTTP POST to /transactions/get
	// 2. Parse response with pagination

	return []*client.Transaction{
		{
			ID:          "txn_" + generateID(),
			AccountID:   accountID,
			Amount:      -4500, // -$45.00 (negative for debit)
			Currency:    "USD",
			Date:        time.Now().Add(-24 * time.Hour),
			Description: "Coffee Shop",
			Category:    "Food and Drink",
			Type:        "debit",
			Pending:     false,
		},
		{
			ID:          "txn_" + generateID(),
			AccountID:   accountID,
			Amount:      50000, // $500.00 (positive for credit)
			Currency:    "USD",
			Date:        time.Now().Add(-48 * time.Hour),
			Description: "Payroll Deposit",
			Category:    "Income",
			Type:        "credit",
			Pending:     false,
		},
	}, nil
}

// GetBalance gets the current balance for an account
func (c *Client) GetBalance(ctx context.Context, accountID string) (int64, string, error) {
	// Real implementation would:
	// 1. Make HTTP POST to /accounts/balance/get
	// 2. Return balance and currency

	return 150000, "USD", nil
}

// GetIdentity retrieves identity information
func (c *Client) GetIdentity(ctx context.Context, accessToken string) (*client.Identity, error) {
	// Real implementation would:
	// 1. Make HTTP POST to /identity/get
	// 2. Parse response

	return &client.Identity{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "+1-555-123-4567",
		Address: client.Address{
			Street:     "123 Main St",
			City:       "San Francisco",
			State:      "CA",
			PostalCode: "94102",
			Country:    "US",
		},
	}, nil
}

// generateID generates a simple ID for demonstration
func generateID() string {
	return "1234567890abcdef"
}
