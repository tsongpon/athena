package repository

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/tsongpon/athena/internal/model"
)

// Helper function to create a test database connection
func setupTestDB(t *testing.T) *sql.DB {
	// This requires a test PostgreSQL instance
	// You can skip this test if TEST_DATABASE_URL is not set
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping PostgreSQL tests")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	return db
}

// Helper function to clean up test data
func cleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM bookmarks")
	if err != nil {
		t.Logf("Warning: Failed to clean up test data: %v", err)
	}
}

// Helper function to get test database URL from environment
func getTestDatabaseURL() string {
	// This can be set via environment variable for CI/CD
	// Example: TEST_DATABASE_URL=postgres://user:pass@localhost/athena_test?sslmode=disable
	return ""
}

func TestBookmarkPostgresRepository_CreateBookmark(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := NewBookmarkPostgresRepository(db)

	bookmark := model.Bookmark{
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	result, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Fatalf("CreateBookmark() unexpected error = %v", err)
	}

	// ID should be auto-generated
	if result.ID == "" {
		t.Error("CreateBookmark() result ID should not be empty")
	}
	if result.UserID != bookmark.UserID {
		t.Errorf("CreateBookmark() result UserID = %v, want %v", result.UserID, bookmark.UserID)
	}
	if result.URL != bookmark.URL {
		t.Errorf("CreateBookmark() result URL = %v, want %v", result.URL, bookmark.URL)
	}
	if result.Title != bookmark.Title {
		t.Errorf("CreateBookmark() result Title = %v, want %v", result.Title, bookmark.Title)
	}
	if result.CreatedAt.IsZero() {
		t.Error("CreateBookmark() result CreatedAt should not be zero")
	}
	if result.UpdatedAt.IsZero() {
		t.Error("CreateBookmark() result UpdatedAt should not be zero")
	}
}

func TestBookmarkPostgresRepository_GetBookmark(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := NewBookmarkPostgresRepository(db)

	// Create a test bookmark
	bookmark := model.Bookmark{
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	created, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Fatalf("Failed to create bookmark: %v", err)
	}

	// Get the bookmark
	result, err := repo.GetBookmark(created.ID)
	if err != nil {
		t.Fatalf("GetBookmark() unexpected error = %v", err)
	}

	if result.ID != created.ID {
		t.Errorf("GetBookmark() result ID = %v, want %v", result.ID, created.ID)
	}
	if result.UserID != bookmark.UserID {
		t.Errorf("GetBookmark() result UserID = %v, want %v", result.UserID, bookmark.UserID)
	}
}

func TestBookmarkPostgresRepository_GetBookmark_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := NewBookmarkPostgresRepository(db)

	_, err := repo.GetBookmark("nonexistent")
	if err == nil {
		t.Error("GetBookmark() with non-existent ID should return error")
	}
}

func TestBookmarkPostgresRepository_ListBookmarks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := NewBookmarkPostgresRepository(db)

	// Create test bookmarks
	bookmarks := []model.Bookmark{
		{UserID: "user1", URL: "https://example1.com", Title: "Example 1"},
		{UserID: "user1", URL: "https://example2.com", Title: "Example 2"},
		{UserID: "user2", URL: "https://example3.com", Title: "Example 3"},
	}

	for _, bookmark := range bookmarks {
		_, err := repo.CreateBookmark(bookmark)
		if err != nil {
			t.Fatalf("Failed to create bookmark: %v", err)
		}
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	// List bookmarks for user1
	result, err := repo.ListBookmarks(model.BookmarkQuery{UserID: "user1", Archived: false})
	if err != nil {
		t.Fatalf("ListBookmarks() unexpected error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("ListBookmarks() count = %v, want 2", len(result))
	}

	// Verify they are sorted by created date descending (newest first)
	if len(result) >= 2 {
		if result[0].CreatedAt.Before(result[1].CreatedAt) {
			t.Error("ListBookmarks() should return bookmarks in descending order by created date")
		}
	}
}

func TestBookmarkPostgresRepository_UpdateBookmark(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := NewBookmarkPostgresRepository(db)

	// Create initial bookmark
	bookmark := model.Bookmark{
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Original Title",
	}

	created, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Fatalf("Failed to create bookmark: %v", err)
	}

	// Update the bookmark
	created.Title = "Updated Title"
	created.IsArchived = true

	result, err := repo.UpdateBookmark(created)
	if err != nil {
		t.Fatalf("UpdateBookmark() unexpected error = %v", err)
	}

	if result.Title != "Updated Title" {
		t.Errorf("UpdateBookmark() result Title = %v, want Updated Title", result.Title)
	}
	if !result.IsArchived {
		t.Error("UpdateBookmark() result IsArchived should be true")
	}
	if result.UpdatedAt.Before(created.CreatedAt) {
		t.Error("UpdateBookmark() UpdatedAt should be after CreatedAt")
	}
}

func TestBookmarkPostgresRepository_UpdateBookmark_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := NewBookmarkPostgresRepository(db)

	bookmark := model.Bookmark{
		ID:     "nonexistent",
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	_, err := repo.UpdateBookmark(bookmark)
	if err == nil {
		t.Error("UpdateBookmark() with non-existent ID should return error")
	}
}

func TestBookmarkPostgresRepository_DeleteBookmark(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := NewBookmarkPostgresRepository(db)

	// Create a test bookmark
	bookmark := model.Bookmark{
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	created, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Fatalf("Failed to create bookmark: %v", err)
	}

	// Delete the bookmark
	err = repo.DeleteBookmark(created.ID)
	if err != nil {
		t.Fatalf("DeleteBookmark() unexpected error = %v", err)
	}

	// Verify bookmark is deleted
	_, err = repo.GetBookmark(created.ID)
	if err == nil {
		t.Error("GetBookmark() after delete should return error")
	}
}

func TestBookmarkPostgresRepository_DeleteBookmark_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := NewBookmarkPostgresRepository(db)

	err := repo.DeleteBookmark("nonexistent")
	if err == nil {
		t.Error("DeleteBookmark() with non-existent ID should return error")
	}
}

func TestBookmarkPostgresRepository_ListBookmarks_OrderedByCreatedDateDesc(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := NewBookmarkPostgresRepository(db)

	// Create bookmarks with specific timestamps
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	bookmarks := []model.Bookmark{
		{
			UserID:    "user1",
			URL:       "https://example1.com",
			Title:     "First (oldest)",
			CreatedAt: baseTime,
		},
		{
			UserID:    "user1",
			URL:       "https://example2.com",
			Title:     "Second",
			CreatedAt: baseTime.Add(1 * time.Hour),
		},
		{
			UserID:    "user1",
			URL:       "https://example3.com",
			Title:     "Third (newest)",
			CreatedAt: baseTime.Add(2 * time.Hour),
		},
	}

	// Add bookmarks in random order
	for _, bookmark := range []int{1, 0, 2} {
		_, err := repo.CreateBookmark(bookmarks[bookmark])
		if err != nil {
			t.Fatalf("Failed to create bookmark: %v", err)
		}
	}

	// List bookmarks
	result, err := repo.ListBookmarks(model.BookmarkQuery{UserID: "user1", Archived: false})
	if err != nil {
		t.Fatalf("ListBookmarks() unexpected error = %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("ListBookmarks() count = %v, want 3", len(result))
	}

	// Verify they are ordered by created date descending (newest first)
	if result[0].Title != "Third (newest)" {
		t.Errorf("First bookmark should be 'Third (newest)', got '%s'", result[0].Title)
	}
	if result[1].Title != "Second" {
		t.Errorf("Second bookmark should be 'Second', got '%s'", result[1].Title)
	}
	if result[2].Title != "First (oldest)" {
		t.Errorf("Third bookmark should be 'First (oldest)', got '%s'", result[2].Title)
	}

	// Verify timestamps are in descending order
	for i := 0; i < len(result)-1; i++ {
		if result[i].CreatedAt.Before(result[i+1].CreatedAt) {
			t.Errorf("Bookmarks not in descending order: result[%d].CreatedAt (%v) is before result[%d].CreatedAt (%v)",
				i, result[i].CreatedAt, i+1, result[i+1].CreatedAt)
		}
	}
}

func TestBookmarkPostgresRepository_ListBookmarks_FilterByArchived(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := NewBookmarkPostgresRepository(db)

	// Create bookmarks with different archived states
	bookmarks := []model.Bookmark{
		{UserID: "user1", URL: "https://example1.com", Title: "Active 1", IsArchived: false},
		{UserID: "user1", URL: "https://example2.com", Title: "Active 2", IsArchived: false},
		{UserID: "user1", URL: "https://example3.com", Title: "Archived 1", IsArchived: true},
	}

	for _, bookmark := range bookmarks {
		_, err := repo.CreateBookmark(bookmark)
		if err != nil {
			t.Fatalf("Failed to create bookmark: %v", err)
		}
	}

	// List active bookmarks
	activeResult, err := repo.ListBookmarks(model.BookmarkQuery{UserID: "user1", Archived: false})
	if err != nil {
		t.Fatalf("ListBookmarks() for active unexpected error = %v", err)
	}

	if len(activeResult) != 2 {
		t.Errorf("ListBookmarks() active count = %v, want 2", len(activeResult))
	}

	// List archived bookmarks
	archivedResult, err := repo.ListBookmarks(model.BookmarkQuery{UserID: "user1", Archived: true})
	if err != nil {
		t.Fatalf("ListBookmarks() for archived unexpected error = %v", err)
	}

	if len(archivedResult) != 1 {
		t.Errorf("ListBookmarks() archived count = %v, want 1", len(archivedResult))
	}
}
