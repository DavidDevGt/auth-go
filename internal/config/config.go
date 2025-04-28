package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpire  int
	RefreshTokenExpire int
	MaxSessionsPerUser int
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, relying on environment variables")
	}
}

func LoadConfig() *Config {
	accessTokenExpiry := getEnvDuration("ACCESS_TOKEN_EXPIRY", 15*time.Minute)
	refreshTokenExpiry := getEnvDuration("REFRESH_TOKEN_EXPIRY", 30*24*time.Hour)

	return &Config{
		Port:               getEnv("PORT", ""),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", ""),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", ""),
		AccessTokenExpire:  int(accessTokenExpiry.Seconds()),
		RefreshTokenExpire: int(refreshTokenExpiry.Seconds()),
		MaxSessionsPerUser: getEnvAsInt("MAX_SESSIONS_PER_USER", 5),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
