/*
Reads environment variables (DB URL, port, JWT secret, etc.) 
	into a typed Config struct so the rest of the app doesn’t call os.Getenv directly.
*/

package config

import (
	"os"
	"strconv" // for converting string to int
)

type Config struct {
	Port 	  string
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   int // hours
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/taskflow?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "default-secret"),
		JWTExpiry: getEnvInt("JWT_EXPIRY", 24),
	}
}

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is present in the environment, its value (which may be empty) is returned.
// Otherwise, it returns the provided fallback value.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvInt retrieves the value of the environment variable named by the key and converts it to an integer.
func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}