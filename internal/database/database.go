package database

import (
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"auth-go/internal/config"
	"auth-go/internal/database/models"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		cfg := config.LoadConfig()
		if cfg.DatabaseURL == "" {
			log.Fatal("DATABASE_URL is not set in configuration")
		}
		db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		// Auto-migrate models
		err = db.AutoMigrate(&models.User{}, &models.Session{})
		if err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
		dbInstance = db
	})
	return dbInstance
}
