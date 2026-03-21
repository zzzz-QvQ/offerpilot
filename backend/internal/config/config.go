package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	MySQLDSN       string
	JWTSecret      string
	JWTExpireHours int
}

func Load() Config {
	return Config{
		Port:           getEnv("APP_PORT", "8080"),
		MySQLDSN:       getEnv("MYSQL_DSN", "root:password@tcp(127.0.0.1:3306)/offerpilot?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret:      getEnv("APP_JWT_SECRET", "change_me"),
		JWTExpireHours: getEnvAsInt("APP_JWT_EXPIRE_HOURS", 72),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
