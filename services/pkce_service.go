package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"idmapp-go/dto"
	"idmapp-go/internal/client"
	"idmapp-go/internal/pkce"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var jwtSigningKey = []byte("your-256-bit-secret") // Replace with a secure key in production

type PKCEService struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func NewPKCEService(db *gorm.DB) *PKCEService {
	return &PKCEService{
		logger: logrus.New(),
		db:     db,
	}
}

// GenerateCodeVerifier generates a random code verifier for PKCE
func (s *PKCEService) GenerateCodeVerifier() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// GenerateCodeChallenge generates a code challenge from the code verifier
func (s *PKCEService) GenerateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// GenerateState generates a random state parameter for CSRF protection
func (s *PKCEService) GenerateState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// CreateAuthorizationCode stores a PKCE code in the DB and returns the code and state
func (s *PKCEService) CreateAuthorizationCode(req dto.PKCEAuthRequest, userID *uuid.UUID) (string, string, string, error) {
	// For PKCE, the client generates the code_challenge
	// We store the challenge and will validate it later when the client sends the code_verifier

	// Use the client's state parameter if provided, otherwise generate one
	state := req.State
	if state == "" {
		var err error
		state, err = s.GenerateState()
		if err != nil {
			return "", "", "", fmt.Errorf("failed to generate state: %w", err)
		}
	}

	// Generate a unique authorization code
	code := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%d-%s", time.Now().UnixNano(), state)))

	pkceCode := pkce.PKCECode{
		Code:                code,
		CodeChallenge:       req.CodeChallenge, // Store the client's code challenge
		CodeChallengeMethod: req.CodeChallengeMethod,
		CodeVerifier:        "", // Will be empty until token exchange
		ClientID:            req.ClientID,
		RedirectURI:         req.RedirectURI,
		State:               &state, // Store the state parameter as pointer
		UserID:              userID,
		ExpiresAt:           time.Now().Add(10 * time.Minute),
		Used:                false,
	}
	if err := s.db.Create(&pkceCode).Error; err != nil {
		return "", "", "", fmt.Errorf("failed to store PKCE code: %w", err)
	}
	return code, state, "", nil
}

// ExchangeCodeForToken validates the code and code_verifier, then issues a JWT
func (s *PKCEService) ExchangeCodeForToken(req dto.PKCETokenRequest) (*dto.PKCETokenResponse, error) {
	s.logger.Info("=== EXCHANGE CODE FOR TOKEN CALLED ===")
	s.logger.Infof("Received request: %+v", req)

	var pkceCode pkce.PKCECode
	if err := s.db.Where("code = ? AND client_id = ? AND redirect_uri = ? AND used = false", req.Code, req.ClientID, req.RedirectURI).First(&pkceCode).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired authorization code")
	}
	if time.Now().After(pkceCode.ExpiresAt) {
		return nil, fmt.Errorf("authorization code expired")
	}

	// Validate state parameter if provided
	if req.State != "" && (pkceCode.State == nil || *pkceCode.State != req.State) {
		return nil, fmt.Errorf("invalid state parameter")
	}

	if pkceCode.CodeChallengeMethod == "S256" {
		challenge := s.GenerateCodeChallenge(req.CodeVerifier)
		s.logger.Debugf("Code verifier received: %s", req.CodeVerifier)
		s.logger.Debugf("Generated challenge: %s", challenge)
		s.logger.Debugf("Stored challenge: %s", pkceCode.CodeChallenge)
		if challenge != pkceCode.CodeChallenge {
			return nil, fmt.Errorf("invalid code_verifier for S256")
		}
	} else if pkceCode.CodeChallengeMethod == "plain" {
		if req.CodeVerifier != pkceCode.CodeChallenge {
			return nil, fmt.Errorf("invalid code_verifier for plain method")
		}
	} else {
		return nil, fmt.Errorf("unsupported code_challenge_method")
	}
	// Mark code as used
	pkceCode.Used = true
	s.db.Save(&pkceCode)
	// Issue JWT
	claims := jwt.MapClaims{
		"sub": pkceCode.UserID,
		"aud": req.ClientID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSigningKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}
	return &dto.PKCETokenResponse{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:       "openid profile email",
	}, nil
}

// ValidatePKCEFlow validates the PKCE flow parameters
func (s *PKCEService) ValidatePKCEFlow(req dto.PKCEAuthRequest) error {
	if req.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if req.RedirectURI == "" {
		return fmt.Errorf("redirect_uri is required")
	}
	// Load the client and check redirect_uri
	var client client.Client
	if err := s.db.Where("client_id = ?", req.ClientID).First(&client).Error; err != nil {
		return fmt.Errorf("invalid client_id")
	}
	found := false
	for _, uri := range client.RedirectURIs {
		if uri == req.RedirectURI {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("redirect_uri is not registered for this client")
	}
	if req.CodeChallenge == "" {
		return fmt.Errorf("code_challenge is required")
	}
	if req.CodeChallengeMethod != "S256" && req.CodeChallengeMethod != "plain" {
		return fmt.Errorf("code_challenge_method must be 'S256' or 'plain'")
	}
	return nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *PKCEService) RefreshToken(refreshToken string, clientID string) (*dto.PKCETokenResponse, error) {
	// Parse and validate the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSigningKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Extract claims from refresh token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Get user ID from claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	// Generate new access token
	accessToken, err := s.GenerateAccessToken(userID, claims["email"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &dto.PKCETokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:       "openid profile email",
	}, nil
}

// GenerateAccessToken generates a JWT access token for a user
func (s *PKCEService) GenerateAccessToken(userID string, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSigningKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}
