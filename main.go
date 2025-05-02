package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"auth-go/internal/database"
	"auth-go/internal/routes"
	"auth-go/internal/services"
	"auth-go/internal/utils"
)

func main() {
	app := fiber.New()
	
	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",                // Allows any origin
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Type",
		MaxAge:           86400, // 24 hours
	}))
	
	db := database.GetDB()
	userService := services.NewUserService(db)

	// If fiber app is not initialized, kill the program
	if app == nil {
		log.Fatal("Failed to create Fiber app")
	}

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
