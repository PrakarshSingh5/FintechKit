# Razorpay Integration with FinTechKit

This example demonstrates how to integrate Razorpay payment gateway into your application using FinTechKit.

## 🎯 What This Example Shows

1. **Payment Creation** - Create Razorpay orders programmatically
2. **Payment Retrieval** - Fetch payment details by ID
3. **Refund Processing** - Issue full or partial refunds
4. **Webhook Handling** - Receive real-time payment notifications
5. **Reliability Features** - Automatic retry logic and error handling

## 📋 Prerequisites

Before running this example, you need:

1. **Razorpay Account** - Sign up at [razorpay.com](https://razorpay.com)
2. **API Credentials** - Get your Key ID and Key Secret from Razorpay Dashboard
3. **Go 1.21+** - Install from [golang.org](https://golang.org)

## 🔑 Getting Your Razorpay Credentials

1. Log in to [Razorpay Dashboard](https://dashboard.razorpay.com)
2. Go to **Settings** → **API Keys**
3. Generate Test/Live API keys
4. You'll get:
   - **Key ID**: `rzp_test_xxxxxxxxxx` (for test mode)
   - **Key Secret**: `xxxxxxxxxxxxxxxxxx`

⚠️ **Important**: Never commit your Key Secret to version control!

## 🚀 Running the Example

### 1. Basic Payment Flow

```bash
# Navigate to the example directory
cd examples/razorpay-integration

# Set your Razorpay credentials as environment variables
export RAZORPAY_KEY_ID="rzp_test_your_key_id"
export RAZORPAY_KEY_SECRET="your_key_secret"

# Run the basic payment example
go run main.go
```

**Expected Output:**

```
✅ Payment created successfully!
   Payment ID: order_1234567890
   Amount: ₹500.00
   Currency: INR
   Status: created

📋 Retrieved Payment Details:
   ID: order_1234567890
   Status: paid

💰 Refund processed:
   Refund ID: rfnd_9876543210
   Amount: ₹250.00
   Status: processed
   Reason: Customer requested partial refund
```

### 2. Webhook Server

```bash
# Run the webhook server
go run webhook_server.go
```

**Expected Output:**

```
🚀 Webhook server starting on :3000
📡 Razorpay webhook endpoint: http://localhost:3000/webhooks/razorpay
```

## 🔗 Configuring Razorpay Webhooks

1. Go to [Razorpay Dashboard](https://dashboard.razorpay.com) → **Settings** → **Webhooks**
2. Click **Add New Webhook**
3. Enter your webhook URL: `https://yourdomain.com/webhooks/razorpay`
   - For local testing, use [ngrok](https://ngrok.com): `ngrok http 3000`
4. Select events to listen for:
   - ✅ `payment.captured`
   - ✅ `payment.failed`
   - ✅ `refund.processed`
   - ✅ `order.paid`
5. Copy the **Webhook Secret** for signature verification

## 💡 Understanding Razorpay Amounts

Razorpay uses **paise** (smallest currency unit) instead of rupees:

```go
// ❌ Wrong
Amount: 500  // This is ₹5.00, not ₹500!

// ✅ Correct
Amount: 50000  // This is ₹500.00 (50000 paise)
```

**Conversion Formula:**

```
Amount in Paise = Amount in Rupees × 100
```

## 🏗️ Integration Steps for Your Website

### Step 1: Install FinTechKit

```bash
go get github.com/yourusername/fintechkit
```

### Step 2: Initialize Razorpay Client

```go
import (
    "github.com/yourusername/fintechkit/pkg/providers/razorpay"
)

client, err := razorpay.NewClient(&razorpay.Config{
    KeyID:     "rzp_test_xxxxx",
    KeySecret: "your_secret",
})
```

### Step 3: Create Payment Order

```go
payment, err := client.CreatePayment(ctx, &client.PaymentRequest{
    Amount:      50000,  // ₹500 in paise
    Currency:    "INR",
    Description: "Order #12345",
    Metadata: map[string]string{
        "order_id": "12345",
    },
})
```

### Step 4: Handle Frontend Integration

On your website's frontend, use Razorpay Checkout:

```html
<script src="https://checkout.razorpay.com/v1/checkout.js"></script>
<script>
  var options = {
    key: "rzp_test_xxxxx", // Your Key ID
    amount: "50000", // Amount in paise
    currency: "INR",
    name: "Your Company",
    description: "Order #12345",
    order_id: "order_xxxxx", // Order ID from backend
    handler: function (response) {
      // Payment successful
      console.log(response.razorpay_payment_id);
      console.log(response.razorpay_order_id);
      console.log(response.razorpay_signature);

      // Send to your backend for verification
      verifyPayment(response);
    },
  };

  var rzp = new Razorpay(options);
  rzp.open();
</script>
```

### Step 5: Verify Payment on Backend

```go
// Retrieve payment to verify status
payment, err := client.GetPayment(ctx, paymentID)
if err != nil {
    // Handle error
}

if payment.Status == "paid" {
    // Payment successful - fulfill order
    fulfillOrder(payment.ID)
}
```

## 🔐 Security Best Practices

1. **Never expose Key Secret** - Keep it server-side only
2. **Verify webhook signatures** - Always validate incoming webhooks
3. **Use HTTPS** - Especially for webhook endpoints
4. **Implement idempotency** - Prevent duplicate payment processing
5. **Log everything** - Maintain audit trail for compliance

## 🎨 Complete E-commerce Flow

```
1. Customer adds items to cart
   ↓
2. Backend creates Razorpay order
   ↓
3. Frontend shows Razorpay Checkout
   ↓
4. Customer completes payment
   ↓
5. Razorpay sends webhook to your server
   ↓
6. Backend verifies payment
   ↓
7. Order is fulfilled
```

## 📊 Testing

### Test Card Numbers

Razorpay provides test cards for development:

| Card Number         | Type       | Result  |
| ------------------- | ---------- | ------- |
| 4111 1111 1111 1111 | Visa       | Success |
| 5555 5555 5555 4444 | Mastercard | Success |
| 4000 0000 0000 0002 | Visa       | Failure |

**CVV**: Any 3 digits  
**Expiry**: Any future date

## 🐛 Troubleshooting

### Issue: "Authentication failed"

**Solution**: Check your Key ID and Key Secret are correct

### Issue: "Webhook signature verification failed"

**Solution**: Ensure you're using the correct webhook secret from Razorpay Dashboard

### Issue: "Amount mismatch"

**Solution**: Remember to use paise (multiply rupees by 100)

## 📚 Additional Resources

- [Razorpay API Documentation](https://razorpay.com/docs/api/)
- [Razorpay Webhooks Guide](https://razorpay.com/docs/webhooks/)
- [FinTechKit Documentation](../../README.md)

## 🤝 Support

For issues specific to:

- **Razorpay**: Contact [Razorpay Support](https://razorpay.com/support/)
- **FinTechKit**: Open an issue on GitHub

---

**Made with ❤️ for Indian developers**
