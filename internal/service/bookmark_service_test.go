package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/tsongpon/athena/internal/model"
)

// MockBookmarkRepository is a mock implementation of BookmarkRepository for testing
type MockBookmarkRepository struct {
	createBookmarkFunc func(bookmark model.Bookmark) (model.Bookmark, error)
	getBookmarkFunc    func(id string) (model.Bookmark, error)
	listBookmarksFunc  func(userID string) ([]model.Bookmark, error)
	updateBookmarkFunc func(bookmark model.Bookmark) (model.Bookmark, error)
	deleteBookmarkFunc func(id string) error
}

func (m *MockBookmarkRepository) CreateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
	if m.createBookmarkFunc != nil {
		return m.createBookmarkFunc(bookmark)
	}
	return model.Bookmark{}, nil
}

func (m *MockBookmarkRepository) GetBookmark(id string) (model.Bookmark, error) {
	if m.getBookmarkFunc != nil {
		return m.getBookmarkFunc(id)
	}
	return model.Bookmark{}, nil
}

func (m *MockBookmarkRepository) ListBookmarks(userID string) ([]model.Bookmark, error) {
	if m.listBookmarksFunc != nil {
		return m.listBookmarksFunc(userID)
	}
	return []model.Bookmark{}, nil
}

func (m *MockBookmarkRepository) UpdateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
	if m.updateBookmarkFunc != nil {
		return m.updateBookmarkFunc(bookmark)
	}
	return model.Bookmark{}, nil
}

func (m *MockBookmarkRepository) DeleteBookmark(id string) error {
	if m.deleteBookmarkFunc != nil {
		return m.deleteBookmarkFunc(id)
	}
	return nil
}

// TestBookmarkService_CreateBookmark tests successful bookmark creation
func TestBookmarkService_CreateBookmark(t *testing.T) {
	expectedBookmark := model.Bookmark{
		ID:        "bookmark-1",
		UserID:    "user-1",
		URL:       "https://example.com",
		Title:     "Example",
		CreatedAt: time.Now(),
	}

	mockRepo := &MockBookmarkRepository{
		createBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			if bookmark.UserID != "user-1" {
				t.Errorf("CreateBookmark() received UserID = %v, want user-1", bookmark.UserID)
			}
			if bookmark.URL != "https://example.com" {
				t.Errorf("CreateBookmark() received URL = %v, want https://example.com", bookmark.URL)
			}
			return expectedBookmark, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.CreateBookmark(model.Bookmark{
		UserID: "user-1",
		URL:    "https://example.com",
	})

	if err != nil {
		t.Errorf("CreateBookmark() unexpected error = %v", err)
		return
	}

	if result.ID != expectedBookmark.ID {
		t.Errorf("CreateBookmark() result ID = %v, want %v", result.ID, expectedBookmark.ID)
	}
	if result.UserID != expectedBookmark.UserID {
		t.Errorf("CreateBookmark() result UserID = %v, want %v", result.UserID, expectedBookmark.UserID)
	}
	if result.URL != expectedBookmark.URL {
		t.Errorf("CreateBookmark() result URL = %v, want %v", result.URL, expectedBookmark.URL)
	}
}

// TestBookmarkService_CreateBookmark_RepositoryError tests error handling when repository fails
func TestBookmarkService_CreateBookmark_RepositoryError(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		createBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			return model.Bookmark{}, fmt.Errorf("database connection failed")
		},
	}

	service := NewBookmarkService(mockRepo)
	_, err := service.CreateBookmark(model.Bookmark{
		UserID: "user-1",
		URL:    "https://example.com",
	})

	if err == nil {
		t.Error("CreateBookmark() should return error when repository fails")
		return
	}

	expectedErrorSubstring := "failed to create bookmark for URL https://example.com"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("CreateBookmark() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_CreateBookmark_EmptyURL tests creating bookmark with empty URL
func TestBookmarkService_CreateBookmark_EmptyURL(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		createBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			if bookmark.URL != "" {
				t.Errorf("CreateBookmark() should pass empty URL to repository")
			}
			return bookmark, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.CreateBookmark(model.Bookmark{
		UserID: "user-1",
		URL:    "",
	})

	if err != nil {
		t.Errorf("CreateBookmark() with empty URL unexpected error = %v", err)
		return
	}

	if result.URL != "" {
		t.Errorf("CreateBookmark() result URL = %v, want empty string", result.URL)
	}
}

// TestBookmarkService_CreateBookmark_EmptyUserID tests creating bookmark with empty userID
func TestBookmarkService_CreateBookmark_EmptyUserID(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		createBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			if bookmark.UserID != "" {
				t.Errorf("CreateBookmark() should pass empty UserID to repository")
			}
			return bookmark, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.CreateBookmark(model.Bookmark{
		UserID: "",
		URL:    "https://example.com",
	})

	if err != nil {
		t.Errorf("CreateBookmark() with empty UserID unexpected error = %v", err)
		return
	}

	if result.UserID != "" {
		t.Errorf("CreateBookmark() result UserID = %v, want empty string", result.UserID)
	}
}

// TestBookmarkService_CreateBookmark_NonEmptyID tests creating bookmark with non-empty ID should fail
func TestBookmarkService_CreateBookmark_NonEmptyID(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		createBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			t.Error("CreateBookmark() should not be called when ID is not empty")
			return model.Bookmark{}, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	_, err := service.CreateBookmark(model.Bookmark{
		ID:     "existing-id",
		UserID: "user-1",
		URL:    "https://example.com",
	})

	if err == nil {
		t.Error("CreateBookmark() should return error when ID is not empty")
		return
	}

	expectedError := "bookmark ID must be empty"
	if err.Error() != expectedError {
		t.Errorf("CreateBookmark() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestBookmarkService_ArchiveBookmark tests successful bookmark archiving
func TestBookmarkService_ArchiveBookmark(t *testing.T) {
	existingBookmark := model.Bookmark{
		ID:         "bookmark-1",
		UserID:     "user-1",
		URL:        "https://example.com",
		Title:      "Example",
		IsArchived: false,
		CreatedAt:  time.Now(),
	}

	archivedBookmark := existingBookmark
	archivedBookmark.IsArchived = true

	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			if id != "bookmark-1" {
				t.Errorf("GetBookmark() received ID = %v, want bookmark-1", id)
			}
			return existingBookmark, nil
		},
		updateBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			if !bookmark.IsArchived {
				t.Error("UpdateBookmark() should receive bookmark with IsArchived = true")
			}
			if bookmark.ID != "bookmark-1" {
				t.Errorf("UpdateBookmark() received ID = %v, want bookmark-1", bookmark.ID)
			}
			return archivedBookmark, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.ArchiveBookmark("bookmark-1")

	if err != nil {
		t.Errorf("ArchiveBookmark() unexpected error = %v", err)
		return
	}

	if !result.IsArchived {
		t.Error("ArchiveBookmark() result should have IsArchived = true")
	}
	if result.ID != "bookmark-1" {
		t.Errorf("ArchiveBookmark() result ID = %v, want bookmark-1", result.ID)
	}
}

// TestBookmarkService_ArchiveBookmark_BookmarkNotFound tests error when bookmark doesn't exist
func TestBookmarkService_ArchiveBookmark_BookmarkNotFound(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			return model.Bookmark{}, fmt.Errorf("bookmark with ID %s not found", id)
		},
	}

	service := NewBookmarkService(mockRepo)
	_, err := service.ArchiveBookmark("nonexistent")

	if err == nil {
		t.Error("ArchiveBookmark() should return error when bookmark not found")
		return
	}

	expectedErrorSubstring := "failed to get bookmark with ID nonexistent"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("ArchiveBookmark() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_ArchiveBookmark_UnauthorizedUser tests that the service archives bookmark regardless of user
// Note: The current implementation doesn't check authorization - it's the caller's responsibility
func TestBookmarkService_ArchiveBookmark_UnauthorizedUser(t *testing.T) {
	unauthorizedBookmark := model.Bookmark{
		ID:         "bookmark-1",
		UserID:     "user-2",
		URL:        "https://example.com",
		Title:      "Example",
		IsArchived: false,
		CreatedAt:  time.Now(),
	}

	archivedBookmark := unauthorizedBookmark
	archivedBookmark.IsArchived = true

	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			if id != "bookmark-1" {
				t.Errorf("GetBookmark() received ID = %v, want bookmark-1", id)
			}
			return unauthorizedBookmark, nil
		},
		updateBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			if !bookmark.IsArchived {
				t.Error("UpdateBookmark() should receive bookmark with IsArchived = true")
			}
			return archivedBookmark, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.ArchiveBookmark("bookmark-1")

	if err != nil {
		t.Errorf("ArchiveBookmark() unexpected error = %v", err)
		return
	}

	if !result.IsArchived {
		t.Error("ArchiveBookmark() result should have IsArchived = true")
	}
}

// TestBookmarkService_ArchiveBookmark_UpdateError tests error handling during update
func TestBookmarkService_ArchiveBookmark_UpdateError(t *testing.T) {
	existingBookmark := model.Bookmark{
		ID:         "bookmark-1",
		UserID:     "user-1",
		URL:        "https://example.com",
		Title:      "Example",
		IsArchived: false,
		CreatedAt:  time.Now(),
	}

	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			return existingBookmark, nil
		},
		updateBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			return model.Bookmark{}, fmt.Errorf("database update failed")
		},
	}

	service := NewBookmarkService(mockRepo)
	_, err := service.ArchiveBookmark("bookmark-1")

	if err == nil {
		t.Error("ArchiveBookmark() should return error when update fails")
		return
	}

	expectedErrorSubstring := "failed to update bookmark with ID bookmark-1"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("ArchiveBookmark() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_ArchiveBookmark_AlreadyArchived tests archiving an already archived bookmark
func TestBookmarkService_ArchiveBookmark_AlreadyArchived(t *testing.T) {
	alreadyArchivedBookmark := model.Bookmark{
		ID:         "bookmark-1",
		UserID:     "user-1",
		URL:        "https://example.com",
		Title:      "Example",
		IsArchived: true,
		CreatedAt:  time.Now(),
	}

	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			return alreadyArchivedBookmark, nil
		},
		updateBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			if !bookmark.IsArchived {
				t.Error("UpdateBookmark() should receive bookmark with IsArchived = true")
			}
			return bookmark, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.ArchiveBookmark("bookmark-1")

	if err != nil {
		t.Errorf("ArchiveBookmark() on already archived bookmark unexpected error = %v", err)
		return
	}

	if !result.IsArchived {
		t.Error("ArchiveBookmark() result should have IsArchived = true")
	}
}

// TestBookmarkService_ArchiveBookmark_EmptyUserID tests archiving with empty userID
func TestBookmarkService_ArchiveBookmark_EmptyUserID(t *testing.T) {
	existingBookmark := model.Bookmark{
		ID:         "bookmark-1",
		UserID:     "",
		URL:        "https://example.com",
		Title:      "Example",
		IsArchived: false,
		CreatedAt:  time.Now(),
	}

	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			return existingBookmark, nil
		},
		updateBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			bookmark.IsArchived = true
			return bookmark, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.ArchiveBookmark("bookmark-1")

	if err != nil {
		t.Errorf("ArchiveBookmark() with empty userID unexpected error = %v", err)
		return
	}

	if !result.IsArchived {
		t.Error("ArchiveBookmark() result should have IsArchived = true")
	}
}

// TestBookmarkService_ArchiveBookmark_EmptyBookmarkID tests archiving with empty bookmarkID
func TestBookmarkService_ArchiveBookmark_EmptyBookmarkID(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			if id != "" {
				t.Errorf("GetBookmark() should receive empty ID")
			}
			return model.Bookmark{}, fmt.Errorf("bookmark with ID  not found")
		},
	}

	service := NewBookmarkService(mockRepo)
	_, err := service.ArchiveBookmark("")

	if err == nil {
		t.Error("ArchiveBookmark() with empty bookmarkID should return error")
		return
	}

	expectedErrorSubstring := "failed to get bookmark with ID "
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("ArchiveBookmark() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_GetBookmark tests successful bookmark retrieval
func TestBookmarkService_GetBookmark(t *testing.T) {
	expectedBookmark := model.Bookmark{
		ID:        "bookmark-1",
		UserID:    "user-1",
		URL:       "https://example.com",
		Title:     "Example",
		CreatedAt: time.Now(),
	}

	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			if id != "bookmark-1" {
				t.Errorf("GetBookmark() received ID = %v, want bookmark-1", id)
			}
			return expectedBookmark, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.GetBookmark("bookmark-1")

	if err != nil {
		t.Errorf("GetBookmark() unexpected error = %v", err)
		return
	}

	if result.ID != expectedBookmark.ID {
		t.Errorf("GetBookmark() result ID = %v, want %v", result.ID, expectedBookmark.ID)
	}
	if result.UserID != expectedBookmark.UserID {
		t.Errorf("GetBookmark() result UserID = %v, want %v", result.UserID, expectedBookmark.UserID)
	}
	if result.URL != expectedBookmark.URL {
		t.Errorf("GetBookmark() result URL = %v, want %v", result.URL, expectedBookmark.URL)
	}
	if result.Title != expectedBookmark.Title {
		t.Errorf("GetBookmark() result Title = %v, want %v", result.Title, expectedBookmark.Title)
	}
}

// TestBookmarkService_GetBookmark_EmptyID tests getting bookmark with empty ID
func TestBookmarkService_GetBookmark_EmptyID(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			t.Error("GetBookmark() should not be called with empty ID")
			return model.Bookmark{}, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	_, err := service.GetBookmark("")

	if err == nil {
		t.Error("GetBookmark() should return error when ID is empty")
		return
	}

	expectedError := "id is required"
	if err.Error() != expectedError {
		t.Errorf("GetBookmark() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestBookmarkService_GetBookmark_RepositoryError tests error handling when repository fails
func TestBookmarkService_GetBookmark_RepositoryError(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			return model.Bookmark{}, fmt.Errorf("database connection failed")
		},
	}

	service := NewBookmarkService(mockRepo)
	_, err := service.GetBookmark("bookmark-1")

	if err == nil {
		t.Error("GetBookmark() should return error when repository fails")
		return
	}

	expectedErrorSubstring := "failed to get bookmarks for id bookmark-1"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("GetBookmark() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_GetBookmark_NotFound tests getting non-existent bookmark
func TestBookmarkService_GetBookmark_NotFound(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			return model.Bookmark{}, fmt.Errorf("bookmark with ID %s not found", id)
		},
	}

	service := NewBookmarkService(mockRepo)
	_, err := service.GetBookmark("nonexistent-id")

	if err == nil {
		t.Error("GetBookmark() should return error when bookmark not found")
		return
	}

	expectedErrorSubstring := "failed to get bookmarks for id nonexistent-id"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("GetBookmark() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_GetAllBookmarks tests successful retrieval of all bookmarks for a user
func TestBookmarkService_GetAllBookmarks(t *testing.T) {
	expectedBookmarks := []model.Bookmark{
		{
			ID:        "bookmark-1",
			UserID:    "user-1",
			URL:       "https://example1.com",
			Title:     "Example 1",
			CreatedAt: time.Now(),
		},
		{
			ID:        "bookmark-2",
			UserID:    "user-1",
			URL:       "https://example2.com",
			Title:     "Example 2",
			CreatedAt: time.Now(),
		},
	}

	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string) ([]model.Bookmark, error) {
			if userID != "user-1" {
				t.Errorf("ListBookmarks() received UserID = %v, want user-1", userID)
			}
			return expectedBookmarks, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.GetAllBookmarks("user-1")

	if err != nil {
		t.Errorf("GetAllBookmarks() unexpected error = %v", err)
		return
	}

	if len(result) != 2 {
		t.Errorf("GetAllBookmarks() returned %d bookmarks, want 2", len(result))
		return
	}

	if result[0].ID != "bookmark-1" {
		t.Errorf("GetAllBookmarks() first bookmark ID = %v, want bookmark-1", result[0].ID)
	}
	if result[1].ID != "bookmark-2" {
		t.Errorf("GetAllBookmarks() second bookmark ID = %v, want bookmark-2", result[1].ID)
	}
}

// TestBookmarkService_GetAllBookmarks_EmptyUserID tests getting bookmarks with empty userID
func TestBookmarkService_GetAllBookmarks_EmptyUserID(t *testing.T) {
	expectedBookmarks := []model.Bookmark{}

	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string) ([]model.Bookmark, error) {
			if userID != "" {
				t.Errorf("ListBookmarks() received UserID = %v, want empty string", userID)
			}
			return expectedBookmarks, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.GetAllBookmarks("")

	if err != nil {
		t.Errorf("GetAllBookmarks() with empty userID unexpected error = %v", err)
		return
	}

	if len(result) != 0 {
		t.Errorf("GetAllBookmarks() returned %d bookmarks, want 0", len(result))
	}
}

// TestBookmarkService_GetAllBookmarks_RepositoryError tests error handling when repository fails
func TestBookmarkService_GetAllBookmarks_RepositoryError(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string) ([]model.Bookmark, error) {
			return nil, fmt.Errorf("database connection failed")
		},
	}

	service := NewBookmarkService(mockRepo)
	_, err := service.GetAllBookmarks("user-1")

	if err == nil {
		t.Error("GetAllBookmarks() should return error when repository fails")
		return
	}

	expectedErrorSubstring := "failed to get all bookmarks"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("GetAllBookmarks() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_GetAllBookmarks_NoBookmarks tests getting bookmarks when user has none
func TestBookmarkService_GetAllBookmarks_NoBookmarks(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string) ([]model.Bookmark, error) {
			return []model.Bookmark{}, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	result, err := service.GetAllBookmarks("user-1")

	if err != nil {
		t.Errorf("GetAllBookmarks() unexpected error = %v", err)
		return
	}

	if len(result) != 0 {
		t.Errorf("GetAllBookmarks() returned %d bookmarks, want 0", len(result))
	}
}

// TestBookmarkService_DeleteBookmark tests successful bookmark deletion
func TestBookmarkService_DeleteBookmark(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		deleteBookmarkFunc: func(id string) error {
			if id != "bookmark-1" {
				t.Errorf("DeleteBookmark() received ID = %v, want bookmark-1", id)
			}
			return nil
		},
	}

	service := NewBookmarkService(mockRepo)
	err := service.DeleteBookmark("bookmark-1")

	if err != nil {
		t.Errorf("DeleteBookmark() unexpected error = %v", err)
	}
}

// TestBookmarkService_DeleteBookmark_EmptyID tests deleting bookmark with empty ID
func TestBookmarkService_DeleteBookmark_EmptyID(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		deleteBookmarkFunc: func(id string) error {
			t.Error("DeleteBookmark() should not be called with empty ID")
			return nil
		},
	}

	service := NewBookmarkService(mockRepo)
	err := service.DeleteBookmark("")

	if err == nil {
		t.Error("DeleteBookmark() should return error when ID is empty")
		return
	}

	expectedError := "id is required"
	if err.Error() != expectedError {
		t.Errorf("DeleteBookmark() error = %v, want %v", err.Error(), expectedError)
	}
}

// TestBookmarkService_DeleteBookmark_RepositoryError tests error handling when repository fails
func TestBookmarkService_DeleteBookmark_RepositoryError(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		deleteBookmarkFunc: func(id string) error {
			return fmt.Errorf("database connection failed")
		},
	}

	service := NewBookmarkService(mockRepo)
	err := service.DeleteBookmark("bookmark-1")

	if err == nil {
		t.Error("DeleteBookmark() should return error when repository fails")
		return
	}

	expectedErrorSubstring := "failed to delete bookmark with ID bookmark-1"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("DeleteBookmark() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_DeleteBookmark_NotFound tests deleting non-existent bookmark
func TestBookmarkService_DeleteBookmark_NotFound(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		deleteBookmarkFunc: func(id string) error {
			return fmt.Errorf("bookmark with ID %s not found", id)
		},
	}

	service := NewBookmarkService(mockRepo)
	err := service.DeleteBookmark("nonexistent-id")

	if err == nil {
		t.Error("DeleteBookmark() should return error when bookmark not found")
		return
	}

	expectedErrorSubstring := "failed to delete bookmark with ID nonexistent-id"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("DeleteBookmark() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestNewBookmarkService tests service initialization
func TestNewBookmarkService(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			return model.Bookmark{ID: "test-id"}, nil
		},
	}
	service := NewBookmarkService(mockRepo)

	// Verify service is usable by calling a method
	result, err := service.GetBookmark("test-id")
	if err != nil {
		t.Errorf("NewBookmarkService() service should be functional, got error = %v", err)
	}
	if result.ID != "test-id" {
		t.Error("NewBookmarkService() service should work correctly")
	}
}
