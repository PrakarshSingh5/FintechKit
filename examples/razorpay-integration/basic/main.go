package main

import (
	"context"
	"log"

	"github.com/PrakarshSingh5/fintechkit/pkg/auth"
	"github.com/PrakarshSingh5/fintechkit/pkg/client"
	"github.com/PrakarshSingh5/fintechkit/pkg/providers/razorpay"
)

func main() {
	// Step 1: Initialize authentication manager
	authManager := auth.NewManager(auth.NewInMemoryStore())

	ctx := context.Background()

	// Step 2: Set your Razorpay credentials
	authManager.SetCredentials(ctx, "razorpay", &auth.Credentials{
		Type:   auth.CredentialTypeAPIKey,
		APIKey: "rzp_test_your_key_id:your_key_secret", // Format: keyID:keySecret
	})

	// Step 3: Create Razorpay client
	razorpayClient, err := razorpay.NewClient(&razorpay.Config{
		KeyID:     "rzp_test_your_key_id",     // Get from Razorpay Dashboard
		KeySecret: "your_razorpay_key_secret", // Get from Razorpay Dashboard
	})
	if err != nil {
		log.Fatal("Failed to create Razorpay client:", err)
	}

	// Step 4: Create a payment
	// (FinTechKit automatically includes retry logic via the reliability package)

	// Create payment request
	paymentReq := &client.PaymentRequest{
		Amount:      50000, // ₹500.00 in paise (Razorpay uses paise, not rupees)
		Currency:    "INR",
		Description: "Order #12345",
		Metadata: map[string]string{
			"order_id":     "12345",
			"customer_id":  "cust_001",
			"product_name": "Premium Subscription",
		},
	}

	// Execute payment creation with retry
	payment, err := razorpayClient.CreatePayment(ctx, paymentReq)
	if err != nil {
		log.Fatal("Payment creation failed:", err)
	}

	log.Printf("✅ Payment created successfully!")
	log.Printf("   Payment ID: %s", payment.ID)
	log.Printf("   Amount: ₹%.2f", float64(payment.Amount)/100)
	log.Printf("   Currency: %s", payment.Currency)
	log.Printf("   Status: %s", payment.Status)

	// Step 5: Retrieve payment details
	retrievedPayment, err := razorpayClient.GetPayment(ctx, payment.ID)
	if err != nil {
		log.Fatal("Failed to retrieve payment:", err)
	}

	log.Printf("\n📋 Retrieved Payment Details:")
	log.Printf("   ID: %s", retrievedPayment.ID)
	log.Printf("   Status: %s", retrievedPayment.Status)

	// Step 6: Process refund (if needed)
	refund, err := razorpayClient.RefundPayment(
		ctx,
		payment.ID,
		25000, // Partial refund: ₹250.00 in paise
		"Customer requested partial refund",
	)
	if err != nil {
		log.Fatal("Refund failed:", err)
	}

	log.Printf("\n💰 Refund processed:")
	log.Printf("   Refund ID: %s", refund.ID)
	log.Printf("   Amount: ₹%.2f", float64(refund.Amount)/100)
	log.Printf("   Status: %s", refund.Status)
	log.Printf("   Reason: %s", refund.Reason)
}
