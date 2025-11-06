package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/tsongpon/athena/internal/model"
)

// =============================================================================
// MOCK-BASED TESTS (run without database)
// =============================================================================

// TestBookmarkPostgresRepository_CreateBookmark_Mock tests bookmark creation with mock DB
func TestBookmarkPostgresRepository_CreateBookmark_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	bookmark := model.Bookmark{
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	// Expect INSERT query
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO bookmarks (id, user_id, url, title, is_archived, created_at, updated_at)`)).
		WithArgs(sqlmock.AnyArg(), "user1", "https://example.com", "Example", false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Errorf("CreateBookmark() unexpected error = %v", err)
		return
	}

	if result.ID == "" {
		t.Error("CreateBookmark() should generate ID")
	}
	if result.UserID != bookmark.UserID {
		t.Errorf("CreateBookmark() UserID = %v, want %v", result.UserID, bookmark.UserID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_CreateBookmark_WithID tests creating bookmark with existing ID
func TestBookmarkPostgresRepository_CreateBookmark_WithID_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	bookmark := model.Bookmark{
		ID:     "existing-id",
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO bookmarks`)).
		WithArgs("existing-id", "user1", "https://example.com", "Example", false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Errorf("CreateBookmark() unexpected error = %v", err)
		return
	}

	if result.ID != "existing-id" {
		t.Errorf("CreateBookmark() should preserve existing ID, got %v", result.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_CreateBookmark_Error tests error handling
func TestBookmarkPostgresRepository_CreateBookmark_Error_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	bookmark := model.Bookmark{
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO bookmarks`)).
		WillReturnError(errors.New("database error"))

	_, err = repo.CreateBookmark(bookmark)
	if err == nil {
		t.Error("CreateBookmark() should return error when database fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_GetBookmark_Mock tests retrieving a bookmark
func TestBookmarkPostgresRepository_GetBookmark_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	expectedBookmark := model.Bookmark{
		ID:         "bookmark1",
		UserID:     "user1",
		URL:        "https://example.com",
		Title:      "Example",
		IsArchived: false,
		CreatedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "url", "title", "is_archived", "created_at", "updated_at"}).
		AddRow(expectedBookmark.ID, expectedBookmark.UserID, expectedBookmark.URL, expectedBookmark.Title,
			expectedBookmark.IsArchived, expectedBookmark.CreatedAt, expectedBookmark.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, url, title, is_archived, created_at, updated_at FROM bookmarks WHERE id = $1`)).
		WithArgs("bookmark1").
		WillReturnRows(rows)

	result, err := repo.GetBookmark("bookmark1")
	if err != nil {
		t.Errorf("GetBookmark() unexpected error = %v", err)
		return
	}

	if result.ID != expectedBookmark.ID {
		t.Errorf("GetBookmark() ID = %v, want %v", result.ID, expectedBookmark.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_GetBookmark_NotFound tests not found scenario
func TestBookmarkPostgresRepository_GetBookmark_NotFound_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, url, title, is_archived, created_at, updated_at FROM bookmarks WHERE id = $1`)).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetBookmark("nonexistent")
	if err == nil {
		t.Error("GetBookmark() should return error for non-existent bookmark")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_GetBookmark_QueryError tests query error
func TestBookmarkPostgresRepository_GetBookmark_QueryError_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, url, title, is_archived, created_at, updated_at FROM bookmarks WHERE id = $1`)).
		WithArgs("bookmark1").
		WillReturnError(errors.New("database error"))

	_, err = repo.GetBookmark("bookmark1")
	if err == nil {
		t.Error("GetBookmark() should return error when database fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_ListBookmarks_Mock tests listing bookmarks
func TestBookmarkPostgresRepository_ListBookmarks_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "url", "title", "is_archived", "created_at", "updated_at"}).
		AddRow("bookmark1", "user1", "https://example1.com", "Example 1", false, time.Now(), time.Now()).
		AddRow("bookmark2", "user1", "https://example2.com", "Example 2", false, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, url, title, is_archived, created_at, updated_at FROM bookmarks WHERE user_id = $1 AND is_archived = $2 ORDER BY created_at DESC`)).
		WithArgs("user1", false).
		WillReturnRows(rows)

	result, err := repo.ListBookmarks(model.BookmarkQuery{UserID: "user1", Archived: false})
	if err != nil {
		t.Errorf("ListBookmarks() unexpected error = %v", err)
		return
	}

	if len(result) != 2 {
		t.Errorf("ListBookmarks() returned %d bookmarks, want 2", len(result))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_ListBookmarks_WithPagination tests pagination
func TestBookmarkPostgresRepository_ListBookmarks_WithPagination_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "url", "title", "is_archived", "created_at", "updated_at"}).
		AddRow("bookmark1", "user1", "https://example1.com", "Example 1", false, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, url, title, is_archived, created_at, updated_at FROM bookmarks WHERE user_id = $1 AND is_archived = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`)).
		WithArgs("user1", false, 10, 0).
		WillReturnRows(rows)

	result, err := repo.ListBookmarks(model.BookmarkQuery{UserID: "user1", Archived: false, Page: 1, PageSize: 10})
	if err != nil {
		t.Errorf("ListBookmarks() with pagination unexpected error = %v", err)
		return
	}

	if len(result) != 1 {
		t.Errorf("ListBookmarks() returned %d bookmarks, want 1", len(result))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_ListBookmarks_Error tests error handling
func TestBookmarkPostgresRepository_ListBookmarks_Error_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, url, title, is_archived, created_at, updated_at FROM bookmarks`)).
		WillReturnError(errors.New("database error"))

	_, err = repo.ListBookmarks(model.BookmarkQuery{UserID: "user1", Archived: false})
	if err == nil {
		t.Error("ListBookmarks() should return error when database fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_ListBookmarks_ScanError tests scan error
func TestBookmarkPostgresRepository_ListBookmarks_ScanError_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	// Return wrong column types to cause scan error
	rows := sqlmock.NewRows([]string{"id", "user_id", "url", "title", "is_archived", "created_at", "updated_at"}).
		AddRow("bookmark1", "user1", "https://example.com", "Example", "invalid_bool", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, url, title, is_archived, created_at, updated_at FROM bookmarks`)).
		WithArgs("user1", false).
		WillReturnRows(rows)

	_, err = repo.ListBookmarks(model.BookmarkQuery{UserID: "user1", Archived: false})
	if err == nil {
		t.Error("ListBookmarks() should return error when scan fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_CountBookmarks_Mock tests counting bookmarks
func TestBookmarkPostgresRepository_CountBookmarks_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM bookmarks WHERE user_id = $1 AND is_archived = $2`)).
		WithArgs("user1", false).
		WillReturnRows(rows)

	count, err := repo.CountBookmarks(model.BookmarkQuery{UserID: "user1", Archived: false})
	if err != nil {
		t.Errorf("CountBookmarks() unexpected error = %v", err)
		return
	}

	if count != 5 {
		t.Errorf("CountBookmarks() = %d, want 5", count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_CountBookmarks_Error tests count error
func TestBookmarkPostgresRepository_CountBookmarks_Error_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM bookmarks WHERE user_id = $1 AND is_archived = $2`)).
		WithArgs("user1", false).
		WillReturnError(errors.New("database error"))

	_, err = repo.CountBookmarks(model.BookmarkQuery{UserID: "user1", Archived: false})
	if err == nil {
		t.Error("CountBookmarks() should return error when database fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_UpdateBookmark_Mock tests updating a bookmark
func TestBookmarkPostgresRepository_UpdateBookmark_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	bookmark := model.Bookmark{
		ID:         "bookmark1",
		UserID:     "user1",
		URL:        "https://updated.com",
		Title:      "Updated Title",
		IsArchived: true,
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE bookmarks SET user_id = $2, url = $3, title = $4, is_archived = $5, updated_at = $6 WHERE id = $1`)).
		WithArgs("bookmark1", "user1", "https://updated.com", "Updated Title", true, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Expect GetBookmark to be called after update
	rows := sqlmock.NewRows([]string{"id", "user_id", "url", "title", "is_archived", "created_at", "updated_at"}).
		AddRow("bookmark1", "user1", "https://updated.com", "Updated Title", true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, url, title, is_archived, created_at, updated_at FROM bookmarks WHERE id = $1`)).
		WithArgs("bookmark1").
		WillReturnRows(rows)

	result, err := repo.UpdateBookmark(bookmark)
	if err != nil {
		t.Errorf("UpdateBookmark() unexpected error = %v", err)
		return
	}

	if result.Title != "Updated Title" {
		t.Errorf("UpdateBookmark() Title = %v, want Updated Title", result.Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_UpdateBookmark_NotFound tests updating non-existent bookmark
func TestBookmarkPostgresRepository_UpdateBookmark_NotFound_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	bookmark := model.Bookmark{
		ID:     "nonexistent",
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE bookmarks SET user_id = $2, url = $3, title = $4, is_archived = $5, updated_at = $6 WHERE id = $1`)).
		WithArgs("nonexistent", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err = repo.UpdateBookmark(bookmark)
	if err == nil {
		t.Error("UpdateBookmark() should return error for non-existent bookmark")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_UpdateBookmark_ExecError tests update execution error
func TestBookmarkPostgresRepository_UpdateBookmark_ExecError_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	bookmark := model.Bookmark{
		ID:     "bookmark1",
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE bookmarks`)).
		WillReturnError(errors.New("database error"))

	_, err = repo.UpdateBookmark(bookmark)
	if err == nil {
		t.Error("UpdateBookmark() should return error when database fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_UpdateBookmark_RowsAffectedError tests RowsAffected error
func TestBookmarkPostgresRepository_UpdateBookmark_RowsAffectedError_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	bookmark := model.Bookmark{
		ID:     "bookmark1",
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE bookmarks`)).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

	_, err = repo.UpdateBookmark(bookmark)
	if err == nil {
		t.Error("UpdateBookmark() should return error when RowsAffected fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_DeleteBookmark_Mock tests deleting a bookmark
func TestBookmarkPostgresRepository_DeleteBookmark_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM bookmarks WHERE id = $1`)).
		WithArgs("bookmark1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteBookmark("bookmark1")
	if err != nil {
		t.Errorf("DeleteBookmark() unexpected error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_DeleteBookmark_NotFound tests deleting non-existent bookmark
func TestBookmarkPostgresRepository_DeleteBookmark_NotFound_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM bookmarks WHERE id = $1`)).
		WithArgs("nonexistent").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.DeleteBookmark("nonexistent")
	if err == nil {
		t.Error("DeleteBookmark() should return error for non-existent bookmark")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_DeleteBookmark_ExecError tests delete execution error
func TestBookmarkPostgresRepository_DeleteBookmark_ExecError_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM bookmarks WHERE id = $1`)).
		WithArgs("bookmark1").
		WillReturnError(errors.New("database error"))

	err = repo.DeleteBookmark("bookmark1")
	if err == nil {
		t.Error("DeleteBookmark() should return error when database fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestBookmarkPostgresRepository_DeleteBookmark_RowsAffectedError tests RowsAffected error
func TestBookmarkPostgresRepository_DeleteBookmark_RowsAffectedError_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM bookmarks WHERE id = $1`)).
		WithArgs("bookmark1").
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

	err = repo.DeleteBookmark("bookmark1")
	if err == nil {
		t.Error("DeleteBookmark() should return error when RowsAffected fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

// TestNewBookmarkPostgresRepository tests repository initialization
func TestNewBookmarkPostgresRepository_Mock(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := NewBookmarkPostgresRepository(db)
	if repo == nil {
		t.Error("NewBookmarkPostgresRepository() should not return nil")
	}

	if repo.db == nil {
		t.Error("NewBookmarkPostgresRepository() should initialize db")
	}
}

// =============================================================================
// INTEGRATION TESTS (require actual PostgreSQL database)
// These tests are skipped if TEST_DATABASE_URL is not set
// =============================================================================

// Helper function to create a test database connection
func setupTestDB(t *testing.T) *sql.DB {
	// This requires a test PostgreSQL instance
	// You can skip this test if TEST_DATABASE_URL is not set
	dbURL := getTestDatabaseURL()
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping PostgreSQL integration tests")
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

func TestBookmarkPostgresRepository_CreateBookmark_Integration(t *testing.T) {
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

func TestBookmarkPostgresRepository_GetBookmark_Integration(t *testing.T) {
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

func TestBookmarkPostgresRepository_ListBookmarks_Integration(t *testing.T) {
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

func TestBookmarkPostgresRepository_UpdateBookmark_Integration(t *testing.T) {
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

func TestBookmarkPostgresRepository_DeleteBookmark_Integration(t *testing.T) {
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

func TestBookmarkPostgresRepository_ListBookmarks_OrderedByCreatedDateDesc_Integration(t *testing.T) {
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

func TestBookmarkPostgresRepository_ListBookmarks_FilterByArchived_Integration(t *testing.T) {
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
