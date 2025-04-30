package main

import (
	"log"

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
	tokenManager, err := utils.NewManager()
	// Api killed if token manager fails to init
	if err != nil {
		log.Fatalf("token manager init: %v", err)
	}
	authService := services.NewAuthService(db, userService, sessionService, tokenManager)

	routes.RegisterRoutes(app, db, tokenManager, authService)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
