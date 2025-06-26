package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// JWKS represents the JSON Web Key Set from Auth0
type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

var jwksCache *JWKS
var jwksCacheTime time.Time

// Local JWT signing key (same as in PKCE service)
var jwtSigningKey = []byte("your-256-bit-secret") // Replace with a secure key in production

func getJWKS(auth0Domain string) (*JWKS, error) {
	// Cache JWKS for 1 hour
	if jwksCache != nil && time.Since(jwksCacheTime) < time.Hour {
		return jwksCache, nil
	}

	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", auth0Domain))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	jwksCache = &jwks
	jwksCacheTime = time.Now()
	return &jwks, nil
}

// Convert JWK to RSA public key
func jwkToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decode the modulus (n)
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %v", err)
	}

	// Decode the exponent (e)
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %v", err)
	}

	// Convert to big integers
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	// Create RSA public key
	publicKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return publicKey, nil
}

func AuthMiddleware(auth0Domain, auth0Audience string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("[DEBUG] Token String:", tokenString)
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

		// Parse token without validation first to get the key ID and algorithm
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

		var validatedToken *jwt.Token

		// Handle different signing algorithms
		switch alg {
		case "HS256":
			// Use local signing key for HS256 tokens (our PKCE implementation)
			validatedToken, err = jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSigningKey, nil
			})
		case "RS256":
			// Use Auth0 JWKS for RS256 tokens
			kid, ok := token.Header["kid"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No key ID in token"})
				c.Abort()
				return
			}

			// Get JWKS from Auth0
			jwks, err := getJWKS(auth0Domain)
			if err != nil {
				logrus.Errorf("Failed to get JWKS: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
				c.Abort()
				return
			}

			// Find the matching key and convert to RSA public key
			var publicKey *rsa.PublicKey
			for _, key := range jwks.Keys {
				if key.Kid == kid {
					publicKey, err = jwkToRSAPublicKey(key)
					if err != nil {
						logrus.Errorf("Failed to convert JWK to RSA key: %v", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process token key"})
						c.Abort()
						return
					}
					break
				}
			}

			if publicKey == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No matching key found"})
				c.Abort()
				return
			}

			// Now validate the token with the correct public key
			validatedToken, err = jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return publicKey, nil
			})
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Unsupported signing algorithm: %s", alg)})
			c.Abort()
			return
		}

		if err != nil {
			logrus.Errorf("Failed to validate token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := validatedToken.Claims.(*Claims); ok && validatedToken.Valid {
			// Validate audience (only for Auth0 tokens)
			if alg == "RS256" && auth0Audience != "" && claims.Audience != nil {
				audienceValid := false
				for _, aud := range claims.Audience {
					if aud == auth0Audience {
						audienceValid = true
						break
					}
				}
				if !audienceValid {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid audience"})
					c.Abort()
					return
				}
			}

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
