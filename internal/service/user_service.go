package service

import (
	"fmt"

	"github.com/tsongpon/athena/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(user model.User) (model.User, error) {
	// Validate required fields
	if user.Email == "" {
		return model.User{}, fmt.Errorf("email is required")
	}
	if user.Password == "" {
		return model.User{}, fmt.Errorf("password is required")
	}
	if user.Name == "" {
		return model.User{}, fmt.Errorf("name is required")
	}

	// User should not provide an ID
	if user.ID != "" {
		return model.User{}, fmt.Errorf("user ID must be empty")
	}

	// Validate password length (bcrypt has a 72-byte limit)
	if len(user.Password) > 72 {
		return model.User{}, fmt.Errorf("password length exceeds 72 bytes")
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)
	user.Tier = "free"

	// Create user in repository
	created, err := s.repo.CreateUser(user)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return created, nil
}

// AuthenticateUser authenticates a user by email and plain text password
func (s *UserService) AuthenticateUser(email, plainPassword string) (model.User, error) {
	// Validate input
	if email == "" {
		return model.User{}, fmt.Errorf("email is required")
	}
	if plainPassword == "" {
		return model.User{}, fmt.Errorf("password is required")
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return model.User{}, fmt.Errorf("invalid email or password")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))
	if err != nil {
		return model.User{}, fmt.Errorf("invalid email or password")
	}

	return user, nil
}
