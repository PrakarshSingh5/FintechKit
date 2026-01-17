package tests

import (
	"context"
	"errors"
	"time"

	"github.com/PrakarshSingh5/fintechkit/pkg/client"
)

// MockPaymentProvider is a mock implementation of PaymentProvider for testing
type MockPaymentProvider struct {
	payments       map[string]*client.Payment
	shouldFail     bool
	failureMessage string
}

// NewMockPaymentProvider creates a new mock payment provider
func NewMockPaymentProvider() *MockPaymentProvider {
	return &MockPaymentProvider{
		payments: make(map[string]*client.Payment),
	}
}

// SetShouldFail configures the mock to fail on next operation
func (m *MockPaymentProvider) SetShouldFail(shouldFail bool, message string) {
	m.shouldFail = shouldFail
	m.failureMessage = message
}

// Name returns the provider name
func (m *MockPaymentProvider) Name() string {
	return "mock"
}

// Authenticate implements Provider
func (m *MockPaymentProvider) Authenticate(ctx context.Context) error {
	if m.shouldFail {
		return errors.New(m.failureMessage)
	}
	return nil
}

// HealthCheck implements Provider
func (m *MockPaymentProvider) HealthCheck(ctx context.Context) error {
	if m.shouldFail {
		return errors.New(m.failureMessage)
	}
	return nil
}

// CreatePayment creates a mock payment
func (m *MockPaymentProvider) CreatePayment(ctx context.Context, req *client.PaymentRequest) (*client.Payment, error) {
	if m.shouldFail {
		return nil, errors.New(m.failureMessage)
	}

	payment := &client.Payment{
		ID:          "mock_payment_" + generateID(),
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      "succeeded",
		Description: req.Description,
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.payments[payment.ID] = payment
	return payment, nil
}

// GetPayment retrieves a mock payment
func (m *MockPaymentProvider) GetPayment(ctx context.Context, id string) (*client.Payment, error) {
	if m.shouldFail {
		return nil, errors.New(m.failureMessage)
	}

	payment, ok := m.payments[id]
	if !ok {
		return nil, errors.New("payment not found")
	}

	return payment, nil
}

// RefundPayment creates a mock refund
func (m *MockPaymentProvider) RefundPayment(ctx context.Context, id string, amount int64, reason string) (*client.Refund, error) {
	if m.shouldFail {
		return nil, errors.New(m.failureMessage)
	}

	payment, ok := m.payments[id]
	if !ok {
		return nil, errors.New("payment not found")
	}

	return &client.Refund{
		ID:        "mock_refund_" + generateID(),
		PaymentID: payment.ID,
		Amount:    amount,
		Status:    "succeeded",
		Reason:    reason,
		CreatedAt: time.Now(),
	}, nil
}

// ListPayments lists mock payments
func (m *MockPaymentProvider) ListPayments(ctx context.Context, filters map[string]string) ([]*client.Payment, error) {
	if m.shouldFail {
		return nil, errors.New(m.failureMessage)
	}

	payments := make([]*client.Payment, 0, len(m.payments))
	for _, payment := range m.payments {
		payments = append(payments, payment)
	}

	return payments, nil
}

// MockBankingProvider is a mock implementation of BankingProvider
type MockBankingProvider struct {
	accounts       map[string]*client.Account
	transactions   map[string][]*client.Transaction
	shouldFail     bool
	failureMessage string
}

// NewMockBankingProvider creates a new mock banking provider
func NewMockBankingProvider() *MockBankingProvider {
	mock := &MockBankingProvider{
		accounts:     make(map[string]*client.Account),
		transactions: make(map[string][]*client.Transaction),
	}

	// Add some default test data
	account := &client.Account{
		ID:       "acc_test_123",
		Name:     "Test Checking",
		Type:     "checking",
		Balance:  100000, // $1,000.00
		Currency: "USD",
	}
	mock.accounts[account.ID] = account

	mock.transactions[account.ID] = []*client.Transaction{
		{
			ID:          "txn_1",
			AccountID:   account.ID,
			Amount:      -5000,
			Currency:    "USD",
			Date:        time.Now().Add(-24 * time.Hour),
			Description: "Coffee Shop",
			Type:        "debit",
		},
	}

	return mock
}

// Name returns the provider name
func (m *MockBankingProvider) Name() string {
	return "mock"
}

// Authenticate implements Provider
func (m *MockBankingProvider) Authenticate(ctx context.Context) error {
	if m.shouldFail {
		return errors.New(m.failureMessage)
	}
	return nil
}

// HealthCheck implements Provider
func (m *MockBankingProvider) HealthCheck(ctx context.Context) error {
	if m.shouldFail {
		return errors.New(m.failureMessage)
	}
	return nil
}

// GetAccounts returns mock accounts
func (m *MockBankingProvider) GetAccounts(ctx context.Context) ([]*client.Account, error) {
	if m.shouldFail {
		return nil, errors.New(m.failureMessage)
	}

	accounts := make([]*client.Account, 0, len(m.accounts))
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// GetAccount returns a specific mock account
func (m *MockBankingProvider) GetAccount(ctx context.Context, accountID string) (*client.Account, error) {
	if m.shouldFail {
		return nil, errors.New(m.failureMessage)
	}

	account, ok := m.accounts[accountID]
	if !ok {
		return nil, errors.New("account not found")
	}

	return account, nil
}

// GetTransactions returns mock transactions
func (m *MockBankingProvider) GetTransactions(ctx context.Context, accountID string, startDate, endDate time.Time) ([]*client.Transaction, error) {
	if m.shouldFail {
		return nil, errors.New(m.failureMessage)
	}

	transactions, ok := m.transactions[accountID]
	if !ok {
		return []*client.Transaction{}, nil
	}

	return transactions, nil
}

// GetBalance returns mock balance
func (m *MockBankingProvider) GetBalance(ctx context.Context, accountID string) (int64, string, error) {
	if m.shouldFail {
		return 0, "", errors.New(m.failureMessage)
	}

	account, ok := m.accounts[accountID]
	if !ok {
		return 0, "", errors.New("account not found")
	}

	return account.Balance, account.Currency, nil
}

// SetShouldFail configures the mock to fail
func (m *MockBankingProvider) SetShouldFail(shouldFail bool, message string) {
	m.shouldFail = shouldFail
	m.failureMessage = message
}

// Helper function
func generateID() string {
	return "1234567890"
}
