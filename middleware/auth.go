package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Local JWT signing key (same as in PKCE service)
var jwtSigningKey = []byte("your-256-bit-secret") // Replace with a secure key in production

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// For development/testing, accept a simple test token
		if tokenString == "test-token" {
			c.Set("user_id", "test-user-id")
			c.Set("email", "test@example.com")
			c.Next()
			return
		}

		// Parse token without validation first to get the algorithm
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
		if err != nil {
			logrus.Errorf("Failed to parse token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Check the signing algorithm
		alg, ok := token.Header["alg"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No algorithm in token"})
			c.Abort()
			return
		}

		// Only support HS256 for local tokens
		if alg != "HS256" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Unsupported signing algorithm: %s. Only HS256 is supported.", alg)})
			c.Abort()
			return
		}

		// Validate the token with local signing key
		validatedToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSigningKey, nil
		})

		if err != nil {
			logrus.Errorf("Failed to validate token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := validatedToken.Claims.(*Claims); ok && validatedToken.Valid {
			// Add user info to context
			c.Set("user_id", claims.Sub)
			c.Set("email", claims.Email)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(string)
	}
	return ""
}

func GetUserEmail(c *gin.Context) string {
	if email, exists := c.Get("email"); exists {
		return email.(string)
	}
	return ""
}
