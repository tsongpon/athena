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
	listBookmarksFunc  func(userID string, archived bool) ([]model.Bookmark, error)
	countBookmarksFunc func(query model.BookmarkQuery) (int, error)
	updateBookmarkFunc func(bookmark model.Bookmark) (model.Bookmark, error)
	deleteBookmarkFunc func(id string) error
}

// MockWebRepository is a mock implementation of WebRepository for testing
type MockWebRepository struct {
	getTitleFunc          func(url string) (string, error)
	getMainImageFunc      func(url string) (string, error)
	getContentSummaryFunc func(url string) (string, error)
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

func (m *MockBookmarkRepository) ListBookmarks(query model.BookmarkQuery) ([]model.Bookmark, error) {
	if m.listBookmarksFunc != nil {
		return m.listBookmarksFunc(query.UserID, query.Archived)
	}
	return []model.Bookmark{}, nil
}

func (m *MockBookmarkRepository) CountBookmarks(query model.BookmarkQuery) (int, error) {
	if m.countBookmarksFunc != nil {
		return m.countBookmarksFunc(query)
	}
	return 0, nil
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

func (m *MockWebRepository) GetTitle(url string) (string, error) {
	if m.getTitleFunc != nil {
		return m.getTitleFunc(url)
	}
	return "Default Title", nil
}

func (m *MockWebRepository) GetMainImage(url string) (string, error) {
	if m.getMainImageFunc != nil {
		return m.getMainImageFunc(url)
	}
	return "", nil
}

func (m *MockWebRepository) GetContentSummary(url string) (string, error) {
	if m.getContentSummaryFunc != nil {
		return m.getContentSummaryFunc(url)
	}
	return "", nil
}

// TestBookmarkService_CreateBookmark tests successful bookmark creation
func TestBookmarkService_CreateBookmark(t *testing.T) {
	expectedBookmark := model.Bookmark{
		ID:             "bookmark-1",
		UserID:         "user-1",
		URL:            "https://example.com",
		Title:          "Example",
		MainImageURL:   "https://example.com/og-image.jpg",
		ContentSummary: "This is an example website with useful content.",
		CreatedAt:      time.Now(),
	}

	mockRepo := &MockBookmarkRepository{
		createBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			if bookmark.UserID != "user-1" {
				t.Errorf("CreateBookmark() received UserID = %v, want user-1", bookmark.UserID)
			}
			if bookmark.URL != "https://example.com" {
				t.Errorf("CreateBookmark() received URL = %v, want https://example.com", bookmark.URL)
			}
			if bookmark.Title != "Example" {
				t.Errorf("CreateBookmark() received Title = %v, want Example", bookmark.Title)
			}
			if bookmark.MainImageURL != "https://example.com/og-image.jpg" {
				t.Errorf("CreateBookmark() received MainImageURL = %v, want https://example.com/og-image.jpg", bookmark.MainImageURL)
			}
			if bookmark.ContentSummary != "This is an example website with useful content." {
				t.Errorf("CreateBookmark() received ContentSummary = %v, want This is an example website with useful content.", bookmark.ContentSummary)
			}
			return expectedBookmark, nil
		},
	}

	mockWebRepo := &MockWebRepository{
		getTitleFunc: func(url string) (string, error) {
			return "Example", nil
		},
		getMainImageFunc: func(url string) (string, error) {
			return "https://example.com/og-image.jpg", nil
		},
		getContentSummaryFunc: func(url string) (string, error) {
			return "This is an example website with useful content.", nil
		},
	}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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
	if result.MainImageURL != expectedBookmark.MainImageURL {
		t.Errorf("CreateBookmark() result MainImageURL = %v, want %v", result.MainImageURL, expectedBookmark.MainImageURL)
	}
	if result.ContentSummary != expectedBookmark.ContentSummary {
		t.Errorf("CreateBookmark() result ContentSummary = %v, want %v", result.ContentSummary, expectedBookmark.ContentSummary)
	}
}

// TestBookmarkService_CreateBookmark_RepositoryError tests error handling when repository fails
func TestBookmarkService_CreateBookmark_RepositoryError(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		createBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			return model.Bookmark{}, fmt.Errorf("database connection failed")
		},
	}

	mockWebRepo := &MockWebRepository{
		getTitleFunc: func(url string) (string, error) {
			return "Test Title", nil
		},
		getMainImageFunc: func(url string) (string, error) {
			return "https://example.com/image.jpg", nil
		},
		getContentSummaryFunc: func(url string) (string, error) {
			return "Summary text", nil
		},
	}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

// TestBookmarkService_CreateBookmark_GetMainImageError tests error handling when GetMainImage fails
func TestBookmarkService_CreateBookmark_GetMainImageError(t *testing.T) {
	mockRepo := &MockBookmarkRepository{}

	mockWebRepo := &MockWebRepository{
		getTitleFunc: func(url string) (string, error) {
			return "Test Title", nil
		},
		getMainImageFunc: func(url string) (string, error) {
			return "", fmt.Errorf("failed to fetch og:image")
		},
		getContentSummaryFunc: func(url string) (string, error) {
			t.Error("GetContentSummary() should not be called when GetMainImage fails first")
			return "", nil
		},
	}
	service := NewBookmarkService(mockRepo, mockWebRepo)
	_, err := service.CreateBookmark(model.Bookmark{
		UserID: "user-1",
		URL:    "https://example.com",
	})

	if err == nil {
		t.Error("CreateBookmark() should return error when GetMainImage fails")
		return
	}

	expectedErrorSubstring := "failed to fetch main image URL for URL https://example.com"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("CreateBookmark() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_CreateBookmark_GetContentSummaryError tests error handling when GetContentSummary fails
func TestBookmarkService_CreateBookmark_GetContentSummaryError(t *testing.T) {
	mockRepo := &MockBookmarkRepository{}

	mockWebRepo := &MockWebRepository{
		getTitleFunc: func(url string) (string, error) {
			return "Test Title", nil
		},
		getMainImageFunc: func(url string) (string, error) {
			return "https://example.com/image.jpg", nil
		},
		getContentSummaryFunc: func(url string) (string, error) {
			return "", fmt.Errorf("failed to generate content summary")
		},
	}
	service := NewBookmarkService(mockRepo, mockWebRepo)
	_, err := service.CreateBookmark(model.Bookmark{
		UserID: "user-1",
		URL:    "https://example.com",
	})

	if err == nil {
		t.Error("CreateBookmark() should return error when GetContentSummary fails")
		return
	}

	expectedErrorSubstring := "failed to fetch content for URL https://example.com"
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

	mockWebRepo := &MockWebRepository{
		getTitleFunc: func(url string) (string, error) {
			return "", nil
		},
		getMainImageFunc: func(url string) (string, error) {
			return "", nil
		},
		getContentSummaryFunc: func(url string) (string, error) {
			return "", nil
		},
	}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{
		getTitleFunc: func(url string) (string, error) {
			return "Title", nil
		},
		getMainImageFunc: func(url string) (string, error) {
			return "https://example.com/image.jpg", nil
		},
		getContentSummaryFunc: func(url string) (string, error) {
			return "Content summary", nil
		},
	}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{
		getTitleFunc: func(url string) (string, error) {
			t.Error("GetTitle() should not be called when ID validation fails")
			return "", nil
		},
		getMainImageFunc: func(url string) (string, error) {
			t.Error("GetMainImage() should not be called when ID validation fails")
			return "", nil
		},
	}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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
		ID:             "bookmark-1",
		UserID:         "user-1",
		URL:            "https://example.com",
		Title:          "Example",
		MainImageURL:   "https://example.com/og-image.jpg",
		ContentSummary: "Example website content summary.",
		CreatedAt:      time.Now(),
	}

	mockRepo := &MockBookmarkRepository{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			if id != "bookmark-1" {
				t.Errorf("GetBookmark() received ID = %v, want bookmark-1", id)
			}
			return expectedBookmark, nil
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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
	if result.MainImageURL != expectedBookmark.MainImageURL {
		t.Errorf("GetBookmark() result MainImageURL = %v, want %v", result.MainImageURL, expectedBookmark.MainImageURL)
	}
	if result.ContentSummary != expectedBookmark.ContentSummary {
		t.Errorf("GetBookmark() result ContentSummary = %v, want %v", result.ContentSummary, expectedBookmark.ContentSummary)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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
			ID:             "bookmark-1",
			UserID:         "user-1",
			URL:            "https://example1.com",
			Title:          "Example 1",
			MainImageURL:   "https://example1.com/image1.jpg",
			ContentSummary: "First example website content.",
			CreatedAt:      time.Now(),
		},
		{
			ID:             "bookmark-2",
			UserID:         "user-1",
			URL:            "https://example2.com",
			Title:          "Example 2",
			MainImageURL:   "https://example2.com/image2.jpg",
			ContentSummary: "Second example website content.",
			CreatedAt:      time.Now(),
		},
	}

	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			if userID != "user-1" {
				t.Errorf("ListBookmarks() received UserID = %v, want user-1", userID)
			}
			return expectedBookmarks, nil
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
	result, err := service.GetAllBookmarks("user-1", false)

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
	if result[0].MainImageURL != "https://example1.com/image1.jpg" {
		t.Errorf("GetAllBookmarks() first bookmark MainImageURL = %v, want https://example1.com/image1.jpg", result[0].MainImageURL)
	}
	if result[0].ContentSummary != "First example website content." {
		t.Errorf("GetAllBookmarks() first bookmark ContentSummary = %v, want First example website content.", result[0].ContentSummary)
	}
	if result[1].ID != "bookmark-2" {
		t.Errorf("GetAllBookmarks() second bookmark ID = %v, want bookmark-2", result[1].ID)
	}
	if result[1].MainImageURL != "https://example2.com/image2.jpg" {
		t.Errorf("GetAllBookmarks() second bookmark MainImageURL = %v, want https://example2.com/image2.jpg", result[1].MainImageURL)
	}
	if result[1].ContentSummary != "Second example website content." {
		t.Errorf("GetAllBookmarks() second bookmark ContentSummary = %v, want Second example website content.", result[1].ContentSummary)
	}
}

// TestBookmarkService_GetAllBookmarks_EmptyUserID tests getting bookmarks with empty userID
func TestBookmarkService_GetAllBookmarks_EmptyUserID(t *testing.T) {
	expectedBookmarks := []model.Bookmark{}

	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			if userID != "" {
				t.Errorf("ListBookmarks() received UserID = %v, want empty string", userID)
			}
			return expectedBookmarks, nil
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
	result, err := service.GetAllBookmarks("", false)

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
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			return nil, fmt.Errorf("database connection failed")
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
	_, err := service.GetAllBookmarks("user-1", false)

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
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			return []model.Bookmark{}, nil
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
	result, err := service.GetAllBookmarks("user-1", false)

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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
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
	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)

	// Verify service is usable by calling a method
	result, err := service.GetBookmark("test-id")
	if err != nil {
		t.Errorf("NewBookmarkService() service should be functional, got error = %v", err)
	}
	if result.ID != "test-id" {
		t.Error("NewBookmarkService() service should work correctly")
	}
}

// TestBookmarkService_GetBookmarksWithPagination tests successful paginated retrieval
func TestBookmarkService_GetBookmarksWithPagination(t *testing.T) {
	expectedBookmarks := []model.Bookmark{
		{ID: "bookmark-1", UserID: "user-1", URL: "https://example1.com", Title: "Example 1", MainImageURL: "https://example1.com/og1.jpg", ContentSummary: "Summary for example 1."},
		{ID: "bookmark-2", UserID: "user-1", URL: "https://example2.com", Title: "Example 2", MainImageURL: "https://example2.com/og2.jpg", ContentSummary: "Summary for example 2."},
	}

	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			return expectedBookmarks, nil
		},
		countBookmarksFunc: func(query model.BookmarkQuery) (int, error) {
			return 10, nil
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
	result, err := service.GetBookmarksWithPagination("user-1", false, 1, 20)

	if err != nil {
		t.Errorf("GetBookmarksWithPagination() unexpected error = %v", err)
		return
	}

	if len(result.Bookmarks) != 2 {
		t.Errorf("GetBookmarksWithPagination() returned %d bookmarks, want 2", len(result.Bookmarks))
	}
	if result.Bookmarks[0].MainImageURL != "https://example1.com/og1.jpg" {
		t.Errorf("GetBookmarksWithPagination() first bookmark MainImageURL = %v, want https://example1.com/og1.jpg", result.Bookmarks[0].MainImageURL)
	}
	if result.Bookmarks[0].ContentSummary != "Summary for example 1." {
		t.Errorf("GetBookmarksWithPagination() first bookmark ContentSummary = %v, want Summary for example 1.", result.Bookmarks[0].ContentSummary)
	}
	if result.Bookmarks[1].MainImageURL != "https://example2.com/og2.jpg" {
		t.Errorf("GetBookmarksWithPagination() second bookmark MainImageURL = %v, want https://example2.com/og2.jpg", result.Bookmarks[1].MainImageURL)
	}
	if result.Bookmarks[1].ContentSummary != "Summary for example 2." {
		t.Errorf("GetBookmarksWithPagination() second bookmark ContentSummary = %v, want Summary for example 2.", result.Bookmarks[1].ContentSummary)
	}
	if result.TotalCount != 10 {
		t.Errorf("GetBookmarksWithPagination() TotalCount = %d, want 10", result.TotalCount)
	}
	if result.Page != 1 {
		t.Errorf("GetBookmarksWithPagination() Page = %d, want 1", result.Page)
	}
	if result.PageSize != 20 {
		t.Errorf("GetBookmarksWithPagination() PageSize = %d, want 20", result.PageSize)
	}
	if result.TotalPages != 1 {
		t.Errorf("GetBookmarksWithPagination() TotalPages = %d, want 1", result.TotalPages)
	}
}

// TestBookmarkService_GetBookmarksWithPagination_PageValidation tests page parameter validation
func TestBookmarkService_GetBookmarksWithPagination_PageValidation(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			return []model.Bookmark{}, nil
		},
		countBookmarksFunc: func(query model.BookmarkQuery) (int, error) {
			return 0, nil
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)

	// Test with page < 1 (should default to 1)
	result, err := service.GetBookmarksWithPagination("user-1", false, 0, 20)
	if err != nil {
		t.Errorf("GetBookmarksWithPagination() unexpected error = %v", err)
		return
	}
	if result.Page != 1 {
		t.Errorf("GetBookmarksWithPagination() with page=0 should default to 1, got %d", result.Page)
	}

	// Test with negative page (should default to 1)
	result, err = service.GetBookmarksWithPagination("user-1", false, -5, 20)
	if err != nil {
		t.Errorf("GetBookmarksWithPagination() unexpected error = %v", err)
		return
	}
	if result.Page != 1 {
		t.Errorf("GetBookmarksWithPagination() with page=-5 should default to 1, got %d", result.Page)
	}
}

// TestBookmarkService_GetBookmarksWithPagination_PageSizeValidation tests pageSize parameter validation
func TestBookmarkService_GetBookmarksWithPagination_PageSizeValidation(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			return []model.Bookmark{}, nil
		},
		countBookmarksFunc: func(query model.BookmarkQuery) (int, error) {
			return 0, nil
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)

	// Test with pageSize < 1 (should default to 20)
	result, err := service.GetBookmarksWithPagination("user-1", false, 1, 0)
	if err != nil {
		t.Errorf("GetBookmarksWithPagination() unexpected error = %v", err)
		return
	}
	if result.PageSize != 20 {
		t.Errorf("GetBookmarksWithPagination() with pageSize=0 should default to 20, got %d", result.PageSize)
	}

	// Test with negative pageSize (should default to 20)
	result, err = service.GetBookmarksWithPagination("user-1", false, 1, -10)
	if err != nil {
		t.Errorf("GetBookmarksWithPagination() unexpected error = %v", err)
		return
	}
	if result.PageSize != 20 {
		t.Errorf("GetBookmarksWithPagination() with pageSize=-10 should default to 20, got %d", result.PageSize)
	}

	// Test with pageSize > 100 (should cap at 100)
	result, err = service.GetBookmarksWithPagination("user-1", false, 1, 200)
	if err != nil {
		t.Errorf("GetBookmarksWithPagination() unexpected error = %v", err)
		return
	}
	if result.PageSize != 100 {
		t.Errorf("GetBookmarksWithPagination() with pageSize=200 should cap at 100, got %d", result.PageSize)
	}
}

// TestBookmarkService_GetBookmarksWithPagination_TotalPagesCalculation tests total pages calculation
func TestBookmarkService_GetBookmarksWithPagination_TotalPagesCalculation(t *testing.T) {
	testCases := []struct {
		name       string
		totalCount int
		pageSize   int
		wantPages  int
	}{
		{"Exact multiple", 100, 20, 5},
		{"Not exact multiple", 95, 20, 5},
		{"Less than page size", 15, 20, 1},
		{"Zero count", 0, 20, 1},
		{"One item", 1, 20, 1},
		{"Edge case 21 items", 21, 20, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockBookmarkRepository{
				listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
					return []model.Bookmark{}, nil
				},
				countBookmarksFunc: func(query model.BookmarkQuery) (int, error) {
					return tc.totalCount, nil
				},
			}

			mockWebRepo := &MockWebRepository{}
			service := NewBookmarkService(mockRepo, mockWebRepo)
			result, err := service.GetBookmarksWithPagination("user-1", false, 1, tc.pageSize)

			if err != nil {
				t.Errorf("GetBookmarksWithPagination() unexpected error = %v", err)
				return
			}

			if result.TotalPages != tc.wantPages {
				t.Errorf("GetBookmarksWithPagination() TotalPages = %d, want %d", result.TotalPages, tc.wantPages)
			}
		})
	}
}

// TestBookmarkService_GetBookmarksWithPagination_EmptyResult tests handling empty result
func TestBookmarkService_GetBookmarksWithPagination_EmptyResult(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			return []model.Bookmark{}, nil
		},
		countBookmarksFunc: func(query model.BookmarkQuery) (int, error) {
			return 0, nil
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
	result, err := service.GetBookmarksWithPagination("user-1", false, 1, 20)

	if err != nil {
		t.Errorf("GetBookmarksWithPagination() with empty result unexpected error = %v", err)
		return
	}

	if len(result.Bookmarks) != 0 {
		t.Errorf("GetBookmarksWithPagination() returned %d bookmarks, want 0", len(result.Bookmarks))
	}
	if result.TotalCount != 0 {
		t.Errorf("GetBookmarksWithPagination() TotalCount = %d, want 0", result.TotalCount)
	}
	if result.TotalPages != 1 {
		t.Errorf("GetBookmarksWithPagination() TotalPages = %d, want 1", result.TotalPages)
	}
}

// TestBookmarkService_GetBookmarksWithPagination_ListError tests error handling when listing fails
func TestBookmarkService_GetBookmarksWithPagination_ListError(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			return nil, fmt.Errorf("database connection failed")
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
	_, err := service.GetBookmarksWithPagination("user-1", false, 1, 20)

	if err == nil {
		t.Error("GetBookmarksWithPagination() should return error when list fails")
		return
	}

	expectedErrorSubstring := "failed to get bookmarks"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("GetBookmarksWithPagination() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_GetBookmarksWithPagination_CountError tests error handling when count fails
func TestBookmarkService_GetBookmarksWithPagination_CountError(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			return []model.Bookmark{}, nil
		},
		countBookmarksFunc: func(query model.BookmarkQuery) (int, error) {
			return 0, fmt.Errorf("database count failed")
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)
	_, err := service.GetBookmarksWithPagination("user-1", false, 1, 20)

	if err == nil {
		t.Error("GetBookmarksWithPagination() should return error when count fails")
		return
	}

	expectedErrorSubstring := "failed to count bookmarks"
	if len(err.Error()) < len(expectedErrorSubstring) || err.Error()[:len(expectedErrorSubstring)] != expectedErrorSubstring {
		t.Errorf("GetBookmarksWithPagination() error should contain '%s', got %v", expectedErrorSubstring, err.Error())
	}
}

// TestBookmarkService_GetBookmarksWithPagination_ArchivedFilter tests archived filtering
func TestBookmarkService_GetBookmarksWithPagination_ArchivedFilter(t *testing.T) {
	mockRepo := &MockBookmarkRepository{
		listBookmarksFunc: func(userID string, archived bool) ([]model.Bookmark, error) {
			if archived {
				return []model.Bookmark{
					{ID: "archived-1", IsArchived: true},
				}, nil
			}
			return []model.Bookmark{
				{ID: "active-1", IsArchived: false},
			}, nil
		},
		countBookmarksFunc: func(query model.BookmarkQuery) (int, error) {
			if query.Archived {
				return 1, nil
			}
			return 1, nil
		},
	}

	mockWebRepo := &MockWebRepository{}
	service := NewBookmarkService(mockRepo, mockWebRepo)

	// Test with archived = false
	result, err := service.GetBookmarksWithPagination("user-1", false, 1, 20)
	if err != nil {
		t.Errorf("GetBookmarksWithPagination() with archived=false unexpected error = %v", err)
		return
	}
	if len(result.Bookmarks) != 1 || result.Bookmarks[0].ID != "active-1" {
		t.Error("GetBookmarksWithPagination() with archived=false should return active bookmarks")
	}

	// Test with archived = true
	result, err = service.GetBookmarksWithPagination("user-1", true, 1, 20)
	if err != nil {
		t.Errorf("GetBookmarksWithPagination() with archived=true unexpected error = %v", err)
		return
	}
	if len(result.Bookmarks) != 1 || result.Bookmarks[0].ID != "archived-1" {
		t.Error("GetBookmarksWithPagination() with archived=true should return archived bookmarks")
	}
}
