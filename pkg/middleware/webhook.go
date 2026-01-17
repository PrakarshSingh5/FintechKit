package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/PrakarshSingh5/fintechkit/pkg/webhook"
)

// WebhookMiddleware creates middleware for webhook validation
func WebhookMiddleware(receiver *webhook.Receiver, provider string, getSignature func(*fiber.Ctx) string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get webhook payload
		payload := c.Body()

		// Get signature from appropriate header
		signature := getSignature(c)

		// Store for downstream handlers
		c.Locals("webhook_payload", payload)
		c.Locals("webhook_signature", signature)
		c.Locals("webhook_provider", provider)

		return c.Next()
	}
}

// StripeWebhookMiddleware validates Stripe webhooks
func StripeWebhookMiddleware() fiber.Handler {
	return WebhookMiddleware(
		nil,
		"stripe",
		func(c *fiber.Ctx) string {
			return c.Get("Stripe-Signature")
		},
	)
}

// PlaidWebhookMiddleware validates Plaid webhooks
func PlaidWebhookMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Plaid doesn't use signature verification
		// Just store the payload
		c.Locals("webhook_payload", c.Body())
		c.Locals("webhook_provider", "plaid")

		return c.Next()
	}
}

// IdempotencyMiddleware ensures webhook idempotency
func IdempotencyMiddleware(tracker *webhook.IdempotencyTracker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get idempotency key from header
		key := c.Get("Idempotency-Key")

		if key != "" {
			if tracker.IsProcessed(key) {
				// Already processed, return success
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status": "already_processed",
				})
			}

			// Mark as processed after successful processing
			defer tracker.MarkProcessed(key)
		}

		return c.Next()
	}
}

// RequestLoggingMiddleware logs all requests
func RequestLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request details
		duration := time.Since(start)

		// In production, use structured logging
		// logger.Info("request",
		// 	"method", c.Method(),
		// 	"path", c.Path(),
		// 	"status", c.Response().StatusCode(),
		// 	"duration", duration,
		// 	"ip", c.IP(),
		// )

		_ = duration // Suppress unused variable warning

		return err
	}
}
