package middleware

import "github.com/gofiber/fiber/v2"

// CORSConfig returns a configured CORS middleware
func CORSConfig() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		c.Set("Access-Control-Max-Age", "86400")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent clickjacking
		c.Set("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Set("X-XSS-Protection", "1; mode=block")

		// Strict transport security
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		return c.Next()
	}
}
