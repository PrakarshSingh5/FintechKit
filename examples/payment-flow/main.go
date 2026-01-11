package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/fintechkit/pkg/auth"
	"github.com/yourusername/fintechkit/pkg/client"
	"github.com/yourusername/fintechkit/pkg/middleware"
	"github.com/yourusername/fintechkit/pkg/providers/stripe"
	"github.com/yourusername/fintechkit/pkg/reliability"
)

func main() {
	// Initialize authentication manager
	authManager := auth.NewManager(auth.NewInMemoryStore())

	// Store Stripe credentials
	ctx := context.Background()
	err := authManager.SetCredentials(ctx, "stripe", &auth.Credentials{
		Type:   auth.CredentialTypeAPIKey,
		APIKey: "sk_test_your_stripe_key_here",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Global middleware
	app.Use(middleware.RecoveryMiddleware())
	app.Use(middleware.RequestLoggingMiddleware())
	app.Use(middleware.SecurityHeadersMiddleware())
	app.Use(middleware.CORSConfig())

	// API routes with authentication
	api := app.Group("/api/v1")
	api.Use(middleware.APIKeyMiddleware("your-api-key-here"))

	// Payment endpoints
	payments := api.Group("/payments")
	payments.Use(middleware.RateLimitMiddleware(10, 20)) // 10 req/sec, burst 20

	payments.Post("/", createPayment(authManager))
	payments.Get("/:id", getPayment(authManager))
	payments.Post("/:id/refund", refundPayment(authManager))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// Start server
	log.Println("Starting payment flow server on :3000")
	log.Fatal(app.Listen(":3000"))
}

func createPayment(authManager *auth.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse request
		var req client.PaymentRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Create Stripe client
		creds, err := authManager.GetCredentials(c.Context(), "stripe")
		if err != nil {
			return err
		}

		stripeClient, err := stripe.NewClient(&stripe.Config{
			APIKey: creds.APIKey,
		})
		if err != nil {
			return err
		}

		// Apply reliability features
		config := &client.ProviderConfig{
			Name:        "stripe",
			Credentials: creds,
			RetryPolicy: reliability.StripeRetryPolicy(),
		}

		factory := client.NewFactory(authManager)
		factory.Register("stripe", func(cfg *client.ProviderConfig) (client.Provider, error) {
			return stripeClient, nil
		})

		provider, err := factory.Create(c.Context(), config)
		if err != nil {
			return err
		}

		// Create payment
		paymentProvider, ok := provider.(client.PaymentProvider)
		if !ok {
			return fmt.Errorf("provider does not support payments")
		}

		payment, err := paymentProvider.CreatePayment(c.Context(), &req)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(payment)
	}
}

func getPayment(authManager *auth.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		paymentID := c.Params("id")

		// Create Stripe client
		creds, err := authManager.GetCredentials(c.Context(), "stripe")
		if err != nil {
			return err
		}

		stripeClient, err := stripe.NewClient(&stripe.Config{
			APIKey: creds.APIKey,
		})
		if err != nil {
			return err
		}

		payment, err := stripeClient.GetPayment(c.Context(), paymentID)
		if err != nil {
			return err
		}

		return c.JSON(payment)
	}
}

func refundPayment(authManager *auth.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		paymentID := c.Params("id")

		var req struct {
			Amount int64  `json:"amount"`
			Reason string `json:"reason"`
		}
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		// Create Stripe client
		creds, err := authManager.GetCredentials(c.Context(), "stripe")
		if err != nil {
			return err
		}

		stripeClient, err := stripe.NewClient(&stripe.Config{
			APIKey: creds.APIKey,
		})
		if err != nil {
			return err
		}

		refund, err := stripeClient.RefundPayment(c.Context(), paymentID, req.Amount, req.Reason)
		if err != nil {
			return err
		}

		return c.JSON(refund)
	}
}
