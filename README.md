# FinTechKit 🚀

A production-ready Go framework for fintech API integration with built-in authentication, reliability patterns, compliance helpers, and webhook management.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 🎯 Features

- **Unified API Interfaces** - Common interfaces for Stripe, Razorpay, Plaid, TrueLayer, and CoinGecko
- **Authentication Management** - OAuth 2.0, API key management, and automatic token rotation
- **Reliability Patterns** - Retry logic, rate limiting, and circuit breakers

- **Webhook Management** - Signature verification, event routing, and idempotency
- **Fiber Integration** - Production-ready middleware for Fiber web framework
- **Testing Utilities** - Mock providers and test helpers for fintech workflows

## 📦 Installation

```bash
go get github.com/PrakarshSingh5/fintechkit
```

## 🚀 Quick Start

### Payment Processing with Stripe

```go
package main

import (
    "context"
    "log"

    "github.com/PrakarshSingh5/fintechkit/pkg/auth"
    "github.com/PrakarshSingh5/fintechkit/pkg/client"
    "github.com/PrakarshSingh5/fintechkit/pkg/providers/stripe"
    "github.com/PrakarshSingh5/fintechkit/pkg/reliability"
)

func main() {
    // Initialize authentication
    authManager := auth.NewManager(auth.NewInMemoryStore())

    ctx := context.Background()
    authManager.SetCredentials(ctx, "stripe", &auth.Credentials{
        Type:   auth.CredentialTypeAPIKey,
        APIKey: "sk_test_your_key",
    })

    // Create Stripe client
    stripeClient, _ := stripe.NewClient(&stripe.Config{
        APIKey: "sk_test_your_key",
    })

    // Create payment with automatic retry
    config := &client.ProviderConfig{
        Name:        "stripe",
        RetryPolicy: reliability.StripeRetryPolicy(),
    }

    payment, err := stripeClient.CreatePayment(ctx, &client.PaymentRequest{
        Amount:   10000, // $100.00 in cents
        Currency: "usd",
    })

    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Payment created: %s", payment.ID)
}
```

### Banking with Plaid

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/PrakarshSingh5/fintechkit/pkg/providers/plaid"
)

func main() {
    client, _ := plaid.NewClient(&plaid.Config{
        ClientID: "your_client_id",
        Secret:   "your_secret",
        Env:      "sandbox",
    })

    ctx := context.Background()

    // Get accounts
    accounts, _ := client.GetAccounts(ctx)

    // Get transactions
    transactions, _ := client.GetTransactions(
        ctx,
        accounts[0].ID,
        time.Now().AddDate(0, -1, 0),
        time.Now(),
    )

    log.Printf("Found %d transactions", len(transactions))
}
```

### Webhook Handling

```go
package main

import (
    "context"
    "log"

    "github.com/gofiber/fiber/v2"
    "github.com/PrakarshSingh5/fintechkit/pkg/webhook"
)

func main() {
    receiver := webhook.NewReceiver()
    router := webhook.NewRouter()

    // Register handler for payment success
    router.Register("stripe", webhook.EventStripePaymentIntentSucceeded,
        webhook.PaymentSucceededHandler(func(paymentID string, amount int64, currency string) error {
            log.Printf("Payment succeeded: %s - %d %s", paymentID, amount, currency)
            return nil
        }),
    )

    app := fiber.New()

    app.Post("/webhooks/stripe", func(c *fiber.Ctx) error {
        payload := c.Body()
        signature := c.Get("Stripe-Signature")

        err := receiver.ProcessEvent(context.Background(), "stripe", payload, signature)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": err.Error()})
        }

        return c.JSON(fiber.Map{"status": "received"})
    })

    app.Listen(":3000")
}
```

## 🏗️ Architecture

```
fintechkit/
├── pkg/
│   ├── auth/           # Authentication & credentials management
│   ├── client/         # Unified provider interfaces
│   ├── providers/      # API integrations
│   │   ├── stripe/     # Stripe payment processing
│   │   ├── razorpay/   # Razorpay payment gateway (India)
│   │   ├── plaid/      # Plaid banking data
│   │   ├── truelayer/  # TrueLayer Open Banking
│   │   └── coingecko/  # CoinGecko crypto data
│   ├── reliability/    # Retry, rate limiting, circuit breakers

│   ├── webhook/        # Webhook management
│   └── middleware/     # Fiber middleware
├── examples/           # Example applications
└── tests/             # Testing utilities
```

## 🛡️ Reliability Features

### Retry Logic with Exponential Backoff

```go
policy := &reliability.RetryPolicy{
    MaxRetries:      3,
    InitialInterval: 1 * time.Second,
    MaxInterval:     30 * time.Second,
    Multiplier:      2.0,
    RandomizeJitter: true,
}

err := reliability.WithRetry(ctx, policy, func() error {
    return doSomething()
})
```

### Rate Limiting

```go
limiter := reliability.NewRateLimiter(&reliability.RateLimitConfig{
    RequestsPerSecond: 10.0,
    Burst:             20,
    WaitTimeout:       5 * time.Second,
})

limiter.Wait(ctx)
```

### Circuit Breaker

```go
breaker := reliability.NewCircuitBreaker("stripe", &reliability.CircuitBreakerConfig{
    FailureThreshold: 5,
    Timeout:          30 * time.Second,
})

result, err := breaker.Execute(func() (interface{}, error) {
    return callExternalAPI()
})
```

## 🧪 Testing

```go
import "github.com/PrakarshSingh5/fintechkit/tests"

func TestPaymentFlow(t *testing.T) {
    helper := tests.NewTestHelper(t)
    provider := tests.NewMockPaymentProvider()

    payment, refund := helper.SimulatePaymentFlow(provider, 10000, "usd")

    helper.AssertPaymentSucceeded(payment)
    helper.AssertNoError(nil)
}
```

## 📚 Examples

See the `/examples` directory for complete working examples:

- **payment-flow**: Payment processing with Stripe
- **razorpay-integration**: Razorpay payment gateway integration (India)
- **banking-aggregator**: Multi-provider banking data aggregation
- **webhook-server**: Webhook receiver with event routing

Run examples:

```bash
cd examples/payment-flow
go run main.go
```

## 🤝 Contributing

This project is perfect for open-source contribution to [FINOS](https://www.finos.org/) (Fintech Open Source Foundation).

## 📄 License

MIT License - see [LICENSE](LICENSE) for details

## 🙏 Acknowledgments

Built for fintech startups, payment platforms, and investment apps requiring robust API integrations.

---

**Made with ❤️ for the fintech community**
