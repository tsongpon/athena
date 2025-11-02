package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsongpon/athena/internal/model"
	"github.com/tsongpon/athena/internal/transport"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) AuthenticateUser(email, password string) (model.User, error) {
	args := m.Called(email, password)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserService) CreateUser(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

// Test NewAuthHandler
func TestNewAuthHandler(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.userService)
}

// Test Login - Success
func TestAuthHandler_Login_Success(t *testing.T) {
	e := echo.New()
	loginJSON := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	expectedUser := model.User{
		ID:        "user123",
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("AuthenticateUser", "test@example.com", "password123").Return(expectedUser, nil)

	err := handler.Login(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response transport.LoginResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "Bearer", response.TokenType)
	assert.Greater(t, response.ExpiresIn, int64(0))
	assert.Equal(t, expectedUser.ID, response.User.ID)
	assert.Equal(t, expectedUser.Name, response.User.Name)
	assert.Equal(t, expectedUser.Email, response.User.Email)

	mockService.AssertExpectations(t)
}

// Test Login - Invalid Credentials
func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	e := echo.New()
	loginJSON := `{"email":"test@example.com","password":"wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	mockService.On("AuthenticateUser", "test@example.com", "wrongpassword").Return(model.User{}, errors.New("invalid credentials"))

	err := handler.Login(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Equal(t, "Invalid email or password", httpErr.Message)

	mockService.AssertExpectations(t)
}

// Test Login - Missing Email
func TestAuthHandler_Login_MissingEmail(t *testing.T) {
	e := echo.New()
	loginJSON := `{"password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	err := handler.Login(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "Email is required", httpErr.Message)
}

// Test Login - Missing Password
func TestAuthHandler_Login_MissingPassword(t *testing.T) {
	e := echo.New()
	loginJSON := `{"email":"test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	err := handler.Login(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "Password is required", httpErr.Message)
}

// Test Login - Invalid JSON
func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	e := echo.New()
	loginJSON := `{"email":"test@example.com","password":}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	err := handler.Login(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "Invalid request body", httpErr.Message)
}

// Test Login - Empty Email
func TestAuthHandler_Login_EmptyEmail(t *testing.T) {
	e := echo.New()
	loginJSON := `{"email":"","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	err := handler.Login(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "Email is required", httpErr.Message)
}

// Test CreateUser - Success
func TestAuthHandler_CreateUser_Success(t *testing.T) {
	e := echo.New()
	userJSON := `{"name":"Test User","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	createdUser := model.User{
		ID:        "user123",
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("CreateUser", mock.MatchedBy(func(u model.User) bool {
		return u.Name == "Test User" && u.Email == "test@example.com" && u.Password == "password123"
	})).Return(createdUser, nil)

	err := handler.CreateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response transport.UserResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, response.ID)
	assert.Equal(t, createdUser.Name, response.Name)
	assert.Equal(t, createdUser.Email, response.Email)

	mockService.AssertExpectations(t)
}

// Test CreateUser - Missing Name
func TestAuthHandler_CreateUser_MissingName(t *testing.T) {
	e := echo.New()
	userJSON := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	err := handler.CreateUser(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "Name is required", httpErr.Message)
}

// Test CreateUser - Missing Email
func TestAuthHandler_CreateUser_MissingEmail(t *testing.T) {
	e := echo.New()
	userJSON := `{"name":"Test User","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	err := handler.CreateUser(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "Email is required", httpErr.Message)
}

// Test CreateUser - Missing Password
func TestAuthHandler_CreateUser_MissingPassword(t *testing.T) {
	e := echo.New()
	userJSON := `{"name":"Test User","email":"test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	err := handler.CreateUser(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "Password is required", httpErr.Message)
}

// Test CreateUser - Duplicate Email
func TestAuthHandler_CreateUser_DuplicateEmail(t *testing.T) {
	e := echo.New()
	userJSON := `{"name":"Test User","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	mockService.On("CreateUser", mock.Anything).Return(model.User{}, errors.New("user with email test@example.com already exists"))

	err := handler.CreateUser(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusConflict, httpErr.Code)
	assert.Equal(t, "User with this email already exists", httpErr.Message)

	mockService.AssertExpectations(t)
}

// Test CreateUser - Password Too Long
func TestAuthHandler_CreateUser_PasswordTooLong(t *testing.T) {
	e := echo.New()
	userJSON := `{"name":"Test User","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	mockService.On("CreateUser", mock.Anything).Return(model.User{}, errors.New("password length exceeds 72 bytes"))

	err := handler.CreateUser(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "password length exceeds 72 bytes", httpErr.Message)

	mockService.AssertExpectations(t)
}

// Test CreateUser - Invalid JSON
func TestAuthHandler_CreateUser_InvalidJSON(t *testing.T) {
	e := echo.New()
	userJSON := `{"name":"Test User","email":"test@example.com","password":}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	err := handler.CreateUser(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "Invalid request body", httpErr.Message)
}

// Test CreateUser - Generic Service Error
func TestAuthHandler_CreateUser_GenericError(t *testing.T) {
	e := echo.New()
	userJSON := `{"name":"Test User","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	mockService.On("CreateUser", mock.Anything).Return(model.User{}, errors.New("database connection failed"))

	err := handler.CreateUser(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
	assert.Equal(t, "Failed to create user", httpErr.Message)

	mockService.AssertExpectations(t)
}

// Test generateJWT
func TestGenerateJWT_Success(t *testing.T) {
	token, expiresAt, err := generateJWT("user123", "test@example.com", "Test User")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))
	assert.True(t, expiresAt.Before(time.Now().Add(25*time.Hour)))

	// Verify token can be parsed
	claims, err := validateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
	assert.Equal(t, "athena", claims.Issuer)
}

// Test generateJWT with custom secret
func TestGenerateJWT_WithCustomSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "custom-test-secret")
	defer os.Unsetenv("JWT_SECRET")

	token, expiresAt, err := generateJWT("user123", "test@example.com", "Test User")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))

	// Verify token with custom secret
	claims, err := validateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claims.UserID)
}

// Test ValidateJWT - Success
func TestValidateJWT_Success(t *testing.T) {
	token, _, err := generateJWT("user123", "test@example.com", "Test User")
	assert.NoError(t, err)

	claims, err := validateJWT(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
}

// Test ValidateJWT - Invalid Token
func TestValidateJWT_InvalidToken(t *testing.T) {
	claims, err := validateJWT("invalid.token.here")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

// Test ValidateJWT - Empty Token
func TestValidateJWT_EmptyToken(t *testing.T) {
	claims, err := validateJWT("")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

// Test ValidateJWT - Expired Token
func TestValidateJWT_ExpiredToken(t *testing.T) {
	// Create an expired token
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = DefaultJWTSecret
	}

	expiresAt := time.Now().Add(-1 * time.Hour)
	claims := JWTClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Name:   "Test User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "athena",
			Subject:   "user123",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	validatedClaims, err := validateJWT(tokenString)

	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
}

// Test ValidateJWT - Wrong Signing Method
func TestValidateJWT_WrongSigningMethod(t *testing.T) {
	// Create token with wrong signing method (RS256 instead of HS256)
	expiresAt := time.Now().Add(24 * time.Hour)
	claims := JWTClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Name:   "Test User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "athena",
			Subject:   "user123",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	validatedClaims, err := validateJWT(tokenString)

	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
}

// Test containsString helper function
func TestContainsString(t *testing.T) {
	assert.True(t, containsString("hello world", "world"))
	assert.True(t, containsString("hello world", "hello"))
	assert.True(t, containsString("hello world", "lo wo"))
	assert.True(t, containsString("already exists", "already exists"))
	assert.False(t, containsString("hello", "world"))
	assert.False(t, containsString("short", "this is longer"))
}

// Test findSubstring helper function
func TestFindSubstring(t *testing.T) {
	assert.True(t, findSubstring("hello world", "world"))
	assert.True(t, findSubstring("hello world", "lo"))
	assert.False(t, findSubstring("hello", "world"))
	assert.True(t, findSubstring("already exists", "exists"))
}

// Test JWTClaims structure
func TestJWTClaims_Structure(t *testing.T) {
	claims := JWTClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Name:   "Test User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "athena",
			Subject:   "user123",
		},
	}

	assert.Equal(t, "user123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
	assert.Equal(t, "athena", claims.Issuer)
	assert.Equal(t, "user123", claims.Subject)
}
