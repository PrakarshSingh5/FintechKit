package middleware

import (
	"golang.org/x/time/rate"

	"github.com/gofiber/fiber/v2"
)

// RateLimitMiddleware creates Fiber middleware for rate limiting
func RateLimitMiddleware(requestsPerSecond float64, burst int) fiber.Handler {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)

	return func(c *fiber.Ctx) error {
		if !limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		}

		return c.Next()
	}
}

// PerUserRateLimitMiddleware creates per-user rate limiting
func PerUserRateLimitMiddleware(requestsPerSecond float64, burst int) fiber.Handler {
	limiters := make(map[string]*rate.Limiter)

	return func(c *fiber.Ctx) error {
		// Get user ID from context (set by auth middleware)
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			userID = c.IP() // Fall back to IP-based limiting
		}

		// Get or create limiter for this user
		limiter, exists := limiters[userID]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
			limiters[userID] = limiter
		}

		if !limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"user_id": userID,
			})
		}

		return c.Next()
	}
}
