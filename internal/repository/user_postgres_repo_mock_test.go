package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tsongpon/athena/internal/model"
)

// TestUserPostgresRepository_CreateUser_Mock tests user creation with mock DB
func TestUserPostgresRepository_CreateUser_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	user := model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	// Expect INSERT query
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, name, email, password, created_at, updated_at)`)).
		WithArgs(sqlmock.AnyArg(), "John Doe", "john@example.com", "hashed_password", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser() unexpected error = %v", err)
		return
	}

	if result.ID == "" {
		t.Error("CreateUser() should generate ID")
	}
	if result.Name != user.Name {
		t.Errorf("CreateUser() Name = %v, want %v", result.Name, user.Name)
	}
	if result.Email != user.Email {
		t.Errorf("CreateUser() Email = %v, want %v", result.Email, user.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_CreateUser_WithID tests creating user with existing ID
func TestUserPostgresRepository_CreateUser_WithID_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	user := model.User{
		ID:       "existing-id",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users`)).
		WithArgs("existing-id", "John Doe", "john@example.com", "hashed_password", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser() unexpected error = %v", err)
		return
	}

	if result.ID != "existing-id" {
		t.Errorf("CreateUser() should preserve existing ID, got %v", result.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_CreateUser_WithTimestamps tests creating user with existing timestamps
func TestUserPostgresRepository_CreateUser_WithTimestamps_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	existingTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	user := model.User{
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		CreatedAt: existingTime,
		UpdatedAt: existingTime,
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users`)).
		WithArgs(sqlmock.AnyArg(), "John Doe", "john@example.com", "hashed_password", existingTime, existingTime).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser() unexpected error = %v", err)
		return
	}

	if !result.CreatedAt.Equal(existingTime) {
		t.Errorf("CreateUser() should preserve CreatedAt, got %v, want %v", result.CreatedAt, existingTime)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_CreateUser_Error tests error handling
func TestUserPostgresRepository_CreateUser_Error_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	user := model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users`)).
		WillReturnError(errors.New("database error"))

	_, err = repo.CreateUser(user)
	if err == nil {
		t.Error("CreateUser() should return error when database fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_CreateUser_DuplicateEmail tests duplicate email error
func TestUserPostgresRepository_CreateUser_DuplicateEmail_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	user := model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	// Simulate unique constraint violation
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users`)).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))

	_, err = repo.CreateUser(user)
	if err == nil {
		t.Error("CreateUser() should return error for duplicate email")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_GetUserByID_Mock tests retrieving a user by ID
func TestUserPostgresRepository_GetUserByID_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	expectedUser := model.User{
		ID:        "user1",
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password,
			expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`)).
		WithArgs("user1").
		WillReturnRows(rows)

	result, err := repo.GetUserByID("user1")
	if err != nil {
		t.Errorf("GetUserByID() unexpected error = %v", err)
		return
	}

	if result.ID != expectedUser.ID {
		t.Errorf("GetUserByID() ID = %v, want %v", result.ID, expectedUser.ID)
	}
	if result.Email != expectedUser.Email {
		t.Errorf("GetUserByID() Email = %v, want %v", result.Email, expectedUser.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_GetUserByID_NotFound tests not found scenario
func TestUserPostgresRepository_GetUserByID_NotFound_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`)).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetUserByID("nonexistent")
	if err == nil {
		t.Error("GetUserByID() should return error for non-existent user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_GetUserByID_QueryError tests query error
func TestUserPostgresRepository_GetUserByID_QueryError_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`)).
		WithArgs("user1").
		WillReturnError(errors.New("database error"))

	_, err = repo.GetUserByID("user1")
	if err == nil {
		t.Error("GetUserByID() should return error when database fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_GetUserByEmail_Mock tests retrieving a user by email
func TestUserPostgresRepository_GetUserByEmail_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	expectedUser := model.User{
		ID:        "user1",
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password,
			expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`)).
		WithArgs("john@example.com").
		WillReturnRows(rows)

	result, err := repo.GetUserByEmail("john@example.com")
	if err != nil {
		t.Errorf("GetUserByEmail() unexpected error = %v", err)
		return
	}

	if result.Email != expectedUser.Email {
		t.Errorf("GetUserByEmail() Email = %v, want %v", result.Email, expectedUser.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_GetUserByEmail_NotFound tests not found scenario
func TestUserPostgresRepository_GetUserByEmail_NotFound_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`)).
		WithArgs("nonexistent@example.com").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetUserByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("GetUserByEmail() should return error for non-existent email")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_GetUserByEmail_QueryError tests query error
func TestUserPostgresRepository_GetUserByEmail_QueryError_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`)).
		WithArgs("john@example.com").
		WillReturnError(errors.New("database error"))

	_, err = repo.GetUserByEmail("john@example.com")
	if err == nil {
		t.Error("GetUserByEmail() should return error when database fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_GetUserByEmailAndPassword_Mock tests authentication
func TestUserPostgresRepository_GetUserByEmailAndPassword_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	expectedUser := model.User{
		ID:        "user1",
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password,
			expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1 AND password = $2`)).
		WithArgs("john@example.com", "hashed_password").
		WillReturnRows(rows)

	result, err := repo.GetUserByEmailAndPassword("john@example.com", "hashed_password")
	if err != nil {
		t.Errorf("GetUserByEmailAndPassword() unexpected error = %v", err)
		return
	}

	if result.Email != expectedUser.Email {
		t.Errorf("GetUserByEmailAndPassword() Email = %v, want %v", result.Email, expectedUser.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_GetUserByEmailAndPassword_WrongPassword tests wrong password
func TestUserPostgresRepository_GetUserByEmailAndPassword_WrongPassword_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1 AND password = $2`)).
		WithArgs("john@example.com", "wrong_password").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetUserByEmailAndPassword("john@example.com", "wrong_password")
	if err == nil {
		t.Error("GetUserByEmailAndPassword() should return error for wrong password")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_GetUserByEmailAndPassword_WrongEmail tests wrong email
func TestUserPostgresRepository_GetUserByEmailAndPassword_WrongEmail_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1 AND password = $2`)).
		WithArgs("wrong@example.com", "hashed_password").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetUserByEmailAndPassword("wrong@example.com", "hashed_password")
	if err == nil {
		t.Error("GetUserByEmailAndPassword() should return error for wrong email")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestUserPostgresRepository_GetUserByEmailAndPassword_QueryError tests query error
func TestUserPostgresRepository_GetUserByEmailAndPassword_QueryError_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1 AND password = $2`)).
		WithArgs("john@example.com", "hashed_password").
		WillReturnError(errors.New("database error"))

	_, err = repo.GetUserByEmailAndPassword("john@example.com", "hashed_password")
	if err == nil {
		t.Error("GetUserByEmailAndPassword() should return error when database fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestNewUserPostgresRepository tests repository initialization
func TestNewUserPostgresRepository_Mock(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db)
	if repo == nil {
		t.Error("NewUserPostgresRepository() should not return nil")
	}

	if repo.db == nil {
		t.Error("NewUserPostgresRepository() should initialize db")
	}
}
