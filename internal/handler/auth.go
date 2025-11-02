package handler

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tsongpon/athena/internal/model"
	"github.com/tsongpon/athena/internal/transport"
)

const (
	// DefaultJWTSecret is used if JWT_SECRET env var is not set
	DefaultJWTSecret = "your-secret-key-change-this-in-production"
	// TokenExpirationHours defines how long the token is valid
	TokenExpirationHours = 24
)

type AuthHandler struct {
	userService UserService
}

func NewAuthHandler(userService UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// JWTClaims represents the claims stored in the JWT token
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c echo.Context) error {
	req := &transport.LoginRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate required fields
	if req.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is required")
	}
	if req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Password is required")
	}

	// Authenticate user via service
	user, err := h.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		// Return generic error to prevent user enumeration
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	// Generate JWT token
	token, expiresAt, err := generateJWT(user.ID, user.Email, user.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Calculate expires_in (seconds until expiration)
	expiresIn := time.Until(expiresAt).Seconds()

	// Build response
	resp := transport.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int64(expiresIn),
		User: transport.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) CreateUser(c echo.Context) error {
	req := &transport.CreateUserRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate required fields
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name is required")
	}
	if req.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is required")
	}
	if req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Password is required")
	}

	// Create user model
	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Will be hashed by service layer
	}

	// Call service to create user
	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		// Check for specific error types
		if err.Error() == "email is required" ||
			err.Error() == "password is required" ||
			err.Error() == "name is required" {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err.Error() == "password length exceeds 72 bytes" {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		// Check for duplicate email error
		if err.Error() != "" && (err.Error() == "user ID must be empty" ||
			containsString(err.Error(), "already exists")) {
			return echo.NewHTTPError(http.StatusConflict, "User with this email already exists")
		}
		// Generic server error
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	// Build response (excluding password)
	resp := transport.UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, resp)
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
