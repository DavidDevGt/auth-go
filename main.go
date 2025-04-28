package main

import (
	"github.com/gofiber/fiber/v2"

	"auth-go/internal/database"
	"auth-go/internal/routes"
	"auth-go/internal/services"
	"auth-go/internal/utils"
)

func main() {
	app := fiber.New()
	db := database.GetDB()
	userService := services.NewUserService(db)
	sessionService := services.NewSessionService(db)
	tokenManager, _ := utils.NewManager()
	authService := services.NewAuthService(db, userService, sessionService, tokenManager)

	routes.RegisterRoutes(app, db, tokenManager, authService)

	app.Listen(":8080")
}
