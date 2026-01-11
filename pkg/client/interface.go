package client

import (
	"context"
	"time"
)

// Provider is the base interface that all API providers must implement
type Provider interface {
	// Name returns the provider's identifier
	Name() string

	// Authenticate performs initial authentication with the provider
	Authenticate(ctx context.Context) error

	// HealthCheck verifies the provider connection is working
	HealthCheck(ctx context.Context) error
}

// Payment represents a payment transaction
type Payment struct {
	ID          string
	Amount      int64 // Amount in smallest currency unit (cents, pence, etc.)
	Currency    string
	Status      string
	Description string
	Metadata    map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PaymentRequest represents a request to create a payment
type PaymentRequest struct {
	Amount        int64
	Currency      string
	Description   string
	Metadata      map[string]string
	CustomerID    string
	PaymentMethod string
}

// Refund represents a refund transaction
type Refund struct {
	ID        string
	PaymentID string
	Amount    int64
	Currency  string
	Status    string
	Reason    string
	CreatedAt time.Time
}

// PaymentProvider defines the interface for payment processing
type PaymentProvider interface {
	Provider

	// CreatePayment initiates a new payment
	CreatePayment(ctx context.Context, req *PaymentRequest) (*Payment, error)

	// GetPayment retrieves payment details
	GetPayment(ctx context.Context, id string) (*Payment, error)

	// RefundPayment refunds a payment (full or partial)
	RefundPayment(ctx context.Context, id string, amount int64, reason string) (*Refund, error)

	// ListPayments lists payments with optional filters
	ListPayments(ctx context.Context, filters map[string]string) ([]*Payment, error)
}

// Account represents a bank account
type Account struct {
	ID            string
	Name          string
	Type          string // checking, savings, credit, etc.
	Balance       int64
	Currency      string
	AccountNumber string // Masked for security
	Institution   string
	Metadata      map[string]string
}

// Transaction represents a bank transaction
type Transaction struct {
	ID          string
	AccountID   string
	Amount      int64
	Currency    string
	Date        time.Time
	Description string
	Category    string
	Type        string // debit, credit
	Pending     bool
	Metadata    map[string]string
}

// BankingProvider defines the interface for banking/account aggregation
type BankingProvider interface {
	Provider

	// GetAccounts retrieves all linked accounts
	GetAccounts(ctx context.Context) ([]*Account, error)

	// GetAccount retrieves a specific account
	GetAccount(ctx context.Context, accountID string) (*Account, error)

	// GetTransactions retrieves transactions for an account
	GetTransactions(ctx context.Context, accountID string, startDate, endDate time.Time) ([]*Transaction, error)

	// GetBalance gets the current balance for an account
	GetBalance(ctx context.Context, accountID string) (int64, string, error)
}

// Price represents cryptocurrency price data
type Price struct {
	CoinID    string
	Currency  string
	Price     float64
	Change24h float64
	Timestamp time.Time
}

// MarketData represents detailed market data for a cryptocurrency
type MarketData struct {
	CoinID             string
	Symbol             string
	Name               string
	CurrentPrice       float64
	MarketCap          float64
	Volume24h          float64
	PriceChange24h     float64
	PriceChangePercent float64
	High24h            float64
	Low24h             float64
	CirculatingSupply  float64
	TotalSupply        float64
	AllTimeHigh        float64
	AllTimeHighDate    time.Time
	LastUpdated        time.Time
}

// CryptoProvider defines the interface for cryptocurrency data
type CryptoProvider interface {
	Provider

	// GetPrice retrieves the current price for a cryptocurrency
	GetPrice(ctx context.Context, coinID string, currency string) (*Price, error)

	// GetPrices retrieves prices for multiple cryptocurrencies
	GetPrices(ctx context.Context, coinIDs []string, currency string) ([]*Price, error)

	// GetMarketData retrieves detailed market data
	GetMarketData(ctx context.Context, coinID string) (*MarketData, error)

	// GetHistoricalPrices retrieves historical price data
	GetHistoricalPrices(ctx context.Context, coinID string, currency string, days int) ([]*Price, error)
}

// Identity represents user identity information
type Identity struct {
	Name    string
	Email   string
	Phone   string
	Address Address
	DOB     time.Time
}

// Address represents a physical address
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// IdentityProvider defines the interface for identity verification
type IdentityProvider interface {
	Provider

	// VerifyIdentity verifies a user's identity
	VerifyIdentity(ctx context.Context, userID string) (*Identity, error)

	// GetIdentity retrieves verified identity information
	GetIdentity(ctx context.Context, userID string) (*Identity, error)
}
