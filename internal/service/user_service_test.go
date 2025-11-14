package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tsongpon/athena/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	createUserFunc                func(user model.User) (model.User, error)
	getUserByIDFunc               func(id string) (model.User, error)
	getUserByEmailFunc            func(email string) (model.User, error)
	getUserByEmailAndPasswordFunc func(email, hashedPassword string) (model.User, error)
}

func (m *MockUserRepository) CreateUser(user model.User) (model.User, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(user)
	}
	return model.User{}, nil
}

func (m *MockUserRepository) GetUserByID(id string) (model.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(id)
	}
	return model.User{}, nil
}

func (m *MockUserRepository) GetUserByEmail(email string) (model.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(email)
	}
	return model.User{}, nil
}

func (m *MockUserRepository) GetUserByEmailAndPassword(email, hashedPassword string) (model.User, error) {
	if m.getUserByEmailAndPasswordFunc != nil {
		return m.getUserByEmailAndPasswordFunc(email, hashedPassword)
	}
	return model.User{}, nil
}

// TestUserService_CreateUser tests successful user creation with password hashing
func TestUserService_CreateUser(t *testing.T) {
	plainPassword := "myPlainPassword123"
	var capturedHashedPassword string

	mockRepo := &MockUserRepository{
		createUserFunc: func(user model.User) (model.User, error) {
			// Verify password is hashed
			if user.Password == plainPassword {
				t.Error("CreateUser() should hash password before storing")
			}

			// Verify password is a valid bcrypt hash
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))
			if err != nil {
				t.Errorf("CreateUser() stored password is not a valid bcrypt hash: %v", err)
			}

			// Verify tier is set to "free" by default
			if user.Tier != "free" {
				t.Errorf("CreateUser() should set Tier to 'free', got %v", user.Tier)
			}

			capturedHashedPassword = user.Password

			// Return the user with ID
			user.ID = "user-123"
			return user, nil
		},
	}

	service := NewUserService(mockRepo)
	result, err := service.CreateUser(model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: plainPassword,
	})

	if err != nil {
		t.Errorf("CreateUser() unexpected error = %v", err)
		return
	}

	if result.ID != "user-123" {
		t.Errorf("CreateUser() result ID = %v, want user-123", result.ID)
	}

	// Verify tier is set to "free"
	if result.Tier != "free" {
		t.Errorf("CreateUser() result Tier = %v, want free", result.Tier)
	}

	// Verify password was hashed
	if result.Password == plainPassword {
		t.Error("CreateUser() should not return plain text password")
	}

	if capturedHashedPassword == "" {
		t.Error("CreateUser() did not capture hashed password")
	}

	// Verify the hash starts with bcrypt prefix
	if !strings.HasPrefix(result.Password, "$2a$") && !strings.HasPrefix(result.Password, "$2b$") {
		t.Errorf("CreateUser() password doesn't look like bcrypt hash: %s", result.Password)
	}
}

// TestUserService_CreateUser_DefaultTierAlwaysFree tests that tier is always set to "free" for new users
func TestUserService_CreateUser_DefaultTierAlwaysFree(t *testing.T) {
	plainPassword := "password123"

	mockRepo := &MockUserRepository{
		createUserFunc: func(user model.User) (model.User, error) {
			// Verify tier is overridden to "free" even if user tries to set it
			if user.Tier != "free" {
				t.Errorf("CreateUser() should override Tier to 'free', got %v", user.Tier)
			}
			user.ID = "user-123"
			return user, nil
		},
	}

	service := NewUserService(mockRepo)

	// Try to create user with "paid" tier
	result, err := service.CreateUser(model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: plainPassword,
		Tier:     "paid", // This should be ignored and overridden to "free"
	})

	if err != nil {
		t.Errorf("CreateUser() unexpected error = %v", err)
		return
	}

	// Verify tier is "free" regardless of what was provided
	if result.Tier != "free" {
		t.Errorf("CreateUser() result Tier = %v, want free (should ignore user-provided tier)", result.Tier)
	}
}

// TestUserService_CreateUser_EmptyEmail tests validation for empty email
func TestUserService_CreateUser_EmptyEmail(t *testing.T) {
	mockRepo := &MockUserRepository{
		createUserFunc: func(user model.User) (model.User, error) {
			t.Error("CreateUser() should not call repository with empty email")
			return model.User{}, nil
		},
	}

	service := NewUserService(mockRepo)
	_, err := service.CreateUser(model.User{
		Name:     "John Doe",
		Email:    "",
		Password: "password123",
	})

	if err == nil {
		t.Error("CreateUser() should return error when email is empty")
		return
	}

	expectedError := "email is required"
	if err.Error() != expectedError {
		t.Errorf("CreateUser() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestUserService_CreateUser_EmptyPassword tests validation for empty password
func TestUserService_CreateUser_EmptyPassword(t *testing.T) {
	mockRepo := &MockUserRepository{
		createUserFunc: func(user model.User) (model.User, error) {
			t.Error("CreateUser() should not call repository with empty password")
			return model.User{}, nil
		},
	}

	service := NewUserService(mockRepo)
	_, err := service.CreateUser(model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "",
	})

	if err == nil {
		t.Error("CreateUser() should return error when password is empty")
		return
	}

	expectedError := "password is required"
	if err.Error() != expectedError {
		t.Errorf("CreateUser() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestUserService_CreateUser_EmptyName tests validation for empty name
func TestUserService_CreateUser_EmptyName(t *testing.T) {
	mockRepo := &MockUserRepository{
		createUserFunc: func(user model.User) (model.User, error) {
			t.Error("CreateUser() should not call repository with empty name")
			return model.User{}, nil
		},
	}

	service := NewUserService(mockRepo)
	_, err := service.CreateUser(model.User{
		Name:     "",
		Email:    "john@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Error("CreateUser() should return error when name is empty")
		return
	}

	expectedError := "name is required"
	if err.Error() != expectedError {
		t.Errorf("CreateUser() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestUserService_CreateUser_WithID tests that providing an ID returns an error
func TestUserService_CreateUser_WithID(t *testing.T) {
	mockRepo := &MockUserRepository{
		createUserFunc: func(user model.User) (model.User, error) {
			t.Error("CreateUser() should not call repository when ID is provided")
			return model.User{}, nil
		},
	}

	service := NewUserService(mockRepo)
	_, err := service.CreateUser(model.User{
		ID:       "existing-id",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Error("CreateUser() should return error when ID is provided")
		return
	}

	expectedError := "user ID must be empty"
	if err.Error() != expectedError {
		t.Errorf("CreateUser() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestUserService_CreateUser_RepositoryError tests error handling when repository fails
func TestUserService_CreateUser_RepositoryError(t *testing.T) {
	mockRepo := &MockUserRepository{
		createUserFunc: func(user model.User) (model.User, error) {
			return model.User{}, fmt.Errorf("database connection failed")
		},
	}

	service := NewUserService(mockRepo)
	_, err := service.CreateUser(model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Error("CreateUser() should return error when repository fails")
		return
	}

	expectedErrorSubstring := "failed to create user"
	if !strings.Contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("CreateUser() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestUserService_CreateUser_DuplicateEmail tests handling of duplicate email
func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	mockRepo := &MockUserRepository{
		createUserFunc: func(user model.User) (model.User, error) {
			return model.User{}, fmt.Errorf("user with email john@example.com already exists")
		},
	}

	service := NewUserService(mockRepo)
	_, err := service.CreateUser(model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Error("CreateUser() should return error for duplicate email")
		return
	}

	expectedErrorSubstring := "already exists"
	if !strings.Contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("CreateUser() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestUserService_AuthenticateUser tests successful user authentication
func TestUserService_AuthenticateUser(t *testing.T) {
	plainPassword := "correctPassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	expectedUser := model.User{
		ID:       "user-123",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: string(hashedPassword),
	}

	mockRepo := &MockUserRepository{
		getUserByEmailFunc: func(email string) (model.User, error) {
			if email != "john@example.com" {
				t.Errorf("AuthenticateUser() received email = %v, want john@example.com", email)
			}
			return expectedUser, nil
		},
	}

	service := NewUserService(mockRepo)
	result, err := service.AuthenticateUser("john@example.com", plainPassword)

	if err != nil {
		t.Errorf("AuthenticateUser() unexpected error = %v", err)
		return
	}

	if result.ID != expectedUser.ID {
		t.Errorf("AuthenticateUser() result ID = %v, want %v", result.ID, expectedUser.ID)
	}
	if result.Email != expectedUser.Email {
		t.Errorf("AuthenticateUser() result Email = %v, want %v", result.Email, expectedUser.Email)
	}
}

// TestUserService_AuthenticateUser_WrongPassword tests authentication with wrong password
func TestUserService_AuthenticateUser_WrongPassword(t *testing.T) {
	correctPassword := "correctPassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	mockRepo := &MockUserRepository{
		getUserByEmailFunc: func(email string) (model.User, error) {
			return model.User{
				ID:       "user-123",
				Email:    "john@example.com",
				Password: string(hashedPassword),
			}, nil
		},
	}

	service := NewUserService(mockRepo)
	_, err := service.AuthenticateUser("john@example.com", "wrongPassword")

	if err == nil {
		t.Error("AuthenticateUser() should return error for wrong password")
		return
	}

	expectedError := "invalid email or password"
	if err.Error() != expectedError {
		t.Errorf("AuthenticateUser() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestUserService_AuthenticateUser_NonExistentEmail tests authentication with non-existent email
func TestUserService_AuthenticateUser_NonExistentEmail(t *testing.T) {
	mockRepo := &MockUserRepository{
		getUserByEmailFunc: func(email string) (model.User, error) {
			return model.User{}, fmt.Errorf("user with email %s not found", email)
		},
	}

	service := NewUserService(mockRepo)
	_, err := service.AuthenticateUser("nonexistent@example.com", "anyPassword")

	if err == nil {
		t.Error("AuthenticateUser() should return error for non-existent email")
		return
	}

	expectedError := "invalid email or password"
	if err.Error() != expectedError {
		t.Errorf("AuthenticateUser() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestUserService_AuthenticateUser_EmptyEmail tests authentication with empty email
func TestUserService_AuthenticateUser_EmptyEmail(t *testing.T) {
	mockRepo := &MockUserRepository{
		getUserByEmailFunc: func(email string) (model.User, error) {
			t.Error("AuthenticateUser() should not call repository with empty email")
			return model.User{}, nil
		},
	}

	service := NewUserService(mockRepo)
	_, err := service.AuthenticateUser("", "password123")

	if err == nil {
		t.Error("AuthenticateUser() should return error when email is empty")
		return
	}

	expectedError := "email is required"
	if err.Error() != expectedError {
		t.Errorf("AuthenticateUser() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestUserService_AuthenticateUser_EmptyPassword tests authentication with empty password
func TestUserService_AuthenticateUser_EmptyPassword(t *testing.T) {
	mockRepo := &MockUserRepository{
		getUserByEmailFunc: func(email string) (model.User, error) {
			t.Error("AuthenticateUser() should not call repository with empty password")
			return model.User{}, nil
		},
	}

	service := NewUserService(mockRepo)
	_, err := service.AuthenticateUser("john@example.com", "")

	if err == nil {
		t.Error("AuthenticateUser() should return error when password is empty")
		return
	}

	expectedError := "password is required"
	if err.Error() != expectedError {
		t.Errorf("AuthenticateUser() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestUserService_AuthenticateUser_ErrorMessageDoesNotLeakInfo tests that error messages don't leak whether email exists
func TestUserService_AuthenticateUser_ErrorMessageDoesNotLeakInfo(t *testing.T) {
	// Test with non-existent email
	mockRepo1 := &MockUserRepository{
		getUserByEmailFunc: func(email string) (model.User, error) {
			return model.User{}, fmt.Errorf("user not found")
		},
	}

	service1 := NewUserService(mockRepo1)
	_, err1 := service1.AuthenticateUser("nonexistent@example.com", "password")

	// Test with wrong password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	mockRepo2 := &MockUserRepository{
		getUserByEmailFunc: func(email string) (model.User, error) {
			return model.User{Password: string(hashedPassword)}, nil
		},
	}

	service2 := NewUserService(mockRepo2)
	_, err2 := service2.AuthenticateUser("exists@example.com", "wrong")

	// Both errors should be the same to prevent email enumeration
	if err1.Error() != err2.Error() {
		t.Errorf("AuthenticateUser() should return same error message for security. Got:\n  non-existent: %v\n  wrong password: %v", err1, err2)
	}

	expectedError := "invalid email or password"
	if err1.Error() != expectedError {
		t.Errorf("AuthenticateUser() error = %v, want %v", err1.Error(), expectedError)
	}
}

// TestUserService_CreateUser_PasswordComplexity tests that various password complexities are hashed correctly
func TestUserService_CreateUser_PasswordComplexity(t *testing.T) {
	testCases := []struct {
		name        string
		password    string
		shouldError bool
	}{
		{"Simple", "password", false},
		{"With Numbers", "password123", false},
		{"With Symbols", "p@ssw0rd!", false},
		{"Max Length", strings.Repeat("a", 72), false},
		{"Too Long", strings.Repeat("a", 100), true},
		{"Unicode", "пароль密码🔐", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				createUserFunc: func(user model.User) (model.User, error) {
					// Verify password can be verified
					err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tc.password))
					if err != nil {
						t.Errorf("CreateUser() failed to hash password '%s': %v", tc.name, err)
					}
					return user, nil
				},
			}

			service := NewUserService(mockRepo)
			_, err := service.CreateUser(model.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: tc.password,
			})

			if tc.shouldError {
				if err == nil {
					t.Errorf("CreateUser() with %s password should return error", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("CreateUser() with %s password failed: %v", tc.name, err)
				}
			}
		})
	}
}

// TestNewUserService tests service initialization
func TestNewUserService(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	if service == nil {
		t.Error("NewUserService() should return non-nil service")
	}

	if service.repo == nil {
		t.Error("NewUserService() should set repository")
	}
}
