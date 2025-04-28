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
	// @Summary Registra un nuevo usuario
	// @Description Crea un usuario con nombre, email y contraseña
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param body body RegisterRequest true "Datos de registro"
	// @Success 201 {object} map[string]string
	// @Failure 400 {object} map[string]string
	// @Router /api/register [post]
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
	// @Summary Login de usuario
	// @Description Obtiene tokens de acceso y refresh al validar credenciales
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param body body LoginRequest true "Credenciales"
	// @Success 200 {object} utils.Pair
	// @Failure 400 {object} map[string]string
	// @Failure 401 {object} map[string]string
	// @Router /api/login [post]
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
	// @Summary Refresca tokens
	// @Description Genera nuevos tokens usando un refresh token válido
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param body body RefreshRequest true "Refresh token y device ID"
	// @Success 200 {object} utils.Pair
	// @Failure 400 {object} map[string]string
	// @Failure 401 {object} map[string]string
	// @Router /api/refresh [post]
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
	// @Summary Logout
	// @Description Revoca un refresh token
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param body body LogoutRequest true "Refresh token"
	// @Success 200 {object} map[string]string
	// @Failure 400 {object} map[string]string
	// @Router /api/logout [post]
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
	// @Summary Revoca sesión específica
	// @Description Revoca un refresh token específico
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param body body RevokeSessionRequest true "Refresh token"
	// @Success 200 {object} map[string]string
	// @Failure 400 {object} map[string]string
	// @Router /api/revoke-session [post]
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
	// @Summary Ruta protegida de ejemplo
	// @Description Endpoint de prueba que requiere autenticación
	// @Tags auth
	// @Produce json
	// @Security ApiKeyAuth
	// @Success 200 {object} map[string]interface{}
	// @Router /api/protected [get]
	api := app.Group("/api", middleware.AuthMiddleware(tokenManager))
	api.Get("/protected", func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		return c.JSON(fiber.Map{"message": "Ruta protegida!", "user_id": userID})
	})
}
