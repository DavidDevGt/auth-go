package main

import (
	"github.com/gofiber/fiber/v2"
	"log"

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

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
