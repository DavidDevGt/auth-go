package routes

import (
	"auth-go/internal/database/models"
	"auth-go/internal/middleware"
	"auth-go/internal/services"
	"auth-go/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RevokeSessionRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RegisterRoutes(app *fiber.App, db *gorm.DB, tokenManager utils.TokenActions, authService services.AuthService) {
	// --- Registro de usuario ---
	app.Post("/api/register", func(c *fiber.Ctx) error {
		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		user := models.User{
			ID:           uuid.New().String(),
			Name:         req.Name,
			Email:        req.Email,
			PasswordHash: req.Password,
		}
		if err := authService.Register(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
	})

	// --- Login ---
	app.Post("/api/login", func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		userAgent := c.Get("User-Agent")
		ip := c.IP()
		pair, err := authService.Login(req.Email, req.Password, userAgent, ip, req.DeviceID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(pair)
	})

	// --- Refresh Token ---
	app.Post("/api/refresh", func(c *fiber.Ctx) error {
		var req RefreshRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		pair, err := authService.Refresh(req.RefreshToken, req.DeviceID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(pair)
	})

	// --- Logout	 ---
	app.Post("/api/logout", func(c *fiber.Ctx) error {
		var req LogoutRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := authService.Logout(req.RefreshToken); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Session revoked"})
	})

	// --- Revocar una sesión específica ---
	app.Post("/api/revoke-session", func(c *fiber.Ctx) error {
		var req RevokeSessionRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := authService.RevokeSession(req.RefreshToken); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Specific session revoked"})
	})

	// --- Ruta protegida de ejemplo ---
	api := app.Group("/api", middleware.AuthMiddleware(tokenManager))
	api.Get("/protected", func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		return c.JSON(fiber.Map{"message": "Ruta protegida!", "user_id": userID})
	})
}
