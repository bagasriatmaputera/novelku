package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppPort     string
	CORSOrigins string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret        string
	AccessTTLMinutes int
	RefreshTTLDays   int
}

func Load() *Config {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	cfg := &Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		AppPort:          getEnv("APP_PORT", "8080"),
		CORSOrigins:      getEnv("CORS_ORIGINS", "http://localhost:3000"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "novelku"),
		DBName:           getEnv("DB_NAME", "novelku"),
		DBPassword:       getEnv("DB_PASSWORD", "novelku_dev_password"),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		AccessTTLMinutes: getEnvInt("ACCESS_TOKEN_TTL_MINUTES", 15),
		RefreshTTLDays:   getEnvInt("REFRESH_TOKEN_TTL_DAYS", 30),
	}

	if len(cfg.JWTSecret) < 32 {
		log.Fatalln("JWT SECRETS TO WEAK")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(key); err != nil {
			return n
		}
	}

	return fallback
}
