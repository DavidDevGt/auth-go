package middleware

import (
	"strings"

	"auth-go/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(tokenManager utils.TokenActions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := tokenManager.ValidateAccessToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}
		c.Locals("userID", claims.Subject)
		return c.Next()
	}
}
