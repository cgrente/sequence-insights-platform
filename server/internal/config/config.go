package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
Package config contains the application configuration loaded from environment variables.

Best practice notes:
- Keep config parsing isolated and unit-testable.
- Use explicit defaults.
- Fail fast for required values in production, but allow local defaults.
*/

// Config holds all runtime configuration for the API service.
type Config struct {
	Env            string
	HTTPHost       string
	HTTPPort       int
	BaseURL        string
	RequestTimeout time.Duration

	APIKey string

	DBUrl             string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxIdleTime time.Duration
	DBConnMaxLifetime time.Duration

	WorkerCount  int
	JobQueueSize int
}

// LoadConfig reads configuration from environment variables and returns a Config struct.
func LoadConfig() (Config, error) {
	cfg := Config{
		Env:            getEnv("APP_ENV", "development"),
		HTTPHost:       getEnv("HTTP_HOST", "localhost"),
		HTTPPort:       getEnvAsInt("HTTP_PORT", 8080),
		RequestTimeout: getDurationEnv("REQUEST_TIMEOUT", 10*time.Second),

		APIKey: getEnv("API_KEY", ""),

		DBUrl:             getEnv("DB_URL", "postgres://user:password@localhost:5432/sequence_insights?sslmode=disable"),
		DBMaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
		DBConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", 15*time.Minute),
		DBConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 1*time.Hour),

		WorkerCount:  getEnvAsInt("WORKER_COUNT", 5),
		JobQueueSize: getEnvAsInt("JOB_QUEUE_SIZE", 100),
	}

	// Validate required configurations
	if strings.EqualFold(cfg.Env, "prod") && strings.TrimSpace(cfg.DBUrl) == "" {
		return Config{}, errors.New("DB_URL is required in production environment")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultVal
	}
	return value
}

func getEnvAsInt(key string, defaultVal int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultVal
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultVal
	}
	return i
}

func getDurationEnv(key string, defaultVal time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultVal
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return defaultVal
	}
	return d
}
