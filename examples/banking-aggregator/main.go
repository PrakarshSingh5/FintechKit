package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/PrakarshSingh5/fintechkit/pkg/auth"
	"github.com/PrakarshSingh5/fintechkit/pkg/middleware"
	"github.com/PrakarshSingh5/fintechkit/pkg/providers/plaid"
	"github.com/PrakarshSingh5/fintechkit/pkg/providers/truelayer"
)

func main() {
	// Initialize authentication manager
	authManager := auth.NewManager(auth.NewInMemoryStore())

	// Store credentials
	ctx := context.Background()

	// Plaid credentials
	authManager.SetCredentials(ctx, "plaid", &auth.Credentials{
		Type: auth.CredentialTypeAPIKey,
		Metadata: map[string]string{
			"client_id": "your_plaid_client_id",
			"secret":    "your_plaid_secret",
		},
	})

	// TrueLayer credentials (OAuth)
	authManager.SetCredentials(ctx, "truelayer", &auth.Credentials{
		Type:        auth.CredentialTypeOAuth,
		AccessToken: "your_truelayer_access_token",
	})

	// Initialize Fiber app
	app := fiber.New()

	// Global middleware
	app.Use(middleware.RecoveryMiddleware())
	app.Use(middleware.SecurityHeadersMiddleware())

	// API routes
	api := app.Group("/api/v1")
	api.Use(middleware.APIKeyMiddleware("your-api-key-here"))

	// Accounts endpoints
	api.Get("/accounts", getAccounts(authManager))
	api.Get("/accounts/:id/transactions", getTransactions(authManager))
	api.Get("/accounts/:id/balance", getBalance(authManager))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// Start server
	log.Println("Starting banking aggregator on :3002")
	log.Fatal(app.Listen(":3002"))
}

func getAccounts(authManager *auth.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		provider := c.Query("provider", "plaid") // Default to Plaid

		var accounts interface{}
		var err error

		switch provider {
		case "plaid":
			client, err := createPlaidClient(authManager)
			if err != nil {
				return err
			}
			accounts, err = client.GetAccounts(c.Context())

		case "truelayer":
			client, err := createTrueLayerClient(authManager)
			if err != nil {
				return err
			}
			accounts, err = client.GetAccounts(c.Context())

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Unsupported provider",
			})
		}

		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"provider": provider,
			"accounts": accounts,
		})
	}
}

func getTransactions(authManager *auth.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accountID := c.Params("id")
		provider := c.Query("provider", "plaid")

		// Parse date range
		startDate := time.Now().AddDate(0, -1, 0) // Last month
		endDate := time.Now()

		var transactions interface{}
		var err error

		switch provider {
		case "plaid":
			client, err := createPlaidClient(authManager)
			if err != nil {
				return err
			}
			transactions, err = client.GetTransactions(c.Context(), accountID, startDate, endDate)

		case "truelayer":
			client, err := createTrueLayerClient(authManager)
			if err != nil {
				return err
			}
			transactions, err = client.GetTransactions(c.Context(), accountID, startDate, endDate)

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Unsupported provider",
			})
		}

		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"provider":     provider,
			"account_id":   accountID,
			"transactions": transactions,
		})
	}
}

func getBalance(authManager *auth.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accountID := c.Params("id")
		provider := c.Query("provider", "plaid")

		var balance int64
		var currency string
		var err error

		switch provider {
		case "plaid":
			client, err := createPlaidClient(authManager)
			if err != nil {
				return err
			}
			balance, currency, err = client.GetBalance(c.Context(), accountID)

		case "truelayer":
			client, err := createTrueLayerClient(authManager)
			if err != nil {
				return err
			}
			balance, currency, err = client.GetBalance(c.Context(), accountID)

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Unsupported provider",
			})
		}

		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"provider":   provider,
			"account_id": accountID,
			"balance":    balance,
			"currency":   currency,
		})
	}
}

func createPlaidClient(authManager *auth.Manager) (*plaid.Client, error) {
	creds, err := authManager.GetCredentials(context.Background(), "plaid")
	if err != nil {
		return nil, err
	}

	return plaid.NewClient(&plaid.Config{
		ClientID: creds.Metadata["client_id"],
		Secret:   creds.Metadata["secret"],
		Env:      "sandbox",
	})
}

func createTrueLayerClient(authManager *auth.Manager) (*truelayer.Client, error) {
	creds, err := authManager.GetCredentials(context.Background(), "truelayer")
	if err != nil {
		return nil, err
	}

	client, err := truelayer.NewClient(&truelayer.Config{
		ClientID:     creds.Metadata["client_id"],
		ClientSecret: creds.Metadata["client_secret"],
		Env:          "sandbox",
	})
	if err != nil {
		return nil, err
	}

	// Set access token
	client.SetAccessToken(creds.AccessToken)
	return client, nil
}
