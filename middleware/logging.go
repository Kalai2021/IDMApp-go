package middleware

import (
	"time"

	"idmapp-go/services"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware creates a middleware that logs HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get user ID from context if available
		userID := ""
		if user, exists := c.Get("user_id"); exists {
			if userStr, ok := user.(string); ok {
				userID = userStr
			}
		}

		// Log the request
		logger := services.GetFluentLogger()
		logger.LogRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
			c.Writer.Status(),
			duration,
			userID,
		)
	}
}

// ErrorLoggingMiddleware creates a middleware that logs errors
func ErrorLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			logger := services.GetFluentLogger()
			for _, err := range c.Errors {
				logger.Error("HTTP Error", map[string]interface{}{
					"error":      err.Error(),
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"status":     c.Writer.Status(),
					"client_ip":  c.ClientIP(),
					"user_agent": c.Request.UserAgent(),
				})
			}
		}
	}
}
