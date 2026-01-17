package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/fintechkit/pkg/webhook"
)

func main() {
	// Step 1: Initialize webhook receiver and router
	receiver := webhook.NewReceiver()
	router := webhook.NewRouter()

	// Step 2: Register handler for Razorpay payment success
	router.Register("razorpay", "payment.captured",
		webhook.PaymentSucceededHandler(func(paymentID string, amount int64, currency string) error {
			log.Printf("🎉 Payment Captured!")
			log.Printf("   Payment ID: %s", paymentID)
			log.Printf("   Amount: ₹%.2f", float64(amount)/100)
			log.Printf("   Currency: %s", currency)

			// Your business logic here:
			// - Update order status in database
			// - Send confirmation email to customer
			// - Trigger fulfillment process
			// - Update inventory

			return nil
		}),
	)

	// Step 3: Register handler for payment failures
	router.Register("razorpay", "payment.failed",
		webhook.PaymentFailedHandler(func(paymentID string, reason string) error {
			log.Printf("❌ Payment Failed!")
			log.Printf("   Payment ID: %s", paymentID)
			log.Printf("   Reason: %s", reason)

			// Your business logic here:
			// - Notify customer about failure
			// - Log for analytics
			// - Retry payment flow if applicable

			return nil
		}),
	)

	// Step 4: Register handler for refunds
	router.Register("razorpay", "refund.processed",
		webhook.RefundCreatedHandler(func(refundID string, paymentID string, amount int64) error {
			log.Printf("💸 Refund Processed!")
			log.Printf("   Refund ID: %s", refundID)
			log.Printf("   Payment ID: %s", paymentID)
			log.Printf("   Amount: ₹%.2f", float64(amount)/100)

			// Your business logic here:
			// - Update order status
			// - Notify customer about refund
			// - Update accounting records

			return nil
		}),
	)

	// Step 5: Create Fiber web server
	app := fiber.New(fiber.Config{
		AppName: "Razorpay Integration with FinTechKit",
	})

	// Step 6: Webhook endpoint
	app.Post("/webhooks/razorpay", func(c *fiber.Ctx) error {
		// Get webhook payload and signature
		payload := c.Body()
		signature := c.Get("X-Razorpay-Signature")

		// Process the webhook event
		err := receiver.ProcessEvent(context.Background(), "razorpay", payload, signature)
		if err != nil {
			log.Printf("Webhook processing error: %v", err)
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid webhook signature or payload",
			})
		}

		return c.JSON(fiber.Map{
			"status": "received",
		})
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "razorpay-webhook-handler",
		})
	})

	// Start server
	log.Println("🚀 Webhook server starting on :3000")
	log.Println("📡 Razorpay webhook endpoint: http://localhost:3000/webhooks/razorpay")
	log.Fatal(app.Listen(":3000"))
}
