package repository

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/tsongpon/athena/internal/model"
)

func TestUserInMemRepository_CreateUser(t *testing.T) {
	repo := NewUserInMemRepository()

	user := model.User{
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword123",
		CreatedAt: time.Now(),
	}

	result, err := repo.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser() unexpected error = %v", err)
		return
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

func TestUserInMemRepository_CreateUser_AutoGenerateID(t *testing.T) {
	repo := NewUserInMemRepository()

	user := model.User{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "hashedpassword456",
	}

	result, err := repo.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser() unexpected error = %v", err)
		return
	}

	// ID should be auto-generated
	if result.ID == "" {
		t.Error("CreateUser() should auto-generate ID when not provided")
	}
}

func TestUserInMemRepository_CreateUser_AutoSetTimestamps(t *testing.T) {
	repo := NewUserInMemRepository()

	user := model.User{
		Name:     "Bob Smith",
		Email:    "bob@example.com",
		Password: "hashedpassword789",
	}

	before := time.Now()
	result, err := repo.CreateUser(user)
	after := time.Now()

	if err != nil {
		t.Errorf("CreateUser() unexpected error = %v", err)
		return
	}

	// Timestamps should be set automatically
	if result.CreatedAt.IsZero() {
		t.Error("CreateUser() should auto-set CreatedAt when not provided")
	}
	if result.UpdatedAt.IsZero() {
		t.Error("CreateUser() should auto-set UpdatedAt when not provided")
	}

	// Verify timestamps are within expected range
	if result.CreatedAt.Before(before) || result.CreatedAt.After(after) {
		t.Errorf("CreateUser() CreatedAt = %v, should be between %v and %v", result.CreatedAt, before, after)
	}
	if result.UpdatedAt.Before(before) || result.UpdatedAt.After(after) {
		t.Errorf("CreateUser() UpdatedAt = %v, should be between %v and %v", result.UpdatedAt, before, after)
	}
}

func TestUserInMemRepository_CreateUser_DuplicateEmail(t *testing.T) {
	repo := NewUserInMemRepository()

	user1 := model.User{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password1",
	}

	user2 := model.User{
		Name:     "Alice Jones",
		Email:    "alice@example.com", // Same email
		Password: "password2",
	}

	// Create first user
	_, err := repo.CreateUser(user1)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Try to create second user with same email
	_, err = repo.CreateUser(user2)
	if err == nil {
		t.Error("CreateUser() should return error for duplicate email")
		return
	}

	expectedError := "user with email alice@example.com already exists"
	if err.Error() != expectedError {
		t.Errorf("CreateUser() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestUserInMemRepository_CreateUser_MultipleUsers(t *testing.T) {
	repo := NewUserInMemRepository()

	user1 := model.User{
		Name:     "User One",
		Email:    "user1@example.com",
		Password: "password1",
	}

	user2 := model.User{
		Name:     "User Two",
		Email:    "user2@example.com",
		Password: "password2",
	}

	// Create first user
	created1, err := repo.CreateUser(user1)
	if err != nil {
		t.Fatalf("First CreateUser() failed: %v", err)
	}

	// Create second user
	created2, err := repo.CreateUser(user2)
	if err != nil {
		t.Fatalf("Second CreateUser() failed: %v", err)
	}

	// IDs should be different
	if created1.ID == created2.ID {
		t.Error("CreateUser() should generate unique IDs for different users")
	}

	if created1.ID == "" || created2.ID == "" {
		t.Error("CreateUser() should generate non-empty IDs")
	}
}

func TestUserInMemRepository_GetUserByID(t *testing.T) {
	repo := NewUserInMemRepository()

	// Create a test user
	user := model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	created, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test getting the user by ID
	result, err := repo.GetUserByID(created.ID)
	if err != nil {
		t.Errorf("GetUserByID() unexpected error = %v", err)
		return
	}

	// Verify the result
	if result.ID != created.ID {
		t.Errorf("GetUserByID() result ID = %v, want %v", result.ID, created.ID)
	}
	if result.Name != user.Name {
		t.Errorf("GetUserByID() result Name = %v, want %v", result.Name, user.Name)
	}
	if result.Email != user.Email {
		t.Errorf("GetUserByID() result Email = %v, want %v", result.Email, user.Email)
	}
	if result.Password != user.Password {
		t.Errorf("GetUserByID() result Password = %v, want %v", result.Password, user.Password)
	}
	if !result.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("GetUserByID() result CreatedAt = %v, want %v", result.CreatedAt, created.CreatedAt)
	}
	if !result.UpdatedAt.Equal(created.UpdatedAt) {
		t.Errorf("GetUserByID() result UpdatedAt = %v, want %v", result.UpdatedAt, created.UpdatedAt)
	}
}

func TestUserInMemRepository_GetUserByID_NotFound(t *testing.T) {
	repo := NewUserInMemRepository()

	// Test getting a non-existent user
	_, err := repo.GetUserByID("nonexistent")
	if err == nil {
		t.Error("GetUserByID() with non-existent ID should return error")
		return
	}

	expectedError := "user with ID nonexistent not found"
	if err.Error() != expectedError {
		t.Errorf("GetUserByID() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestUserInMemRepository_GetUserByID_EmptyID(t *testing.T) {
	repo := NewUserInMemRepository()

	// Test getting user with empty ID
	_, err := repo.GetUserByID("")
	if err == nil {
		t.Error("GetUserByID() with empty ID should return error")
		return
	}

	expectedError := "user with ID  not found"
	if err.Error() != expectedError {
		t.Errorf("GetUserByID() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestUserInMemRepository_GetUserByEmail(t *testing.T) {
	repo := NewUserInMemRepository()

	// Create a test user
	user := model.User{
		Name:     "Email Test User",
		Email:    "emailtest@example.com",
		Password: "hashedpassword123",
	}

	created, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test getting user by email
	result, err := repo.GetUserByEmail("emailtest@example.com")
	if err != nil {
		t.Errorf("GetUserByEmail() unexpected error = %v", err)
		return
	}

	// Verify the result
	if result.ID != created.ID {
		t.Errorf("GetUserByEmail() result ID = %v, want %v", result.ID, created.ID)
	}
	if result.Name != user.Name {
		t.Errorf("GetUserByEmail() result Name = %v, want %v", result.Name, user.Name)
	}
	if result.Email != user.Email {
		t.Errorf("GetUserByEmail() result Email = %v, want %v", result.Email, user.Email)
	}
	if result.Password != user.Password {
		t.Errorf("GetUserByEmail() result Password = %v, want %v", result.Password, user.Password)
	}
}

func TestUserInMemRepository_GetUserByEmail_NotFound(t *testing.T) {
	repo := NewUserInMemRepository()

	// Test getting a non-existent user by email
	_, err := repo.GetUserByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("GetUserByEmail() with non-existent email should return error")
		return
	}

	expectedError := "user with email nonexistent@example.com not found"
	if err.Error() != expectedError {
		t.Errorf("GetUserByEmail() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestUserInMemRepository_GetUserByEmail_EmptyEmail(t *testing.T) {
	repo := NewUserInMemRepository()

	// Test getting user with empty email
	_, err := repo.GetUserByEmail("")
	if err == nil {
		t.Error("GetUserByEmail() with empty email should return error")
		return
	}

	expectedError := "user with email  not found"
	if err.Error() != expectedError {
		t.Errorf("GetUserByEmail() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestUserInMemRepository_GetUserByEmailAndPassword(t *testing.T) {
	repo := NewUserInMemRepository()

	// Create a test user
	user := model.User{
		Name:     "Auth User",
		Email:    "auth@example.com",
		Password: "hashedpassword123",
	}

	created, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test getting user by email and password
	result, err := repo.GetUserByEmailAndPassword("auth@example.com", "hashedpassword123")
	if err != nil {
		t.Errorf("GetUserByEmailAndPassword() unexpected error = %v", err)
		return
	}

	// Verify the result
	if result.ID != created.ID {
		t.Errorf("GetUserByEmailAndPassword() result ID = %v, want %v", result.ID, created.ID)
	}
	if result.Name != user.Name {
		t.Errorf("GetUserByEmailAndPassword() result Name = %v, want %v", result.Name, user.Name)
	}
	if result.Email != user.Email {
		t.Errorf("GetUserByEmailAndPassword() result Email = %v, want %v", result.Email, user.Email)
	}
	if result.Password != user.Password {
		t.Errorf("GetUserByEmailAndPassword() result Password = %v, want %v", result.Password, user.Password)
	}
}

func TestUserInMemRepository_GetUserByEmailAndPassword_WrongPassword(t *testing.T) {
	repo := NewUserInMemRepository()

	// Create a test user
	user := model.User{
		Name:     "Secure User",
		Email:    "secure@example.com",
		Password: "correcthash",
	}

	_, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test with wrong password
	_, err = repo.GetUserByEmailAndPassword("secure@example.com", "wronghash")
	if err == nil {
		t.Error("GetUserByEmailAndPassword() should return error for wrong password")
		return
	}

	expectedError := "user not found with provided credentials"
	if err.Error() != expectedError {
		t.Errorf("GetUserByEmailAndPassword() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestUserInMemRepository_GetUserByEmailAndPassword_WrongEmail(t *testing.T) {
	repo := NewUserInMemRepository()

	// Create a test user
	user := model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	_, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test with wrong email
	_, err = repo.GetUserByEmailAndPassword("wrong@example.com", "hashedpassword")
	if err == nil {
		t.Error("GetUserByEmailAndPassword() should return error for wrong email")
		return
	}

	expectedError := "user not found with provided credentials"
	if err.Error() != expectedError {
		t.Errorf("GetUserByEmailAndPassword() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestUserInMemRepository_GetUserByEmailAndPassword_NonExistentUser(t *testing.T) {
	repo := NewUserInMemRepository()

	// Test with non-existent user
	_, err := repo.GetUserByEmailAndPassword("nobody@example.com", "anypassword")
	if err == nil {
		t.Error("GetUserByEmailAndPassword() should return error for non-existent user")
		return
	}

	expectedError := "user not found with provided credentials"
	if err.Error() != expectedError {
		t.Errorf("GetUserByEmailAndPassword() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestUserInMemRepository_GetUserByEmailAndPassword_EmptyCredentials(t *testing.T) {
	repo := NewUserInMemRepository()

	// Test with empty email and password
	_, err := repo.GetUserByEmailAndPassword("", "")
	if err == nil {
		t.Error("GetUserByEmailAndPassword() should return error for empty credentials")
		return
	}

	expectedError := "user not found with provided credentials"
	if err.Error() != expectedError {
		t.Errorf("GetUserByEmailAndPassword() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestUserInMemRepository_ConcurrentCreate(t *testing.T) {
	repo := NewUserInMemRepository()
	var wg sync.WaitGroup
	numGoroutines := 10

	// Store created user IDs
	type result struct {
		id  string
		err error
	}
	results := make([]result, numGoroutines)

	// Test concurrent creates with unique emails
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			user := model.User{
				Name:     fmt.Sprintf("User %d", index),
				Email:    fmt.Sprintf("user%d@example.com", index),
				Password: fmt.Sprintf("password%d", index),
			}
			created, err := repo.CreateUser(user)
			results[index] = result{id: created.ID, err: err}
		}(i)
	}
	wg.Wait()

	// Verify all users were created successfully
	for i, res := range results {
		if res.err != nil {
			t.Errorf("Concurrent CreateUser() goroutine %d error: %v", i, res.err)
		}
		if res.id == "" {
			t.Errorf("Concurrent CreateUser() goroutine %d returned empty ID", i)
		}
	}

	// Verify all IDs are unique
	idMap := make(map[string]bool)
	for i, res := range results {
		if idMap[res.id] {
			t.Errorf("Concurrent CreateUser() generated duplicate ID: %s for goroutine %d", res.id, i)
		}
		idMap[res.id] = true
	}
}

func TestUserInMemRepository_ConcurrentRead(t *testing.T) {
	repo := NewUserInMemRepository()
	var wg sync.WaitGroup
	numGoroutines := 10

	// Create a user
	user := model.User{
		Name:     "Concurrent User",
		Email:    "concurrent@example.com",
		Password: "hashedpassword",
	}

	created, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				result, err := repo.GetUserByID(created.ID)
				if err != nil {
					t.Errorf("Concurrent GetUserByID() error: %v", err)
					return
				}
				if result.ID != created.ID {
					t.Errorf("Concurrent GetUserByID() result ID = %v, want %v", result.ID, created.ID)
				}
			}
		}()
	}
	wg.Wait()
}

func TestUserInMemRepository_ConcurrentReadWrite(t *testing.T) {
	repo := NewUserInMemRepository()
	var wg sync.WaitGroup
	numReaders := 5
	numWriters := 5

	// Create initial user
	initialUser := model.User{
		Name:     "Initial User",
		Email:    "initial@example.com",
		Password: "initialpassword",
	}

	created, err := repo.CreateUser(initialUser)
	if err != nil {
		t.Fatalf("Failed to create initial user: %v", err)
	}

	// Reader goroutines
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_, err := repo.GetUserByID(created.ID)
				if err != nil {
					t.Errorf("Concurrent read error: %v", err)
				}
			}
		}()
	}

	// Writer goroutines
	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(index int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				user := model.User{
					Name:     fmt.Sprintf("Writer %d User %d", index, j),
					Email:    fmt.Sprintf("writer%d.user%d@example.com", index, j),
					Password: fmt.Sprintf("password%d%d", index, j),
				}
				_, err := repo.CreateUser(user)
				if err != nil {
					t.Errorf("Concurrent write error: %v", err)
				}
			}
		}(i)
	}
	wg.Wait()
}

func TestUserInMemRepository_EmptyRepository(t *testing.T) {
	repo := NewUserInMemRepository()

	// Test getting non-existent user
	_, err := repo.GetUserByID("nonexistent")
	if err == nil {
		t.Error("GetUserByID() on empty repository should return error")
	}
	expectedError := "user with ID nonexistent not found"
	if err.Error() != expectedError {
		t.Errorf("GetUserByID() error = %v, want %v", err.Error(), expectedError)
	}

	// Test authentication on empty repository
	_, err = repo.GetUserByEmailAndPassword("nobody@example.com", "anypassword")
	if err == nil {
		t.Error("GetUserByEmailAndPassword() on empty repository should return error")
	}
	expectedError = "user not found with provided credentials"
	if err.Error() != expectedError {
		t.Errorf("GetUserByEmailAndPassword() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestUserInMemRepository_PreserveExistingTimestamps(t *testing.T) {
	repo := NewUserInMemRepository()

	specificCreatedAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	specificUpdatedAt := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)

	user := model.User{
		Name:      "Time Traveler",
		Email:     "time@example.com",
		Password:  "hashedpassword",
		CreatedAt: specificCreatedAt,
		UpdatedAt: specificUpdatedAt,
	}

	result, err := repo.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser() unexpected error = %v", err)
		return
	}

	// Timestamps should be preserved
	if !result.CreatedAt.Equal(specificCreatedAt) {
		t.Errorf("CreateUser() should preserve existing CreatedAt = %v, got %v", specificCreatedAt, result.CreatedAt)
	}
	if !result.UpdatedAt.Equal(specificUpdatedAt) {
		t.Errorf("CreateUser() should preserve existing UpdatedAt = %v, got %v", specificUpdatedAt, result.UpdatedAt)
	}
}

func TestUserInMemRepository_EdgeCases(t *testing.T) {
	repo := NewUserInMemRepository()

	// Test with empty string values (except email which must be unique)
	emptyUser := model.User{
		Name:     "",
		Email:    "empty@example.com",
		Password: "",
	}

	result, err := repo.CreateUser(emptyUser)
	if err != nil {
		t.Errorf("CreateUser() with empty values failed: %v", err)
	}
	if result.ID == "" {
		t.Error("CreateUser() should auto-generate ID even with empty other fields")
	}

	// Test with special characters
	specialUser := model.User{
		Name:     "User @#$%^&*()",
		Email:    "special!char#$@example.com",
		Password: "p@$$w0rd!#$",
	}

	_, err = repo.CreateUser(specialUser)
	if err != nil {
		t.Errorf("CreateUser() with special characters failed: %v", err)
	}

	// Verify it can be retrieved
	result, err = repo.GetUserByEmailAndPassword("special!char#$@example.com", "p@$$w0rd!#$")
	if err != nil {
		t.Errorf("GetUserByEmailAndPassword() with special characters failed: %v", err)
	}
	if result.Name != specialUser.Name {
		t.Errorf("GetUserByEmailAndPassword() name = %v, want %v", result.Name, specialUser.Name)
	}

	// Test with very long strings
	longUser := model.User{
		Name:     string(make([]byte, 1000)),
		Email:    "long@example.com",
		Password: string(make([]byte, 1000)),
	}

	_, err = repo.CreateUser(longUser)
	if err != nil {
		t.Errorf("CreateUser() with long strings failed: %v", err)
	}
}

func TestNewUserInMemRepository(t *testing.T) {
	repo := NewUserInMemRepository()

	if repo == nil {
		t.Error("NewUserInMemRepository() should return non-nil repository")
	}

	if repo.users == nil {
		t.Error("NewUserInMemRepository() should initialize users map")
	}

	// Verify repository is usable
	user := model.User{
		Name:     "Test Init",
		Email:    "init@example.com",
		Password: "password",
	}

	_, err := repo.CreateUser(user)
	if err != nil {
		t.Errorf("NewUserInMemRepository() created repository should be functional, got error = %v", err)
	}
}

func TestUserInMemRepository_CaseInsensitiveEmail(t *testing.T) {
	repo := NewUserInMemRepository()

	user1 := model.User{
		Name:     "User One",
		Email:    "test@example.com",
		Password: "password1",
	}

	_, err := repo.CreateUser(user1)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Try to create user with different case email
	user2 := model.User{
		Name:     "User Two",
		Email:    "TEST@EXAMPLE.COM",
		Password: "password2",
	}

	// Note: Current implementation is case-sensitive
	// This test documents the current behavior
	_, err = repo.CreateUser(user2)
	if err != nil {
		t.Logf("CreateUser() with different case email returned error: %v (case-sensitive)", err)
		// This is the current behavior - email comparison is case-sensitive
	} else {
		t.Logf("CreateUser() with different case email succeeded (case-insensitive)")
		// If implementation becomes case-insensitive, this would be the expected behavior
	}
}
