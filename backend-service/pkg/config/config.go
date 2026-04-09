package config

import (
	"log"
	"os"
	"strconv"

	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/pkg/envload"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	JWTExpirationHours int
	// FrontendURL is a comma-separated list of allowed browser origins for CORS.
	FrontendURL string
	Environment string
}

func Load() *Config {
	envload.Load()
	if os.Getenv("DATABASE_URL") == "" {
		log.Println("Warning: DATABASE_URL is not set — use backend-service/.env or environment variables")
	}

	expirationHours := 24
	if exp := os.Getenv("JWT_EXPIRATION_HOURS"); exp != "" {
		if hours, err := strconv.Atoi(exp); err == nil {
			expirationHours = hours
		}
	}

	return &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpirationHours: expirationHours,
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		Environment: getEnv("ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
