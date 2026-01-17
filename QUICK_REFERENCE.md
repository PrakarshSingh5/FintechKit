# FinTechKit Quick Reference - Razorpay Integration

## 🚀 Quick Start (Copy-Paste Ready)

### 1. Install

```bash
go get github.com/yourusername/fintechkit
```

### 2. Create Client

```go
import "github.com/yourusername/fintechkit/pkg/providers/razorpay"

client, err := razorpay.NewClient(&razorpay.Config{
    KeyID:     "rzp_test_xxxxx",
    KeySecret: "your_secret",
})
```

### 3. Create Payment

```go
payment, err := client.CreatePayment(ctx, &client.PaymentRequest{
    Amount:   50000,  // ₹500 (in paise)
    Currency: "INR",
})
```

### 4. Get Payment Status

```go
payment, err := client.GetPayment(ctx, "order_xxxxx")
```

### 5. Refund Payment

```go
refund, err := client.RefundPayment(ctx, "order_xxxxx", 25000, "Customer request")
```

## 💰 Amount Conversion

| Rupees | Paise (for API) |
| ------ | --------------- |
| ₹1     | 100             |
| ₹10    | 1,000           |
| ₹100   | 10,000          |
| ₹500   | 50,000          |
| ₹1,000 | 1,00,000        |

**Formula**: `Paise = Rupees × 100`

## 🎯 Common Use Cases

### E-commerce Checkout

```go
// Customer buys product for ₹999
payment, _ := client.CreatePayment(ctx, &client.PaymentRequest{
    Amount:      99900,  // ₹999
    Currency:    "INR",
    Description: "Order #12345",
    Metadata: map[string]string{
        "order_id": "12345",
        "user_id":  "user_789",
    },
})
```

### Subscription Payment

```go
// Monthly subscription ₹499
payment, _ := client.CreatePayment(ctx, &client.PaymentRequest{
    Amount:      49900,
    Currency:    "INR",
    Description: "Premium Subscription - Monthly",
})
```

### Partial Refund

```go
// Refund ₹200 from ₹500 order
refund, _ := client.RefundPayment(ctx, paymentID, 20000, "Partial refund")
```

## 🔔 Webhook Events

| Event              | When It Fires      |
| ------------------ | ------------------ |
| `payment.captured` | Payment successful |
| `payment.failed`   | Payment failed     |
| `refund.processed` | Refund completed   |
| `order.paid`       | Order fully paid   |

### Webhook Handler Example

```go
router.Register("razorpay", "payment.captured",
    webhook.PaymentSucceededHandler(func(paymentID string, amount int64, currency string) error {
        log.Printf("Payment %s succeeded: ₹%.2f", paymentID, float64(amount)/100)
        return nil
    }),
)
```

## 🔐 Security Checklist

- [ ] Never commit API keys to Git
- [ ] Use environment variables for credentials
- [ ] Verify webhook signatures
- [ ] Use HTTPS for webhook endpoints
- [ ] Test with test keys before going live
- [ ] Implement idempotency for payments

## 🧪 Testing

### Test Credentials

```go
KeyID:     "rzp_test_xxxxx"  // From Razorpay Dashboard
KeySecret: "test_secret"      // From Razorpay Dashboard
```

### Test Cards

| Card Number         | Result  |
| ------------------- | ------- |
| 4111 1111 1111 1111 | Success |
| 5555 5555 5555 4444 | Success |
| 4000 0000 0000 0002 | Failure |

**CVV**: Any 3 digits  
**Expiry**: Any future date

## 🐛 Common Errors

### "Authentication failed"

```go
// ❌ Wrong
KeyID: "your_key_id"

// ✅ Correct
KeyID: "rzp_test_xxxxx"  // Must start with rzp_test_ or rzp_live_
```

### "Amount mismatch"

```go
// ❌ Wrong
Amount: 500  // This is ₹5, not ₹500!

// ✅ Correct
Amount: 50000  // ₹500 in paise
```

### "Webhook signature verification failed"

```go
// Make sure you're using the webhook secret from Razorpay Dashboard
// Not the API Key Secret!
```

## 📊 Payment Flow Diagram

```
Customer → Your Website → FinTechKit → Razorpay → Bank
                ↓                          ↓
           Order Created              Money Transferred
                ↓                          ↓
           Show Checkout              Webhook Sent
                ↓                          ↓
           Wait for Payment           Verify & Process
                ↓                          ↓
           Fulfill Order              Update Database
```

## 🎓 Learning Path

1. **Day 1**: Read `RAZORPAY_INTEGRATION_GUIDE.md`
2. **Day 2**: Run `examples/razorpay-integration/basic/main.go`
3. **Day 3**: Set up webhook server
4. **Day 4**: Integrate into your website
5. **Day 5**: Test thoroughly
6. **Day 6**: Go live!

## 📞 Support

- **Razorpay Issues**: support@razorpay.com
- **FinTechKit Issues**: GitHub Issues
- **Documentation**: `/examples/razorpay-integration/README.md`

## 🎉 Success Checklist

- [ ] Razorpay account created
- [ ] API keys obtained
- [ ] FinTechKit installed
- [ ] Basic payment working
- [ ] Webhook handler set up
- [ ] Test payments successful
- [ ] Production keys configured
- [ ] Go live!

---

**Pro Tip**: Start with test mode, perfect your integration, then switch to live keys!
