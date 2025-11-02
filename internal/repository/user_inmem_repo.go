package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tsongpon/athena/internal/model"
)

// UserInMemRepository implements UserRepository interface using an in-memory map
type UserInMemRepository struct {
	users map[string]model.User
	mutex sync.RWMutex
}

// NewUserInMemRepository creates a new instance of UserInMemRepository
func NewUserInMemRepository() *UserInMemRepository {
	return &UserInMemRepository{
		users: make(map[string]model.User),
		mutex: sync.RWMutex{},
	}
}

// CreateUser creates a new user in the repository
func (r *UserInMemRepository) CreateUser(user model.User) (model.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Check if user with same email already exists
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return model.User{}, fmt.Errorf("user with email %s already exists", user.Email)
		}
	}

	// Set creation and update times
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	// Store the user
	r.users[user.ID] = user

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (r *UserInMemRepository) GetUserByID(id string) (model.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return model.User{}, fmt.Errorf("user with ID %s not found", id)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (r *UserInMemRepository) GetUserByEmail(email string) (model.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return model.User{}, fmt.Errorf("user with email %s not found", email)
}

// GetUserByEmailAndPassword retrieves a user by email and password
func (r *UserInMemRepository) GetUserByEmailAndPassword(email, hashedPassword string) (model.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email && user.Password == hashedPassword {
			return user, nil
		}
	}

	return model.User{}, fmt.Errorf("user not found with provided credentials")
}
