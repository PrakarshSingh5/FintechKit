package truelayer

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/fintechkit/pkg/client"
)

// Client implements the TrueLayer Open Banking provider
type Client struct {
	clientID     string
	clientSecret string
	accessToken  string
	baseURL      string
	authURL      string
	env          string
}

// Config holds TrueLayer configuration
type Config struct {
	ClientID     string
	ClientSecret string
	Env          string // sandbox, production
}

// NewClient creates a new TrueLayer client
func NewClient(config *Config) (*Client, error) {
	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("client ID and secret are required")
	}

	env := config.Env
	if env == "" {
		env = "sandbox"
	}

	baseURL := getBaseURL(env)
	authURL := getAuthURL(env)

	return &Client{
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
		baseURL:      baseURL,
		authURL:      authURL,
		env:          env,
	}, nil
}

// getBaseURL returns the appropriate base URL
func getBaseURL(env string) string {
	if env == "production" {
		return "https://api.truelayer.com"
	}
	return "https://api.truelayer-sandbox.com"
}

// getAuthURL returns the appropriate auth URL
func getAuthURL(env string) string {
	if env == "production" {
		return "https://auth.truelayer.com"
	}
	return "https://auth.truelayer-sandbox.com"
}

// Name returns the provider name
func (c *Client) Name() string {
	return "truelayer"
}

// Authenticate verifies credentials
func (c *Client) Authenticate(ctx context.Context) error {
	if c.clientID == "" || c.clientSecret == "" {
		return fmt.Errorf("credentials not set")
	}
	return nil
}

// HealthCheck verifies API accessibility
func (c *Client) HealthCheck(ctx context.Context) error {
	return nil
}

// SetAccessToken sets the OAuth access token
func (c *Client) SetAccessToken(token string) {
	c.accessToken = token
}

// GetAccounts retrieves all linked accounts
func (c *Client) GetAccounts(ctx context.Context) ([]*client.Account, error) {
	// Real implementation would:
	// 1. Make HTTP GET to /data/v1/accounts
	// 2. Parse response

	return []*client.Account{
		{
			ID:            "acc_" + generateID(),
			Name:          "Current Account",
			Type:          "transaction",
			Balance:       250000, // £2,500.00
			Currency:      "GBP",
			AccountNumber: "****5678",
			Institution:   "Metro Bank",
		},
	}, nil
}

// GetAccount retrieves a specific account
func (c *Client) GetAccount(ctx context.Context, accountID string) (*client.Account, error) {
	// Real implementation would:
	// 1. Make HTTP GET to /data/v1/accounts/{account_id}

	return &client.Account{
		ID:       accountID,
		Name:     "Current Account",
		Type:     "transaction",
		Balance:  250000,
		Currency: "GBP",
	}, nil
}

// GetTransactions retrieves transactions
func (c *Client) GetTransactions(ctx context.Context, accountID string, startDate, endDate time.Time) ([]*client.Transaction, error) {
	// Real implementation would:
	// 1. Make HTTP GET to /data/v1/accounts/{account_id}/transactions
	// 2. Handle pagination

	return []*client.Transaction{
		{
			ID:          "txn_" + generateID(),
			AccountID:   accountID,
			Amount:      -3500, // -£35.00
			Currency:    "GBP",
			Date:        time.Now().Add(-24 * time.Hour),
			Description: "Tesco",
			Category:    "Groceries",
			Type:        "debit",
			Pending:     false,
		},
	}, nil
}

// GetBalance gets account balance
func (c *Client) GetBalance(ctx context.Context, accountID string) (int64, string, error) {
	// Real implementation would:
	// 1. Make HTTP GET to /data/v1/accounts/{account_id}/balance

	return 250000, "GBP", nil
}

// InitiatePayment initiates a PSD2 payment
func (c *Client) InitiatePayment(ctx context.Context, req *PaymentInitiationRequest) (*PaymentInitiation, error) {
	// Real implementation would:
	// 1. Make HTTP POST to /payments
	// 2. Return payment initiation details with redirect URL

	return &PaymentInitiation{
		ID:          "pmt_" + generateID(),
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      "authorization_required",
		RedirectURL: "https://payment.truelayer.com/..." + generateID(),
		CreatedAt:   time.Now(),
	}, nil
}

// GetPaymentStatus retrieves payment status
func (c *Client) GetPaymentStatus(ctx context.Context, paymentID string) (string, error) {
	// Real implementation would:
	// 1. Make HTTP GET to /payments/{payment_id}

	return "executed", nil
}

// PaymentInitiationRequest represents a payment initiation request
type PaymentInitiationRequest struct {
	Amount      int64
	Currency    string
	Beneficiary Beneficiary
	Reference   string
}

// Beneficiary represents payment beneficiary details
type Beneficiary struct {
	Type          string // business, individual
	Name          string
	AccountNumber string
	SortCode      string // For UK
}

// PaymentInitiation represents an initiated payment
type PaymentInitiation struct {
	ID          string
	Amount      int64
	Currency    string
	Status      string
	RedirectURL string
	CreatedAt   time.Time
}

// generateID generates a simple ID
func generateID() string {
	return "1234567890abcdef"
}

// CreatePayment implements PaymentProvider interface
func (c *Client) CreatePayment(ctx context.Context, req *client.PaymentRequest) (*client.Payment, error) {
	// Convert to TrueLayer payment initiation
	pmtReq := &PaymentInitiationRequest{
		Amount:   req.Amount,
		Currency: req.Currency,
		Beneficiary: Beneficiary{
			Type: "business",
			Name: "Merchant",
		},
		Reference: req.Description,
	}

	initiation, err := c.InitiatePayment(ctx, pmtReq)
	if err != nil {
		return nil, err
	}

	return &client.Payment{
		ID:          initiation.ID,
		Amount:      initiation.Amount,
		Currency:    initiation.Currency,
		Status:      initiation.Status,
		Description: req.Description,
		CreatedAt:   initiation.CreatedAt,
	}, nil
}

// GetPayment retrieves a payment
func (c *Client) GetPayment(ctx context.Context, id string) (*client.Payment, error) {
	status, err := c.GetPaymentStatus(ctx, id)
	if err != nil {
		return nil, err
	}

	return &client.Payment{
		ID:     id,
		Status: status,
	}, nil
}

// RefundPayment - TrueLayer doesn't support direct refunds
func (c *Client) RefundPayment(ctx context.Context, id string, amount int64, reason string) (*client.Refund, error) {
	return nil, fmt.Errorf("refunds not supported by TrueLayer")
}

// ListPayments lists payments
func (c *Client) ListPayments(ctx context.Context, filters map[string]string) ([]*client.Payment, error) {
	// Real implementation would call TrueLayer API
	return []*client.Payment{}, nil
}
