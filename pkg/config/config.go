package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port            int
	Host            string
	GracefulTimeout int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
	MaxConns     int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnvAsInt("SERVER_PORT", 8080),
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			GracefulTimeout: getEnvAsInt("SERVER_GRACEFUL_TIMEOUT", 30),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnvAsInt("DB_PORT", 5432),
			Username:     getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", "postgres"),
			DatabaseName: getEnv("DB_NAME", "scheme_manager"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxConns:     getEnvAsInt("DB_MAX_CONNS", 10),
		},
	}

	return cfg, nil
}

// PostgresConnectionString returns the connection string for PostgreSQL
func (c *DatabaseConfig) PostgresConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.DatabaseName, c.SSLMode,
	)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}