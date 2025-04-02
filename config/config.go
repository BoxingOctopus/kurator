package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config represents the application configuration
type Config struct {
	// Database configuration
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DSN        string // Complete database connection string

	// API configuration
	APIPort int

	// JWT configuration
	JWTSecret      string
	JWTExpiryHours int
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		// Database configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "media_collection"),
		DSN:        getEnv("DATABASE_URL", ""), // Allow for a complete connection string

		// API configuration
		APIPort: getEnvAsInt("API_PORT", 8080),

		// JWT configuration
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
	}

	// Validate required fields
	if cfg.DSN == "" && cfg.DBPassword == "" {
		return nil, fmt.Errorf("either DATABASE_URL or DB_PASSWORD is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	// If a complete connection string is provided, use it
	if c.DSN != "" {
		return c.DSN
	}

	// Otherwise, build the connection string from individual parameters
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

// GetAPIAddr returns the API server address
func (c *Config) GetAPIAddr() string {
	return fmt.Sprintf(":%d", c.APIPort)
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
