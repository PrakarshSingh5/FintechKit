package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/fintechkit/pkg/middleware"
	"github.com/yourusername/fintechkit/pkg/webhook"
)

func main() {
	// Initialize webhook receiver
	receiver := webhook.NewReceiver()

	// Register Stripe webhook verifier
	stripeSecret := "whsec_your_stripe_webhook_secret"
	receiver.RegisterVerifier("stripe", webhook.NewStripeVerifier(stripeSecret))

	// Initialize webhook router
	router := webhook.NewRouter()

	// Register handlers for different event types
	router.Register("stripe", webhook.EventStripePaymentIntentSucceeded, webhook.PaymentSucceededHandler(
		func(paymentID string, amount int64, currency string) error {
			log.Printf("Payment succeeded: %s, Amount: %d %s", paymentID, amount, currency)
			// Handle successful payment (e.g., fulfill order, send confirmation)
			return nil
		},
	))

	router.Register("stripe", webhook.EventStripePaymentIntentFailed, webhook.PaymentFailedHandler(
		func(paymentID string, reason string) error {
			log.Printf("Payment failed: %s, Reason: %s", paymentID, reason)
			// Handle failed payment (e.g., notify customer)
			return nil
		},
	))

	router.Register("stripe", webhook.EventStripeChargeRefunded, webhook.RefundCreatedHandler(
		func(refundID string, paymentID string, amount int64) error {
			log.Printf("Refund created: %s for payment %s, Amount: %d", refundID, paymentID, amount)
			// Handle refund (e.g., update order status)
			return nil
		},
	))

	// Initialize Fiber app
	app := fiber.New()

	// Global middleware
	app.Use(middleware.RecoveryMiddleware())
	app.Use(middleware.RequestLoggingMiddleware())

	// Webhook endpoint
	app.Post("/webhooks/stripe", handleStripeWebhook(receiver, router))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// Start server
	log.Println("Starting webhook server on :3001")
	log.Fatal(app.Listen(":3001"))
}

func handleStripeWebhook(receiver *webhook.Receiver, router *webhook.Router) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := c.Body()
		signature := c.Get("Stripe-Signature")

		// Process webhook event
		err := receiver.ProcessEvent(context.Background(), "stripe", payload, signature)
		if err != nil {
			log.Printf("Webhook processing error: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Webhook processing failed",
			})
		}

		// Route event to handlers (receiver already validated)
		// In production, you might want to handle this asynchronously

		return c.JSON(fiber.Map{
			"status": "received",
		})
	}
}
