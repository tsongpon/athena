package repository

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/tsongpon/athena/internal/model"
)

// Helper function to clean up user test data
func cleanupUserTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		t.Logf("Warning: Failed to clean up user test data: %v", err)
	}
}

func TestUserPostgresRepository_CreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupUserTestDB(t, db)

	repo := NewUserPostgresRepository(db)

	user := model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	result, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser() unexpected error = %v", err)
	}

	// ID should be auto-generated
	if result.ID == "" {
		t.Error("CreateUser() result ID should not be empty")
	}
	if result.Name != user.Name {
		t.Errorf("CreateUser() result Name = %v, want %v", result.Name, user.Name)
	}
	if result.Email != user.Email {
		t.Errorf("CreateUser() result Email = %v, want %v", result.Email, user.Email)
	}
	if result.Password != user.Password {
		t.Errorf("CreateUser() result Password = %v, want %v", result.Password, user.Password)
	}
	if result.CreatedAt.IsZero() {
		t.Error("CreateUser() result CreatedAt should not be zero")
	}
	if result.UpdatedAt.IsZero() {
		t.Error("CreateUser() result UpdatedAt should not be zero")
	}
}

func TestUserPostgresRepository_GetUserByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupUserTestDB(t, db)

	repo := NewUserPostgresRepository(db)

	// Create a test user
	user := model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	created, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get the user
	result, err := repo.GetUserByID(created.ID)
	if err != nil {
		t.Fatalf("GetUserByID() unexpected error = %v", err)
	}

	if result.ID != created.ID {
		t.Errorf("GetUserByID() result ID = %v, want %v", result.ID, created.ID)
	}
	if result.Name != user.Name {
		t.Errorf("GetUserByID() result Name = %v, want %v", result.Name, user.Name)
	}
	if result.Email != user.Email {
		t.Errorf("GetUserByID() result Email = %v, want %v", result.Email, user.Email)
	}
}

func TestUserPostgresRepository_GetUserByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupUserTestDB(t, db)

	repo := NewUserPostgresRepository(db)

	_, err := repo.GetUserByID("nonexistent")
	if err == nil {
		t.Error("GetUserByID() with non-existent ID should return error")
	}
}

func TestUserPostgresRepository_GetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupUserTestDB(t, db)

	repo := NewUserPostgresRepository(db)

	// Create a test user
	user := model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	created, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get the user by email
	result, err := repo.GetUserByEmail(user.Email)
	if err != nil {
		t.Fatalf("GetUserByEmail() unexpected error = %v", err)
	}

	if result.ID != created.ID {
		t.Errorf("GetUserByEmail() result ID = %v, want %v", result.ID, created.ID)
	}
	if result.Email != user.Email {
		t.Errorf("GetUserByEmail() result Email = %v, want %v", result.Email, user.Email)
	}
}

func TestUserPostgresRepository_GetUserByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupUserTestDB(t, db)

	repo := NewUserPostgresRepository(db)

	_, err := repo.GetUserByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("GetUserByEmail() with non-existent email should return error")
	}
}

func TestUserPostgresRepository_GetUserByEmailAndPassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupUserTestDB(t, db)

	repo := NewUserPostgresRepository(db)

	// Create a test user
	user := model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	created, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get the user by email and password
	result, err := repo.GetUserByEmailAndPassword(user.Email, user.Password)
	if err != nil {
		t.Fatalf("GetUserByEmailAndPassword() unexpected error = %v", err)
	}

	if result.ID != created.ID {
		t.Errorf("GetUserByEmailAndPassword() result ID = %v, want %v", result.ID, created.ID)
	}
	if result.Email != user.Email {
		t.Errorf("GetUserByEmailAndPassword() result Email = %v, want %v", result.Email, user.Email)
	}
}

func TestUserPostgresRepository_GetUserByEmailAndPassword_WrongPassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupUserTestDB(t, db)

	repo := NewUserPostgresRepository(db)

	// Create a test user
	user := model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	_, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Try to get user with wrong password
	_, err = repo.GetUserByEmailAndPassword(user.Email, "wrong_password")
	if err == nil {
		t.Error("GetUserByEmailAndPassword() with wrong password should return error")
	}
}

func TestUserPostgresRepository_GetUserByEmailAndPassword_WrongEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupUserTestDB(t, db)

	repo := NewUserPostgresRepository(db)

	// Create a test user
	user := model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	_, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Try to get user with wrong email
	_, err = repo.GetUserByEmailAndPassword("wrong@example.com", user.Password)
	if err == nil {
		t.Error("GetUserByEmailAndPassword() with wrong email should return error")
	}
}

func TestUserPostgresRepository_CreateUser_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupUserTestDB(t, db)

	repo := NewUserPostgresRepository(db)

	// Create a test user
	user := model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	_, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Try to create another user with the same email
	user2 := model.User{
		Name:     "Jane Doe",
		Email:    "john@example.com", // Same email
		Password: "another_password",
	}

	_, err = repo.CreateUser(user2)
	if err == nil {
		t.Error("CreateUser() with duplicate email should return error")
	}
}
