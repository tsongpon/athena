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
	result, err := service.CreateBookmark("user-1", "https://example.com")

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
	_, err := service.CreateBookmark("user-1", "https://example.com")

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
	result, err := service.CreateBookmark("user-1", "")

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
	result, err := service.CreateBookmark("", "https://example.com")

	if err != nil {
		t.Errorf("CreateBookmark() with empty UserID unexpected error = %v", err)
		return
	}

	if result.UserID != "" {
		t.Errorf("CreateBookmark() result UserID = %v, want empty string", result.UserID)
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
	result, err := service.ArchiveBookmark("user-1", "bookmark-1")

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
	_, err := service.ArchiveBookmark("user-1", "nonexistent")

	if err == nil {
		t.Error("ArchiveBookmark() should return error when bookmark not found")
		return
	}

	expectedErrorSubstring := "failed to get bookmark with ID nonexistent"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("ArchiveBookmark() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_ArchiveBookmark_UnauthorizedUser tests authorization check
func TestBookmarkService_ArchiveBookmark_UnauthorizedUser(t *testing.T) {
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
			t.Error("UpdateBookmark() should not be called for unauthorized user")
			return model.Bookmark{}, nil
		},
	}

	service := NewBookmarkService(mockRepo)
	_, err := service.ArchiveBookmark("user-2", "bookmark-1")

	if err == nil {
		t.Error("ArchiveBookmark() should return error for unauthorized user")
		return
	}

	expectedError := "user user-2 is not authorized to archive bookmark bookmark-1"
	if err.Error() != expectedError {
		t.Errorf("ArchiveBookmark() error = %v, want %v", err.Error(), expectedError)
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
	_, err := service.ArchiveBookmark("user-1", "bookmark-1")

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
	result, err := service.ArchiveBookmark("user-1", "bookmark-1")

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
	result, err := service.ArchiveBookmark("", "bookmark-1")

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
	_, err := service.ArchiveBookmark("user-1", "")

	if err == nil {
		t.Error("ArchiveBookmark() with empty bookmarkID should return error")
		return
	}
}

// TestNewBookmarkService tests service initialization
func TestNewBookmarkService(t *testing.T) {
	mockRepo := &MockBookmarkRepository{}
	service := NewBookmarkService(mockRepo)

	if service.repository == nil {
		t.Error("NewBookmarkService() should initialize repository")
	}
}
