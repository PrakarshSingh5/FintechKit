package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler creates a centralized error handling middleware
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			// Log error (in production, use structured logging)
			// logger.Error(err)

			// Determine status code
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Return error response
			return c.Status(code).JSON(fiber.Map{
				"error":     err.Error(),
				"timestamp": time.Now().Unix(),
			})
		}

		return nil
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log panic (in production, use structured logging)
				// logger.Error("panic recovered", "panic", r)

				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error",
				})
			}
		}()

		return c.Next()
	}
}
