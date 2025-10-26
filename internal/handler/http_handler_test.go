package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tsongpon/athena/internal/model"
)

// MockBookmarkService is a mock implementation of BookmarkService interface for testing
type MockBookmarkService struct {
	createBookmarkFunc  func(bookmark model.Bookmark) (model.Bookmark, error)
	archiveBookmarkFunc func(id string) (model.Bookmark, error)
	getBookmarkFunc     func(id string) (model.Bookmark, error)
	getAllBookmarksFunc func(userID string) ([]model.Bookmark, error)
	deleteBookmarkFunc  func(id string) error
}

func (m *MockBookmarkService) CreateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
	if m.createBookmarkFunc != nil {
		return m.createBookmarkFunc(bookmark)
	}
	return model.Bookmark{}, nil
}

func (m *MockBookmarkService) ArchiveBookmark(id string) (model.Bookmark, error) {
	if m.archiveBookmarkFunc != nil {
		return m.archiveBookmarkFunc(id)
	}
	return model.Bookmark{}, nil
}

func (m *MockBookmarkService) GetBookmark(id string) (model.Bookmark, error) {
	if m.getBookmarkFunc != nil {
		return m.getBookmarkFunc(id)
	}
	return model.Bookmark{}, nil
}

func (m *MockBookmarkService) GetAllBookmarks(userID string) ([]model.Bookmark, error) {
	if m.getAllBookmarksFunc != nil {
		return m.getAllBookmarksFunc(userID)
	}
	return []model.Bookmark{}, nil
}

func (m *MockBookmarkService) DeleteBookmark(id string) error {
	if m.deleteBookmarkFunc != nil {
		return m.deleteBookmarkFunc(id)
	}
	return nil
}

// TestHTTPHandler_Ping tests the Ping endpoint
func TestHTTPHandler_Ping(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &MockBookmarkService{}
	handler := NewHTTPHandler(mockService)

	err := handler.Ping(c)
	if err != nil {
		t.Errorf("Ping() unexpected error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Ping() status code = %d, want %d", rec.Code, http.StatusOK)
	}

	expectedBody := "pong"
	if rec.Body.String() != expectedBody {
		t.Errorf("Ping() body = %v, want %v", rec.Body.String(), expectedBody)
	}
}

// TestHTTPHandler_CreateBookmark tests successful bookmark creation
func TestHTTPHandler_CreateBookmark(t *testing.T) {
	e := echo.New()
	requestBody := `{"url":"https://example.com","user_id":"user-123"}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	createdBookmark := model.Bookmark{
		ID:        "bookmark-1",
		URL:       "https://example.com",
		UserID:    "user-123",
		CreatedAt: time.Now(),
	}

	mockService := &MockBookmarkService{
		createBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			if bookmark.URL != "https://example.com" {
				t.Errorf("CreateBookmark() received URL = %v, want https://example.com", bookmark.URL)
			}
			if bookmark.UserID != "user-123" {
				t.Errorf("CreateBookmark() received UserID = %v, want user-123", bookmark.UserID)
			}
			if bookmark.IsArchived != false {
				t.Error("CreateBookmark() should set IsArchived to false")
			}
			return createdBookmark, nil
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.CreateBookmark(c)

	if err != nil {
		t.Errorf("CreateBookmark() unexpected error = %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("CreateBookmark() status code = %d, want %d", rec.Code, http.StatusCreated)
	}
}

// TestHTTPHandler_CreateBookmark_InvalidJSON tests creating bookmark with invalid JSON
func TestHTTPHandler_CreateBookmark_InvalidJSON(t *testing.T) {
	e := echo.New()
	requestBody := `{"url":"https://example.com","user_id":}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &MockBookmarkService{
		createBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			t.Error("CreateBookmark() should not be called with invalid JSON")
			return model.Bookmark{}, nil
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.CreateBookmark(c)

	if err == nil {
		t.Error("CreateBookmark() should return error with invalid JSON")
	}
}

// TestHTTPHandler_CreateBookmark_ServiceError tests error handling from service
func TestHTTPHandler_CreateBookmark_ServiceError(t *testing.T) {
	e := echo.New()
	requestBody := `{"url":"https://example.com","user_id":"user-123"}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &MockBookmarkService{
		createBookmarkFunc: func(bookmark model.Bookmark) (model.Bookmark, error) {
			return model.Bookmark{}, fmt.Errorf("service error")
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.CreateBookmark(c)

	if err == nil {
		t.Error("CreateBookmark() should return error when service fails")
	}
}

// TestHTTPHandler_GetBookmark tests successful bookmark retrieval
func TestHTTPHandler_GetBookmark(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks/bookmark-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark-1")

	expectedBookmark := model.Bookmark{
		ID:        "bookmark-1",
		URL:       "https://example.com",
		Title:     "Example",
		UserID:    "user-123",
		CreatedAt: time.Now(),
	}

	mockService := &MockBookmarkService{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			if id != "bookmark-1" {
				t.Errorf("GetBookmark() received ID = %v, want bookmark-1", id)
			}
			return expectedBookmark, nil
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.GetBookmark(c)

	if err != nil {
		t.Errorf("GetBookmark() unexpected error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("GetBookmark() status code = %d, want %d", rec.Code, http.StatusOK)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["id"] != "bookmark-1" {
		t.Errorf("GetBookmark() response ID = %v, want bookmark-1", response["id"])
	}
	if response["url"] != "https://example.com" {
		t.Errorf("GetBookmark() response URL = %v, want https://example.com", response["url"])
	}
}

// TestHTTPHandler_GetBookmark_EmptyID tests getting bookmark with empty ID
func TestHTTPHandler_GetBookmark_EmptyID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("")

	mockService := &MockBookmarkService{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			t.Error("GetBookmark() should not be called with empty ID")
			return model.Bookmark{}, nil
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.GetBookmark(c)

	if err == nil {
		t.Error("GetBookmark() should return error with empty ID")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Error("GetBookmark() should return HTTPError")
	} else if httpErr.Code != http.StatusBadRequest {
		t.Errorf("GetBookmark() error code = %d, want %d", httpErr.Code, http.StatusBadRequest)
	}
}

// TestHTTPHandler_GetBookmark_ServiceError tests error handling from service
func TestHTTPHandler_GetBookmark_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks/bookmark-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark-1")

	mockService := &MockBookmarkService{
		getBookmarkFunc: func(id string) (model.Bookmark, error) {
			return model.Bookmark{}, fmt.Errorf("bookmark not found")
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.GetBookmark(c)

	if err == nil {
		t.Error("GetBookmark() should return error when service fails")
	}
}

// TestHTTPHandler_GetBookmarks tests successful retrieval of all bookmarks
func TestHTTPHandler_GetBookmarks(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks?userid=user-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedBookmarks := []model.Bookmark{
		{
			ID:        "bookmark-1",
			URL:       "https://example1.com",
			Title:     "Example 1",
			UserID:    "user-123",
			CreatedAt: time.Now(),
		},
		{
			ID:        "bookmark-2",
			URL:       "https://example2.com",
			Title:     "Example 2",
			UserID:    "user-123",
			CreatedAt: time.Now(),
		},
	}

	mockService := &MockBookmarkService{
		getAllBookmarksFunc: func(userID string) ([]model.Bookmark, error) {
			if userID != "user-123" {
				t.Errorf("GetAllBookmarks() received UserID = %v, want user-123", userID)
			}
			return expectedBookmarks, nil
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.GetBookmarks(c)

	if err != nil {
		t.Errorf("GetBookmarks() unexpected error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("GetBookmarks() status code = %d, want %d", rec.Code, http.StatusOK)
	}

	// Verify response body
	var response []map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("GetBookmarks() returned %d bookmarks, want 2", len(response))
	}

	if response[0]["id"] != "bookmark-1" {
		t.Errorf("GetBookmarks() first bookmark ID = %v, want bookmark-1", response[0]["id"])
	}
}

// TestHTTPHandler_GetBookmarks_EmptyUserID tests getting bookmarks with empty userID
func TestHTTPHandler_GetBookmarks_EmptyUserID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &MockBookmarkService{
		getAllBookmarksFunc: func(userID string) ([]model.Bookmark, error) {
			if userID != "" {
				t.Errorf("GetAllBookmarks() received UserID = %v, want empty string", userID)
			}
			return []model.Bookmark{}, nil
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.GetBookmarks(c)

	if err != nil {
		t.Errorf("GetBookmarks() unexpected error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("GetBookmarks() status code = %d, want %d", rec.Code, http.StatusOK)
	}
}

// TestHTTPHandler_GetBookmarks_ServiceError tests error handling from service
func TestHTTPHandler_GetBookmarks_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks?userid=user-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &MockBookmarkService{
		getAllBookmarksFunc: func(userID string) ([]model.Bookmark, error) {
			return nil, fmt.Errorf("database error")
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.GetBookmarks(c)

	if err == nil {
		t.Error("GetBookmarks() should return error when service fails")
	}
}

// TestHTTPHandler_GetBookmarks_EmptyResult tests getting bookmarks with no results
func TestHTTPHandler_GetBookmarks_EmptyResult(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks?userid=user-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &MockBookmarkService{
		getAllBookmarksFunc: func(userID string) ([]model.Bookmark, error) {
			return []model.Bookmark{}, nil
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.GetBookmarks(c)

	if err != nil {
		t.Errorf("GetBookmarks() unexpected error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("GetBookmarks() status code = %d, want %d", rec.Code, http.StatusOK)
	}

	// Verify response is empty array
	var response []map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("GetBookmarks() returned %d bookmarks, want 0", len(response))
	}
}

// TestHTTPHandler_ArchiveBookmark tests successful bookmark archiving
func TestHTTPHandler_ArchiveBookmark(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bookmarks/bookmark-1/archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark-1")

	archivedBookmark := model.Bookmark{
		ID:         "bookmark-1",
		URL:        "https://example.com",
		UserID:     "user-123",
		IsArchived: true,
		CreatedAt:  time.Now(),
	}

	mockService := &MockBookmarkService{
		archiveBookmarkFunc: func(id string) (model.Bookmark, error) {
			if id != "bookmark-1" {
				t.Errorf("ArchiveBookmark() received ID = %v, want bookmark-1", id)
			}
			return archivedBookmark, nil
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.ArchiveBookmark(c)

	if err != nil {
		t.Errorf("ArchiveBookmark() unexpected error = %v", err)
	}

	if rec.Code != http.StatusNoContent {
		t.Errorf("ArchiveBookmark() status code = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

// TestHTTPHandler_ArchiveBookmark_EmptyID tests archiving bookmark with empty ID
func TestHTTPHandler_ArchiveBookmark_EmptyID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bookmarks//archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("")

	mockService := &MockBookmarkService{
		archiveBookmarkFunc: func(id string) (model.Bookmark, error) {
			t.Error("ArchiveBookmark() should not be called with empty ID")
			return model.Bookmark{}, nil
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.ArchiveBookmark(c)

	if err == nil {
		t.Error("ArchiveBookmark() should return error with empty ID")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Error("ArchiveBookmark() should return HTTPError")
	} else if httpErr.Code != http.StatusBadRequest {
		t.Errorf("ArchiveBookmark() error code = %d, want %d", httpErr.Code, http.StatusBadRequest)
	}
}

// TestHTTPHandler_ArchiveBookmark_ServiceError tests error handling from service
func TestHTTPHandler_ArchiveBookmark_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bookmarks/bookmark-1/archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark-1")

	mockService := &MockBookmarkService{
		archiveBookmarkFunc: func(id string) (model.Bookmark, error) {
			return model.Bookmark{}, fmt.Errorf("bookmark not found")
		},
	}

	handler := NewHTTPHandler(mockService)
	err := handler.ArchiveBookmark(c)

	if err == nil {
		t.Error("ArchiveBookmark() should return error when service fails")
	}
}

// TestNewHTTPHandler tests handler initialization
func TestNewHTTPHandler(t *testing.T) {
	mockService := &MockBookmarkService{}
	handler := NewHTTPHandler(mockService)

	// Test that handler can be called successfully
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Ping(c)
	if err != nil {
		t.Errorf("NewHTTPHandler() initialized handler should work, got error = %v", err)
	}
}
