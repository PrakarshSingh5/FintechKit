package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/PrakarshSingh5/fintechkit/pkg/auth"
)

// AuthMiddleware creates Fiber middleware for API key authentication
func AuthMiddleware(manager *auth.Manager, providerID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get credentials
		creds, err := manager.GetCredentials(c.Context(), providerID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Store credentials in context for downstream handlers
		c.Locals("credentials", creds)
		c.Locals("provider_id", providerID)

		return c.Next()
	}
}

// BearerTokenMiddleware extracts and validates bearer tokens
func BearerTokenMiddleware(manager *auth.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract Bearer token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Parse Bearer token
		var token string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		// Store token in context
		c.Locals("access_token", token)

		return c.Next()
	}
}

// APIKeyMiddleware validates API key from header or query param
func APIKeyMiddleware(expectedKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check header first
		apiKey := c.Get("X-API-Key")

		// Fall back to query parameter
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key required",
			})
		}

		if apiKey != expectedKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		return c.Next()
	}
}
