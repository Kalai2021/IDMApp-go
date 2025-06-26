package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"idmapp-go/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserController_Login(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	
	// Create a mock user service (you would need to create a mock service)
	// For now, we'll just test the request binding
	userController := &UserController{}
	
	// Add the login route
	router.POST("/login", userController.Login)

	// Test valid login request
	t.Run("Valid Login Request", func(t *testing.T) {
		loginReq := dto.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		
		jsonData, _ := json.Marshal(loginReq)
		
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		// Since we don't have a mock service, we expect an error
		// In a real test, you would mock the service and test the actual response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test invalid login request (missing email)
	t.Run("Invalid Login Request - Missing Email", func(t *testing.T) {
		loginReq := dto.LoginRequest{
			Password: "password123",
		}
		
		jsonData, _ := json.Marshal(loginReq)
		
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test invalid login request (invalid email format)
	t.Run("Invalid Login Request - Invalid Email", func(t *testing.T) {
		loginReq := dto.LoginRequest{
			Email:    "invalid-email",
			Password: "password123",
		}
		
		jsonData, _ := json.Marshal(loginReq)
		
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
} 