package tests

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/fintechkit/pkg/client"
	"github.com/yourusername/fintechkit/pkg/reliability"
)

// TestHelper provides utilities for fintech workflow testing
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// AssertPaymentSucceeded asserts that a payment succeeded
func (h *TestHelper) AssertPaymentSucceeded(payment *client.Payment) {
	if payment == nil {
		h.t.Fatal("Payment is nil")
	}

	if payment.Status != "succeeded" {
		h.t.Errorf("Expected payment status 'succeeded', got '%s'", payment.Status)
	}
}

// AssertPaymentAmount asserts the payment amount
func (h *TestHelper) AssertPaymentAmount(payment *client.Payment, expectedAmount int64, expectedCurrency string) {
	if payment.Amount != expectedAmount {
		h.t.Errorf("Expected amount %d, got %d", expectedAmount, payment.Amount)
	}

	if payment.Currency != expectedCurrency {
		h.t.Errorf("Expected currency %s, got %s", expectedCurrency, payment.Currency)
	}
}

// AssertAccountBalance asserts account balance
func (h *TestHelper) AssertAccountBalance(balance int64, expectedBalance int64) {
	if balance != expectedBalance {
		h.t.Errorf("Expected balance %d, got %d", expectedBalance, balance)
	}
}

// SimulatePaymentFlow simulates a complete payment flow
func (h *TestHelper) SimulatePaymentFlow(provider client.PaymentProvider, amount int64, currency string) (*client.Payment, *client.Refund) {
	ctx := context.Background()

	// Create payment
	payment, err := provider.CreatePayment(ctx, &client.PaymentRequest{
		Amount:   amount,
		Currency: currency,
	})
	if err != nil {
		h.t.Fatalf("Failed to create payment: %v", err)
	}

	h.AssertPaymentSucceeded(payment)
	h.AssertPaymentAmount(payment, amount, currency)

	// Simulate a refund
	refund, err := provider.RefundPayment(ctx, payment.ID, amount/2, "customer_request")
	if err != nil {
		h.t.Fatalf("Failed to refund payment: %v", err)
	}

	if refund.Amount != amount/2 {
		h.t.Errorf("Expected refund amount %d, got %d", amount/2, refund.Amount)
	}

	return payment, refund
}

// SimulateBankingFlow simulates account and transaction retrieval
func (h *TestHelper) SimulateBankingFlow(provider client.BankingProvider) {
	ctx := context.Background()

	// Get accounts
	accounts, err := provider.GetAccounts(ctx)
	if err != nil {
		h.t.Fatalf("Failed to get accounts: %v", err)
	}

	if len(accounts) == 0 {
		h.t.Fatal("No accounts returned")
	}

	// Get transactions for first account
	account := accounts[0]
	transactions, err := provider.GetTransactions(ctx, account.ID, time.Now().AddDate(0, -1, 0), time.Now())
	if err != nil {
		h.t.Fatalf("Failed to get transactions: %v", err)
	}

	// Assert transactions exist (optional based on test data)
	_ = transactions

	// Get balance
	balance, currency, err := provider.GetBalance(ctx, account.ID)
	if err != nil {
		h.t.Fatalf("Failed to get balance: %v", err)
	}

	if currency == "" {
		h.t.Error("Currency should not be empty")
	}

	_ = balance
}

// TestRetryBehavior tests retry logic
func (h *TestHelper) TestRetryBehavior(fn func() error, expectedRetries int) {
	attempts := 0
	policy := &reliability.RetryPolicy{
		MaxRetries:      expectedRetries,
		InitialInterval: 10 * time.Millisecond,
		MaxInterval:     100 * time.Millisecond,
		Multiplier:      2.0,
	}

	testFn := func() error {
		attempts++
		return fn()
	}

	reliability.WithRetry(context.Background(), policy, testFn)

	if attempts != expectedRetries+1 {
		h.t.Errorf("Expected %d attempts, got %d", expectedRetries+1, attempts)
	}
}

// AssertError asserts that an error occurred
func (h *TestHelper) AssertError(err error, expectedMessage string) {
	if err == nil {
		h.t.Fatal("Expected error, got nil")
	}

	if expectedMessage != "" && err.Error() != expectedMessage {
		h.t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

// AssertNoError asserts no error occurred
func (h *TestHelper) AssertNoError(err error) {
	if err != nil {
		h.t.Fatalf("Expected no error, got: %v", err)
	}
}
