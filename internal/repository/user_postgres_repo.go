package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tsongpon/athena/internal/logger"
	"github.com/tsongpon/athena/internal/model"
	"go.uber.org/zap"
)

// UserPostgresRepository implements UserRepository interface using PostgreSQL
type UserPostgresRepository struct {
	db *sql.DB
}

// NewUserPostgresRepository creates a new instance of UserPostgresRepository
func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{
		db: db,
	}
}

// CreateUser creates a new user in PostgreSQL
func (r *UserPostgresRepository) CreateUser(user model.User) (model.User, error) {
	// Generate ID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Set creation and update times
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	query := `
		INSERT INTO users (id, name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		logger.Error("Failed to create user in database",
			zap.String("user_id", user.ID),
			zap.String("email", user.Email),
			zap.Error(err))
		return model.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Debug("Created user in database", zap.String("id", user.ID))
	return user, nil
}

// GetUserByID retrieves a user by their ID from PostgreSQL
func (r *UserPostgresRepository) GetUserByID(id string) (model.User, error) {
	logger.Debug("Getting user from database", zap.String("id", id))

	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return model.User{}, fmt.Errorf("user with ID %s not found", id)
	}

	if err != nil {
		logger.Error("Failed to get user from database",
			zap.String("id", id),
			zap.Error(err))
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email address from PostgreSQL
func (r *UserPostgresRepository) GetUserByEmail(email string) (model.User, error) {
	logger.Debug("Getting user by email from database", zap.String("email", email))

	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return model.User{}, fmt.Errorf("user with email %s not found", email)
	}

	if err != nil {
		logger.Error("Failed to get user by email from database",
			zap.String("email", email),
			zap.Error(err))
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmailAndPassword retrieves a user by email and password from PostgreSQL
func (r *UserPostgresRepository) GetUserByEmailAndPassword(email, hashedPassword string) (model.User, error) {
	logger.Debug("Getting user by email and password from database", zap.String("email", email))

	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1 AND password = $2
	`

	var user model.User
	err := r.db.QueryRow(query, email, hashedPassword).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return model.User{}, fmt.Errorf("user not found with provided credentials")
	}

	if err != nil {
		logger.Error("Failed to get user by credentials from database",
			zap.String("email", email),
			zap.Error(err))
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
