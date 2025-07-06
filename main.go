package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"idmapp-go/config"
	"idmapp-go/database"
	"idmapp-go/middleware"
	"idmapp-go/routes"
	"idmapp-go/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Fluentd logger
	services.InitFluentLogger()

	// Set up slog logger
	var level slog.Level
	switch cfg.Server.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create slog logger with JSON handler
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

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
	routes.SetupRoutes(router)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	slog.Info("Starting IDM App server", "port", cfg.Server.Port)

	// Log server startup
	fluentLogger := services.GetFluentLogger()
	fluentLogger.Info("Server starting", map[string]interface{}{
		"port":      cfg.Server.Port,
		"log_level": cfg.Server.LogLevel,
		"mode":      gin.Mode(),
	})

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
