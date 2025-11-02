package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// generateJWT creates a new JWT token for the authenticated user
func generateJWT(userID, email, name string) (string, time.Time, error) {
	// Get JWT secret from environment or use default
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = DefaultJWTSecret
	}

	// Set expiration time
	expiresAt := time.Now().Add(TokenExpirationHours * time.Hour)

	// Create claims
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "athena",
			Subject:   userID,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// getAuthenticatedUser extracts JWTClaims from Echo context
// Echo JWT middleware v4 stores *jwt.Token in context, so we need to extract claims
func getAuthenticatedUser(c echo.Context) (*JWTClaims, error) {
	userFromContext := c.Get("user")

	// Try to get claims directly (in case middleware config changes)
	if claims, ok := userFromContext.(*JWTClaims); ok {
		return claims, nil
	}

	// Extract claims from token (Echo JWT v4 behavior)
	if token, ok := userFromContext.(*jwt.Token); ok {
		if claims, ok := token.Claims.(*JWTClaims); ok {
			return claims, nil
		}
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}

	return nil, echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
}

// ValidateJWT validates a JWT token and returns the claims
// This can be used in middleware for protected routes
func validateJWT(tokenString string) (*JWTClaims, error) {
	// Get JWT secret from environment or use default
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = DefaultJWTSecret
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and return claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
}
