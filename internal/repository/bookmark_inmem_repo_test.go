package repository

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/tsongpon/athena/internal/model"
)

func TestBookmarkInMemRepository_CreateBookmark(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	bookmark := model.Bookmark{
		UserID:    "user1",
		URL:       "https://example.com",
		Title:     "Example",
		CreatedAt: time.Now(),
	}

	result, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Errorf("CreateBookmark() unexpected error = %v", err)
		return
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
}

func TestBookmarkInMemRepository_GetBookmark(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Create a test bookmark
	bookmark := model.Bookmark{
		UserID:    "user1",
		URL:       "https://example.com",
		Title:     "Example",
		CreatedAt: time.Now(),
	}

	created, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Fatalf("Failed to create bookmark: %v", err)
	}

	// Test getting the bookmark using the auto-generated ID
	result, err := repo.GetBookmark(created.ID)
	if err != nil {
		t.Errorf("GetBookmark() unexpected error = %v", err)
		return
	}

	// Verify the result
	if result.ID != created.ID {
		t.Errorf("GetBookmark() result ID = %v, want %v", result.ID, created.ID)
	}
	if result.UserID != bookmark.UserID {
		t.Errorf("GetBookmark() result UserID = %v, want %v", result.UserID, bookmark.UserID)
	}
	if result.URL != bookmark.URL {
		t.Errorf("GetBookmark() result URL = %v, want %v", result.URL, bookmark.URL)
	}
	if result.Title != bookmark.Title {
		t.Errorf("GetBookmark() result Title = %v, want %v", result.Title, bookmark.Title)
	}
	if !result.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("GetBookmark() result CreatedAt = %v, want %v", result.CreatedAt, created.CreatedAt)
	}
}

func TestBookmarkInMemRepository_GetBookmark_NotFound(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Test getting a non-existent bookmark
	_, err := repo.GetBookmark("nonexistent")
	if err == nil {
		t.Error("GetBookmark() with non-existent ID should return error")
		return
	}

	expectedError := "bookmark with ID nonexistent not found"
	if err.Error() != expectedError {
		t.Errorf("GetBookmark() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestBookmarkInMemRepository_GetBookmark_AfterUpdate(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Create initial bookmark
	originalBookmark := model.Bookmark{
		UserID:    "user1",
		URL:       "https://example.com",
		Title:     "Original Title",
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	created, err := repo.CreateBookmark(originalBookmark)
	if err != nil {
		t.Fatalf("Failed to create bookmark: %v", err)
	}

	// Update the bookmark
	updatedBookmark := model.Bookmark{
		ID:     created.ID,
		UserID: "user1",
		URL:    "https://updated.com",
		Title:  "Updated Title",
	}

	_, err = repo.UpdateBookmark(updatedBookmark)
	if err != nil {
		t.Fatalf("Failed to update bookmark: %v", err)
	}

	// Get the bookmark and verify it returns the updated version
	result, err := repo.GetBookmark(created.ID)
	if err != nil {
		t.Errorf("GetBookmark() after update unexpected error = %v", err)
		return
	}

	if result.URL != updatedBookmark.URL {
		t.Errorf("GetBookmark() after update result URL = %v, want %v", result.URL, updatedBookmark.URL)
	}
	if result.Title != updatedBookmark.Title {
		t.Errorf("GetBookmark() after update result Title = %v, want %v", result.Title, updatedBookmark.Title)
	}
	// CreatedAt should still be the original value
	if !result.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("GetBookmark() after update CreatedAt = %v, want %v", result.CreatedAt, created.CreatedAt)
	}
}

func TestBookmarkInMemRepository_GetBookmark_AfterDelete(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Create a bookmark
	bookmark := model.Bookmark{
		UserID:    "user1",
		URL:       "https://example.com",
		Title:     "Example",
		CreatedAt: time.Now(),
	}

	created, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Fatalf("Failed to create bookmark: %v", err)
	}

	// Verify bookmark exists
	_, err = repo.GetBookmark(created.ID)
	if err != nil {
		t.Fatalf("GetBookmark() before delete failed: %v", err)
	}

	// Delete the bookmark
	err = repo.DeleteBookmark(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete bookmark: %v", err)
	}

	// Try to get the deleted bookmark
	_, err = repo.GetBookmark(created.ID)
	if err == nil {
		t.Error("GetBookmark() after delete should return error")
		return
	}

	expectedError := fmt.Sprintf("bookmark with ID %s not found", created.ID)
	if err.Error() != expectedError {
		t.Errorf("GetBookmark() after delete error = %v, want %v", err.Error(), expectedError)
	}
}

func TestBookmarkInMemRepository_GetBookmark_ConcurrentAccess(t *testing.T) {
	repo := NewBookmarkInMemRepository()
	var wg sync.WaitGroup
	numGoroutines := 10

	// Create a bookmark
	bookmark := model.Bookmark{
		UserID:    "user1",
		URL:       "https://example.com",
		Title:     "Example",
		CreatedAt: time.Now(),
	}

	created, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Fatalf("Failed to create bookmark: %v", err)
	}

	// Test concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				result, err := repo.GetBookmark(created.ID)
				if err != nil {
					t.Errorf("Concurrent GetBookmark() error: %v", err)
					return
				}
				if result.ID != created.ID {
					t.Errorf("Concurrent GetBookmark() result ID = %v, want %v", result.ID, created.ID)
				}
			}
		}()
	}
	wg.Wait()
}

func TestBookmarkInMemRepository_GetBookmark_WithEmptyID(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Test getting bookmark with empty ID (should fail since IDs are auto-generated)
	_, err := repo.GetBookmark("")
	if err == nil {
		t.Error("GetBookmark() with empty ID should return error")
		return
	}

	expectedError := "bookmark with ID  not found"
	if err.Error() != expectedError {
		t.Errorf("GetBookmark() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestBookmarkInMemRepository_CreateBookmark_MultipleBookmarks tests creating multiple bookmarks
// Since IDs are auto-generated, duplicate IDs are not possible
func TestBookmarkInMemRepository_CreateBookmark_MultipleBookmarks(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	bookmark1 := model.Bookmark{
		UserID: "user1",
		URL:    "https://example1.com",
		Title:  "Example 1",
	}

	bookmark2 := model.Bookmark{
		UserID: "user2",
		URL:    "https://example2.com",
		Title:  "Example 2",
	}

	// Create first bookmark
	created1, err := repo.CreateBookmark(bookmark1)
	if err != nil {
		t.Fatalf("First CreateBookmark() failed: %v", err)
	}

	// Create second bookmark
	created2, err := repo.CreateBookmark(bookmark2)
	if err != nil {
		t.Fatalf("Second CreateBookmark() failed: %v", err)
	}

	// IDs should be different (auto-generated)
	if created1.ID == created2.ID {
		t.Error("CreateBookmark() should generate unique IDs for different bookmarks")
	}

	if created1.ID == "" || created2.ID == "" {
		t.Error("CreateBookmark() should generate non-empty IDs")
	}
}

func TestBookmarkInMemRepository_ListBookmarks(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Create test bookmarks
	bookmarks := []model.Bookmark{
		{
			UserID:    "user1",
			URL:       "https://example1.com",
			Title:     "Example 1",
			CreatedAt: time.Now(),
		},
		{
			UserID:    "user1",
			URL:       "https://example2.com",
			Title:     "Example 2",
			CreatedAt: time.Now(),
		},
		{
			UserID:    "user2",
			URL:       "https://example3.com",
			Title:     "Example 3",
			CreatedAt: time.Now(),
		},
	}

	// Add bookmarks to repository
	for _, bookmark := range bookmarks {
		_, err := repo.CreateBookmark(bookmark)
		if err != nil {
			t.Fatalf("Failed to create bookmark: %v", err)
		}
	}

	// Test listing bookmarks for user1
	result, err := repo.ListBookmarks("user1")
	if err != nil {
		t.Errorf("ListBookmarks() unexpected error = %v", err)
		return
	}

	if len(result) != 2 {
		t.Errorf("ListBookmarks() count = %v, want 2", len(result))
		return
	}

	// Check if both user1 bookmarks are present
	for _, bookmark := range result {
		if bookmark.UserID != "user1" {
			t.Errorf("ListBookmarks() bookmark %s has wrong UserID = %v, want user1", bookmark.ID, bookmark.UserID)
		}
	}

	// Test listing bookmarks for user2
	result, err = repo.ListBookmarks("user2")
	if err != nil {
		t.Errorf("ListBookmarks() unexpected error = %v", err)
		return
	}

	if len(result) != 1 {
		t.Errorf("ListBookmarks() count = %v, want 1", len(result))
		return
	}

	if result[0].UserID != "user2" {
		t.Errorf("ListBookmarks() bookmark UserID = %v, want user2", result[0].UserID)
	}

	// Test listing bookmarks for non-existent user
	result, err = repo.ListBookmarks("nonexistent")
	if err != nil {
		t.Errorf("ListBookmarks() unexpected error = %v", err)
		return
	}

	if len(result) != 0 {
		t.Errorf("ListBookmarks() for non-existent user count = %v, want 0", len(result))
	}
}

func TestBookmarkInMemRepository_UpdateBookmark(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Create initial bookmark
	originalBookmark := model.Bookmark{
		UserID:    "user1",
		URL:       "https://example.com",
		Title:     "Original Title",
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	created, err := repo.CreateBookmark(originalBookmark)
	if err != nil {
		t.Fatalf("Failed to create initial bookmark: %v", err)
	}

	// Update the bookmark
	updatedBookmark := model.Bookmark{
		ID:     created.ID,
		UserID: "user1",
		URL:    "https://updated.com",
		Title:  "Updated Title",
	}

	result, err := repo.UpdateBookmark(updatedBookmark)
	if err != nil {
		t.Errorf("UpdateBookmark() unexpected error = %v", err)
		return
	}

	// Verify the result
	if result.ID != updatedBookmark.ID {
		t.Errorf("UpdateBookmark() result ID = %v, want %v", result.ID, updatedBookmark.ID)
	}
	if result.UserID != updatedBookmark.UserID {
		t.Errorf("UpdateBookmark() result UserID = %v, want %v", result.UserID, updatedBookmark.UserID)
	}
	if result.URL != updatedBookmark.URL {
		t.Errorf("UpdateBookmark() result URL = %v, want %v", result.URL, updatedBookmark.URL)
	}
	if result.Title != updatedBookmark.Title {
		t.Errorf("UpdateBookmark() result Title = %v, want %v", result.Title, updatedBookmark.Title)
	}

	// Verify CreatedAt is preserved from original
	if !result.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("UpdateBookmark() should preserve CreatedAt = %v, got %v", created.CreatedAt, result.CreatedAt)
	}

	// Test updating non-existent bookmark
	nonExistentBookmark := model.Bookmark{
		ID:     "999",
		UserID: "user1",
		URL:    "https://nonexistent.com",
		Title:  "Non-existent",
	}

	_, err = repo.UpdateBookmark(nonExistentBookmark)
	if err == nil {
		t.Error("UpdateBookmark() with non-existent ID should return error")
		return
	}

	expectedError := "bookmark with ID 999 not found"
	if err.Error() != expectedError {
		t.Errorf("UpdateBookmark() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestBookmarkInMemRepository_DeleteBookmark(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Create test bookmark
	bookmark := model.Bookmark{
		UserID:    "user1",
		URL:       "https://example.com",
		Title:     "Example",
		CreatedAt: time.Now(),
	}

	created, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Fatalf("Failed to create bookmark: %v", err)
	}

	// Delete the bookmark
	err = repo.DeleteBookmark(created.ID)
	if err != nil {
		t.Errorf("DeleteBookmark() unexpected error = %v", err)
		return
	}

	// Verify bookmark is actually deleted
	bookmarks, err := repo.ListBookmarks("user1")
	if err != nil {
		t.Errorf("ListBookmarks() after delete failed: %v", err)
		return
	}

	if len(bookmarks) != 0 {
		t.Errorf("DeleteBookmark() bookmark still exists after deletion")
	}

	// Test deleting non-existent bookmark
	err = repo.DeleteBookmark("999")
	if err == nil {
		t.Error("DeleteBookmark() with non-existent ID should return error")
		return
	}

	expectedError := "bookmark with ID 999 not found"
	if err.Error() != expectedError {
		t.Errorf("DeleteBookmark() error = %v, want %v", err.Error(), expectedError)
	}

	// Test deleting already deleted bookmark
	err = repo.DeleteBookmark(created.ID)
	if err == nil {
		t.Error("DeleteBookmark() on already deleted bookmark should return error")
		return
	}

	expectedError = fmt.Sprintf("bookmark with ID %s not found", created.ID)
	if err.Error() != expectedError {
		t.Errorf("DeleteBookmark() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestBookmarkInMemRepository_EmptyRepository(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Test listing bookmarks in empty repository
	bookmarks, err := repo.ListBookmarks("user1")
	if err != nil {
		t.Errorf("ListBookmarks() on empty repository failed: %v", err)
	}
	if len(bookmarks) != 0 {
		t.Errorf("ListBookmarks() on empty repository should return empty slice, got %d bookmarks", len(bookmarks))
	}

	// Test updating non-existent bookmark
	nonExistentBookmark := model.Bookmark{
		ID:     "nonexistent",
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Non-existent",
	}
	_, err = repo.UpdateBookmark(nonExistentBookmark)
	if err == nil {
		t.Error("UpdateBookmark() on empty repository should return error")
	}
	expectedError := "bookmark with ID nonexistent not found"
	if err.Error() != expectedError {
		t.Errorf("UpdateBookmark() error = %v, want %v", err.Error(), expectedError)
	}

	// Test deleting non-existent bookmark
	err = repo.DeleteBookmark("nonexistent")
	if err == nil {
		t.Error("DeleteBookmark() on empty repository should return error")
	}
	if err.Error() != expectedError {
		t.Errorf("DeleteBookmark() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestBookmarkInMemRepository_CreatedAtHandling(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Test bookmark without CreatedAt
	bookmark := model.Bookmark{
		ID:     "1",
		UserID: "user1",
		URL:    "https://example.com",
		Title:  "Example",
	}

	result, err := repo.CreateBookmark(bookmark)
	if err != nil {
		t.Errorf("CreateBookmark() unexpected error = %v", err)
		return
	}

	// CreatedAt should be set automatically
	if result.CreatedAt.IsZero() {
		t.Error("CreateBookmark() should set CreatedAt when not provided")
	}

	// Test bookmark with existing CreatedAt
	specificTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	bookmark2 := model.Bookmark{
		ID:        "2",
		UserID:    "user1",
		URL:       "https://example2.com",
		Title:     "Example 2",
		CreatedAt: specificTime,
	}

	result2, err := repo.CreateBookmark(bookmark2)
	if err != nil {
		t.Errorf("CreateBookmark() unexpected error = %v", err)
		return
	}

	// CreatedAt should be preserved
	if !result2.CreatedAt.Equal(specificTime) {
		t.Errorf("CreateBookmark() should preserve existing CreatedAt = %v, got %v", specificTime, result2.CreatedAt)
	}
}

func TestBookmarkInMemRepository_ConcurrentAccess(t *testing.T) {
	repo := NewBookmarkInMemRepository()
	var wg sync.WaitGroup
	numGoroutines := 10
	bookmarksPerGoroutine := 10

	// Store created bookmark IDs for each user
	type userBookmarks struct {
		sync.Mutex
		ids []string
	}
	bookmarkIDsPerUser := make(map[int]*userBookmarks)
	for i := 0; i < numGoroutines; i++ {
		bookmarkIDsPerUser[i] = &userBookmarks{ids: make([]string, 0, bookmarksPerGoroutine)}
	}

	// Test concurrent creates
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < bookmarksPerGoroutine; j++ {
				bookmark := model.Bookmark{
					UserID: fmt.Sprintf("user_%d", goroutineID),
					URL:    fmt.Sprintf("https://example%d-%d.com", goroutineID, j),
					Title:  fmt.Sprintf("Title %d-%d", goroutineID, j),
				}
				created, err := repo.CreateBookmark(bookmark)
				if err != nil {
					t.Errorf("Concurrent CreateBookmark() error: %v", err)
				} else {
					// Store the auto-generated ID
					bookmarkIDsPerUser[goroutineID].Lock()
					bookmarkIDsPerUser[goroutineID].ids = append(bookmarkIDsPerUser[goroutineID].ids, created.ID)
					bookmarkIDsPerUser[goroutineID].Unlock()
				}
			}
		}(i)
	}
	wg.Wait()

	// Verify all bookmarks were created
	totalBookmarks := 0
	for i := 0; i < numGoroutines; i++ {
		userID := fmt.Sprintf("user_%d", i)
		bookmarks, err := repo.ListBookmarks(userID)
		if err != nil {
			t.Errorf("ListBookmarks() after concurrent creates failed: %v", err)
		}
		totalBookmarks += len(bookmarks)
		if len(bookmarks) != bookmarksPerGoroutine {
			t.Errorf("Expected %d bookmarks for user %s, got %d", bookmarksPerGoroutine, userID, len(bookmarks))
		}
	}

	expectedTotal := numGoroutines * bookmarksPerGoroutine
	if totalBookmarks != expectedTotal {
		t.Errorf("Expected total %d bookmarks, got %d", expectedTotal, totalBookmarks)
	}

	// Test concurrent reads and writes using actual bookmark IDs
	wg.Add(numGoroutines * 2)
	for i := 0; i < numGoroutines; i++ {
		// Reader goroutine
		go func(goroutineID int) {
			defer wg.Done()
			userID := fmt.Sprintf("user_%d", goroutineID)
			for j := 0; j < 5; j++ {
				_, err := repo.ListBookmarks(userID)
				if err != nil {
					t.Errorf("Concurrent ListBookmarks() error: %v", err)
				}
			}
		}(i)

		// Writer goroutine (update) - use actual bookmark IDs
		go func(goroutineID int) {
			defer wg.Done()
			bookmarkIDsPerUser[goroutineID].Lock()
			ids := bookmarkIDsPerUser[goroutineID].ids
			bookmarkIDsPerUser[goroutineID].Unlock()

			for j := 0; j < 5 && j < len(ids); j++ {
				updatedBookmark := model.Bookmark{
					ID:     ids[j],
					UserID: fmt.Sprintf("user_%d", goroutineID),
					URL:    fmt.Sprintf("https://updated%d-%d.com", goroutineID, j),
					Title:  fmt.Sprintf("Updated Title %d-%d", goroutineID, j),
				}
				_, err := repo.UpdateBookmark(updatedBookmark)
				if err != nil {
					t.Errorf("Concurrent UpdateBookmark() error: %v", err)
				}
			}
		}(i)
	}
	wg.Wait()
}

func TestBookmarkInMemRepository_EdgeCases(t *testing.T) {
	repo := NewBookmarkInMemRepository()

	// Test with empty string values (ID will be auto-generated)
	emptyBookmark := model.Bookmark{
		UserID: "",
		URL:    "",
		Title:  "",
	}

	result, err := repo.CreateBookmark(emptyBookmark)
	if err != nil {
		t.Errorf("CreateBookmark() with empty values failed: %v", err)
	}
	// ID should be auto-generated even if other fields are empty
	if result.ID == "" {
		t.Error("CreateBookmark() should auto-generate ID even with empty other fields")
	}
	if result.UserID != "" || result.URL != "" || result.Title != "" {
		t.Error("CreateBookmark() should preserve empty string values for non-ID fields")
	}

	// Test listing bookmarks for empty userID
	bookmarks, err := repo.ListBookmarks("")
	if err != nil {
		t.Errorf("ListBookmarks() with empty userID failed: %v", err)
	}
	if len(bookmarks) != 1 {
		t.Errorf("ListBookmarks() with empty userID should return 1 bookmark, got %d", len(bookmarks))
	}

	// Test with special characters in UserID
	specialBookmark := model.Bookmark{
		UserID: "user@#$%^&*()",
		URL:    "https://example.com",
		Title:  "Special Characters",
	}

	_, err = repo.CreateBookmark(specialBookmark)
	if err != nil {
		t.Errorf("CreateBookmark() with special characters failed: %v", err)
	}

	// Verify it can be retrieved
	bookmarks, err = repo.ListBookmarks("user@#$%^&*()")
	if err != nil {
		t.Errorf("ListBookmarks() with special character userID failed: %v", err)
	}
	if len(bookmarks) != 1 {
		t.Errorf("ListBookmarks() should return 1 bookmark with special characters, got %d", len(bookmarks))
	}
}
