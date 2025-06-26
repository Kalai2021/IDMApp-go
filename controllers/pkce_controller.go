package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"idmapp-go/dto"
	"idmapp-go/internal/user"
	"idmapp-go/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PKCEController struct {
	pkceService *services.PKCEService
	userService *user.UserService
	logger      *logrus.Logger
}

func NewPKCEController(pkceService *services.PKCEService, userService *user.UserService) *PKCEController {
	return &PKCEController{
		pkceService: pkceService,
		userService: userService,
		logger:      logrus.New(),
	}
}

// Shared handler for PKCE authorization logic
func (c *PKCEController) handlePKCEAuth(ctx *gin.Context, req dto.PKCEAuthRequest) {
	if err := c.pkceService.ValidatePKCEFlow(req); err != nil {
		c.logger.Errorf("PKCE validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *user.User
	userEmail, err := ctx.Cookie("session_user")
	if err != nil || userEmail == "" {
		c.logger.Errorf("No authenticated user in session for PKCE authorize")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, err = c.userService.GetUserByEmail(userEmail)
	if err != nil || user == nil {
		c.logger.Errorf("Failed to get user for PKCE authorize: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	// Use the real user ID for PKCE code generation
	code, state, codeVerifier, err := c.pkceService.CreateAuthorizationCode(req, &user.ID)
	if err != nil {
		c.logger.Errorf("Failed to create authorization code: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create authorization code"})
		return
	}

	response := dto.PKCEAuthResponse{
		AuthorizationURL: fmt.Sprintf("http://localhost:3000/callback?code=%s&state=%s", code, state),
		State:            state,
		CodeVerifier:     codeVerifier,
	}

	ctx.JSON(http.StatusOK, response)
}

// POST handler (existing)
func (c *PKCEController) InitiatePKCEAuth(ctx *gin.Context) {
	var req dto.PKCEAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Errorf("Invalid PKCE auth request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.handlePKCEAuth(ctx, req)
}

// GET handler for OIDC compliance
func (c *PKCEController) InitiatePKCEAuthGET(ctx *gin.Context) {
	// Check for session cookie
	userEmail, err := ctx.Cookie("session_user")
	c.logger.Debug("Session user in Initiate PKCE AuthGET " + userEmail)
	c.logger.Debug("All cookies: " + fmt.Sprintf("%v", ctx.Request.Cookies()))

	if err != nil || userEmail == "" {
		// Check if this is a browser request or API request
		acceptHeader := ctx.GetHeader("Accept")
		userAgent := ctx.GetHeader("User-Agent")

		c.logger.Debugf("Accept header: %s", acceptHeader)
		c.logger.Debugf("User-Agent: %s", userAgent)

		// Consider it a browser request if:
		// 1. Accept header contains text/html, OR
		// 2. User-Agent contains browser indicators, OR
		// 3. No Accept header (direct browser navigation)
		isBrowserRequest := (acceptHeader != "" && strings.Contains(acceptHeader, "text/html")) ||
			(userAgent != "" && (strings.Contains(userAgent, "Mozilla") ||
				strings.Contains(userAgent, "Chrome") ||
				strings.Contains(userAgent, "Safari") ||
				strings.Contains(userAgent, "Firefox") ||
				strings.Contains(userAgent, "Edge"))) ||
			acceptHeader == ""

		c.logger.Debugf("Is browser request: %v", isBrowserRequest)

		if isBrowserRequest {
			// Browser request - redirect to login
			redirectURL := "/login?redirect=" + url.QueryEscape(ctx.Request.RequestURI)
			c.logger.Debugf("Redirecting to login: %s", redirectURL)
			ctx.Redirect(http.StatusFound, redirectURL)
			return
		} else {
			// API request - return JSON error
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":             "authentication_required",
				"error_description": "User authentication required",
				"login_url":         "/login?redirect=" + url.QueryEscape(ctx.Request.RequestURI),
			})
			return
		}
	}

	req := dto.PKCEAuthRequest{
		ClientID:            ctx.Query("client_id"),
		RedirectURI:         ctx.Query("redirect_uri"),
		Scope:               ctx.Query("scope"),
		State:               ctx.Query("state"),
		CodeChallenge:       ctx.Query("code_challenge"),
		CodeChallengeMethod: ctx.Query("code_challenge_method"),
	}

	// Validate PKCE flow
	if err := c.pkceService.ValidatePKCEFlow(req); err != nil {
		c.logger.Errorf("PKCE validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *user.User
	userEmail, err = ctx.Cookie("session_user")
	if err != nil || userEmail == "" {
		c.logger.Errorf("No authenticated user in session for PKCE authorize")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, err = c.userService.GetUserByEmail(userEmail)
	if err != nil || user == nil {
		c.logger.Errorf("Failed to get user for PKCE authorize: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	// Use the real user ID for PKCE code generation
	code, state, codeVerifier, err := c.pkceService.CreateAuthorizationCode(req, &user.ID)
	if err != nil {
		c.logger.Errorf("Failed to create authorization code: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create authorization code"})
		return
	}

	// Check if this is a browser request again
	acceptHeader := ctx.GetHeader("Accept")
	userAgent := ctx.GetHeader("User-Agent")
	isBrowserRequest := (acceptHeader != "" && strings.Contains(acceptHeader, "text/html")) ||
		(userAgent != "" && (strings.Contains(userAgent, "Mozilla") ||
			strings.Contains(userAgent, "Chrome") ||
			strings.Contains(userAgent, "Safari") ||
			strings.Contains(userAgent, "Firefox") ||
			strings.Contains(userAgent, "Edge"))) ||
		acceptHeader == ""

	if isBrowserRequest {
		// Redirect browser to the callback URL with code and state
		callbackURL := fmt.Sprintf("%s?code=%s&state=%s", req.RedirectURI, code, state)
		ctx.Redirect(http.StatusFound, callbackURL)
		return
	}

	// API clients get JSON
	response := dto.PKCEAuthResponse{
		AuthorizationURL: fmt.Sprintf("%s?code=%s&state=%s", req.RedirectURI, code, state),
		State:            state,
		CodeVerifier:     codeVerifier,
	}
	ctx.JSON(http.StatusOK, response)
}

// ExchangeCodeForToken exchanges authorization code for tokens
func (c *PKCEController) ExchangeCodeForToken(ctx *gin.Context) {
	var req dto.PKCETokenRequest

	// Log the raw request for debugging
	c.logger.Debugf("Token exchange request - Content-Type: %s", ctx.GetHeader("Content-Type"))
	c.logger.Debugf("Token exchange request - Raw body: %s", ctx.Request.Body)

	// Accept both form-encoded and JSON payloads
	if err := ctx.ShouldBind(&req); err != nil {
		c.logger.Errorf("Form binding failed: %v", err)
		// Fallback to JSON
		if err := ctx.ShouldBindJSON(&req); err != nil {
			c.logger.Errorf("JSON binding also failed: %v", err)
			c.logger.Errorf("Invalid token exchange request: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Log the parsed request
	c.logger.Debugf("Parsed token request: %+v", req)

	// Validate required fields
	if req.GrantType != "authorization_code" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "grant_type must be 'authorization_code'"})
		return
	}

	if req.RedirectURI == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "redirect_uri is required"})
		return
	}

	// Exchange code for token
	tokenResponse, err := c.pkceService.ExchangeCodeForToken(req)
	if err != nil {
		c.logger.Errorf("Token exchange failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tokenResponse)
}

// RefreshToken refreshes an access token
func (c *PKCEController) RefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
		ClientID     string `json:"client_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Errorf("Invalid refresh token request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For now, return an error since we haven't implemented refresh tokens yet
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "Refresh tokens not implemented yet"})
}

// GetPKCEConfig returns PKCE configuration for clients
func (c *PKCEController) GetPKCEConfig(ctx *gin.Context) {
	config := gin.H{
		"issuer": "http://localhost:8090",
		"pkce": gin.H{
			"code_challenge_method": "S256",
			"supported_scopes": []string{
				"openid",
				"profile",
				"email",
				"offline_access",
			},
		},
	}

	ctx.JSON(http.StatusOK, config)
}

// GetJWKS returns the JSON Web Key Set for JWT validation
func (c *PKCEController) GetJWKS(ctx *gin.Context) {
	// For HS256, we don't need to expose the key in JWKS
	// This is just for demonstration - in production, use RS256
	jwks := gin.H{
		"keys": []gin.H{
			{
				"kty": "oct",
				"use": "sig",
				"alg": "HS256",
				"kid": "default",
				// Note: In production with HS256, you wouldn't expose the key
				// This is just for testing - use RS256 for real applications
			},
		},
	}

	ctx.JSON(http.StatusOK, jwks)
}

// GetOIDCConfig returns the OIDC discovery document
func (c *PKCEController) GetOIDCConfig(ctx *gin.Context) {
	issuer := "http://localhost:8090"
	ctx.JSON(200, gin.H{
		"issuer":                                issuer,
		"authorization_endpoint":                issuer + "/api/v1/auth/pkce/authorize",
		"token_endpoint":                        issuer + "/api/v1/auth/pkce/token",
		"jwks_uri":                              issuer + "/api/v1/auth/pkce/jwks",
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"HS256"},
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"token_endpoint_auth_methods_supported": []string{"none"},
		"claims_supported":                      []string{"sub", "email"},
	})
}
