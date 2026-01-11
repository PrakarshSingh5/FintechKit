package webhook

import (
	"context"
	"encoding/json"
)

// Common event types across providers

// PaymentSucceededHandler handles payment.succeeded events
func PaymentSucceededHandler(onSuccess func(paymentID string, amount int64, currency string) error) Handler {
	return func(ctx context.Context, event *Event) error {
		var data struct {
			ID       string `json:"id"`
			Amount   int64  `json:"amount"`
			Currency string `json:"currency"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			return err
		}

		return onSuccess(data.ID, data.Amount, data.Currency)
	}
}

// PaymentFailedHandler handles payment.failed events
func PaymentFailedHandler(onFailure func(paymentID string, reason string) error) Handler {
	return func(ctx context.Context, event *Event) error {
		var data struct {
			ID     string `json:"id"`
			Reason string `json:"failure_reason"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			return err
		}

		return onFailure(data.ID, data.Reason)
	}
}

// RefundCreatedHandler handles refund.created events
func RefundCreatedHandler(onRefund func(refundID string, paymentID string, amount int64) error) Handler {
	return func(ctx context.Context, event *Event) error {
		var data struct {
			ID        string `json:"id"`
			PaymentID string `json:"payment_id"`
			Amount    int64  `json:"amount"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			return err
		}

		return onRefund(data.ID, data.PaymentID, data.Amount)
	}
}

// TransactionCreatedHandler handles transaction.created events
func TransactionCreatedHandler(onTransaction func(transactionID string, accountID string, amount int64) error) Handler {
	return func(ctx context.Context, event *Event) error {
		var data struct {
			ID        string `json:"id"`
			AccountID string `json:"account_id"`
			Amount    int64  `json:"amount"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			return err
		}

		return onTransaction(data.ID, data.AccountID, data.Amount)
	}
}

// AccountUpdatedHandler handles account.updated events
func AccountUpdatedHandler(onUpdate func(accountID string) error) Handler {
	return func(ctx context.Context, event *Event) error {
		var data struct {
			ID string `json:"id"`
		}

		if err := json.Unmarshal(event.Data, &data); err != nil {
			return err
		}

		return onUpdate(data.ID)
	}
}

// Common event type constants
const (
	// Stripe events
	EventStripePaymentIntentSucceeded = "payment_intent.succeeded"
	EventStripePaymentIntentFailed    = "payment_intent.payment_failed"
	EventStripeChargeRefunded         = "charge.refunded"
	EventStripeCustomerCreated        = "customer.created"

	// Plaid events
	EventPlaidItemError                 = "ITEM_ERROR"
	EventPlaidTransactionsReady         = "DEFAULT_UPDATE"
	EventPlaidWebhookUpdateAcknowledged = "WEBHOOK_UPDATE_ACKNOWLEDGED"

	// TrueLayer events
	EventTrueLayerPaymentExecuted   = "payment_executed"
	EventTrueLayerPaymentFailed     = "payment_failed"
	EventTrueLayerPaymentAuthorized = "payment_authorized"
)
