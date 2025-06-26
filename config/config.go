package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Auth0    Auth0Config
	OpenFGA  OpenFGAConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

type Auth0Config struct {
	Domain       string
	Audience     string
	ClientID     string
	ClientSecret string
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

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{}

	// Database config
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	config.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		Name:     getEnv("DB_NAME", "sairam"),
		Username: getEnv("DB_USERNAME", "sairam"),
		Password: getEnv("DB_PASSWORD", "sairam"),
	}

	// Auth0 config
	config.Auth0 = Auth0Config{
		Domain:       getEnv("AUTH0_DOMAIN", "dev-df4lud4n6zz4i5tg.us.auth0.com"),
		Audience:     getEnv("AUTH0_AUDIENCE", "/api/v1/"),
		ClientID:     getEnv("AUTH0_CLIENT_ID", "nkMspWucfXetzOpiBXi5sunnmRQNP5QZ"),
		ClientSecret: getEnv("AUTH0_CLIENT_SECRET", "CxsP1j4OL_V3FtK0YyJbEa2qYRF2gSV07WPQzUtPVE-VMYEcmGNKQnFTA5eWs3cx"),
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
