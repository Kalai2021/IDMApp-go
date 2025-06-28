package main

import (
	"fmt"
	"log"
	"net/http"

	"idmapp-go/config"
	"idmapp-go/database"
	"idmapp-go/middleware"
	"idmapp-go/routes"
	"idmapp-go/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Fluentd logger
	services.InitFluentLogger()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Initialize database
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set Gin mode
	if cfg.Server.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Add logging middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorLoggingMiddleware())

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(router, cfg.Auth0.Domain, cfg.Auth0.Audience)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("Starting IDM App server on port %d", cfg.Server.Port)

	// Log server startup
	logger := services.GetFluentLogger()
	logger.Info("Server starting", map[string]interface{}{
		"port":      cfg.Server.Port,
		"log_level": cfg.Server.LogLevel,
		"mode":      gin.Mode(),
	})

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
