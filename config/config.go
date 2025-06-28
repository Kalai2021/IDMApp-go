package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	OpenFGA  OpenFGAConfig
	Server   ServerConfig
	Logging  LoggingConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

type OpenFGAConfig struct {
	APIURL   string
	StoreID  string
	APIToken string
}

type ServerConfig struct {
	Port     int
	LogLevel string
}

type LoggingConfig struct {
	FluentEnabled  bool
	FluentEndpoint string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{}

	// Database config
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	config.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		Name:     getEnv("DB_NAME", "iamdb"),
		Username: getEnv("DB_USERNAME", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
	}

	// OpenFGA config
	config.OpenFGA = OpenFGAConfig{
		APIURL:   getEnv("OPENFGA_API_URL", "http://localhost:8080"),
		StoreID:  getEnv("OPENFGA_STORE_ID", ""),
		APIToken: getEnv("OPENFGA_API_TOKEN", ""),
	}

	// Server config
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	config.Server = ServerConfig{
		Port:     serverPort,
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	// Logging config
	config.Logging = LoggingConfig{
		FluentEnabled:  getEnv("FLUENT_ENABLED", "false") == "true",
		FluentEndpoint: getEnv("FLUENT_ENDPOINT", "http://localhost:24224"),
	}

	return config, nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.Username, c.Database.Password, c.Database.Name)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
