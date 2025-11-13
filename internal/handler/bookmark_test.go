package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsongpon/athena/internal/model"
	"github.com/tsongpon/athena/internal/transport"
)

// MockBookmarkService is a mock implementation of BookmarkService
type MockBookmarkService struct {
	mock.Mock
}

func (m *MockBookmarkService) CreateBookmark(b model.Bookmark) (model.Bookmark, error) {
	args := m.Called(b)
	return args.Get(0).(model.Bookmark), args.Error(1)
}

func (m *MockBookmarkService) GetBookmark(id string) (model.Bookmark, error) {
	args := m.Called(id)
	return args.Get(0).(model.Bookmark), args.Error(1)
}

func (m *MockBookmarkService) DeleteBookmark(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookmarkService) GetAllBookmarks(userID string, archived bool) ([]model.Bookmark, error) {
	args := m.Called(userID, archived)
	return args.Get(0).([]model.Bookmark), args.Error(1)
}

func (m *MockBookmarkService) GetBookmarksWithPagination(userID string, archived bool, page, pageSize int) (model.BookmarkListResponse, error) {
	args := m.Called(userID, archived, page, pageSize)
	return args.Get(0).(model.BookmarkListResponse), args.Error(1)
}

func (m *MockBookmarkService) ArchiveBookmark(id string) (model.Bookmark, error) {
	args := m.Called(id)
	return args.Get(0).(model.Bookmark), args.Error(1)
}

func TestNewBookmarkHandler(t *testing.T) {
	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.bookmarkService)
}

func TestBookmarkHandler_Ping(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.Ping(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "pong", rec.Body.String())
}

func TestBookmarkHandler_CreateBookmark_Success(t *testing.T) {
	e := echo.New()
	bookmarkJSON := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmark := model.Bookmark{
		ID:           "bookmark123",
		URL:          "https://example.com",
		Title:        "Example Site",
		UserID:       "user123",
		MainImageURL: "https://example.com/image.png",
		IsArchived:   false,
		CreatedAt:    time.Now(),
	}

	mockService.On("CreateBookmark", mock.MatchedBy(func(b model.Bookmark) bool {
		return b.URL == "https://example.com" && b.UserID == "user123"
	})).Return(expectedBookmark, nil)

	err := handler.CreateBookmark(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var responseTransport transport.BookmarkTransport
	err = json.Unmarshal(rec.Body.Bytes(), &responseTransport)
	assert.NoError(t, err)
	assert.Equal(t, expectedBookmark.ID, responseTransport.ID)
	assert.Equal(t, expectedBookmark.URL, responseTransport.URL)
	assert.Equal(t, expectedBookmark.Title, responseTransport.Title)
	assert.Equal(t, expectedBookmark.UserID, responseTransport.UserID)
	assert.Equal(t, expectedBookmark.MainImageURL, responseTransport.MainImageURL)
	assert.Equal(t, expectedBookmark.IsArchived, responseTransport.IsArchived)

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_CreateBookmark_MissingURL(t *testing.T) {
	e := echo.New()
	bookmarkJSON := `{}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.CreateBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "URL is required", httpErr.Message)
}

func TestBookmarkHandler_CreateBookmark_MissingAuthentication(t *testing.T) {
	e := echo.New()
	bookmarkJSON := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Do not set user_id in context to simulate missing authentication

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.CreateBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Equal(t, "User not authenticated", httpErr.Message)
}

func TestBookmarkHandler_CreateBookmark_ServiceError(t *testing.T) {
	e := echo.New()
	bookmarkJSON := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	mockService.On("CreateBookmark", mock.Anything).Return(model.Bookmark{}, errors.New("database error"))

	err := handler.CreateBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_CreateBookmark_InvalidJSON(t *testing.T) {
	e := echo.New()
	bookmarkJSON := `{"url":"https://example.com","invalid":}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.CreateBookmark(c)

	assert.Error(t, err)
}

func TestBookmarkHandler_GetBookmark_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks/bookmark123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmark := model.Bookmark{
		ID:           "bookmark123",
		URL:          "https://example.com",
		Title:        "Example Site",
		UserID:       "user123",
		MainImageURL: "https://example.com/og-image.jpg",
		IsArchived:   false,
		CreatedAt:    time.Now(),
	}

	mockService.On("GetBookmark", "bookmark123").Return(expectedBookmark, nil)

	err := handler.GetBookmark(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseTransport transport.BookmarkTransport
	err = json.Unmarshal(rec.Body.Bytes(), &responseTransport)
	assert.NoError(t, err)
	assert.Equal(t, expectedBookmark.ID, responseTransport.ID)
	assert.Equal(t, expectedBookmark.URL, responseTransport.URL)
	assert.Equal(t, expectedBookmark.Title, responseTransport.Title)
	assert.Equal(t, expectedBookmark.UserID, responseTransport.UserID)
	assert.Equal(t, expectedBookmark.MainImageURL, responseTransport.MainImageURL)

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmark_MissingID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.GetBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "ID is required", httpErr.Message)
}

func TestBookmarkHandler_GetBookmark_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks/bookmark123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	mockService.On("GetBookmark", "bookmark123").Return(model.Bookmark{}, errors.New("not found"))

	err := handler.GetBookmark(c)

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmarks_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks?archived=false", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmarks := []model.Bookmark{
		{
			ID:           "bookmark1",
			URL:          "https://example1.com",
			Title:        "Example 1",
			UserID:       "user123",
			MainImageURL: "https://example1.com/image1.png",
			IsArchived:   false,
			CreatedAt:    time.Now(),
		},
		{
			ID:           "bookmark2",
			URL:          "https://example2.com",
			Title:        "Example 2",
			UserID:       "user123",
			MainImageURL: "https://example2.com/image2.png",
			IsArchived:   false,
			CreatedAt:    time.Now(),
		},
	}

	mockService.On("GetAllBookmarks", "user123", false).Return(expectedBookmarks, nil)

	err := handler.GetBookmarks(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseTransports []transport.BookmarkTransport
	err = json.Unmarshal(rec.Body.Bytes(), &responseTransports)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(responseTransports))
	assert.Equal(t, expectedBookmarks[0].ID, responseTransports[0].ID)
	assert.Equal(t, expectedBookmarks[0].MainImageURL, responseTransports[0].MainImageURL)
	assert.Equal(t, expectedBookmarks[1].ID, responseTransports[1].ID)
	assert.Equal(t, expectedBookmarks[1].MainImageURL, responseTransports[1].MainImageURL)

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmarks_ArchivedTrue(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks?archived=true", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmarks := []model.Bookmark{
		{
			ID:           "bookmark1",
			URL:          "https://example1.com",
			Title:        "Example 1",
			UserID:       "user123",
			MainImageURL: "https://example1.com/archived-image.png",
			IsArchived:   true,
			CreatedAt:    time.Now(),
		},
	}

	mockService.On("GetAllBookmarks", "user123", true).Return(expectedBookmarks, nil)

	err := handler.GetBookmarks(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseTransports []transport.BookmarkTransport
	err = json.Unmarshal(rec.Body.Bytes(), &responseTransports)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(responseTransports))
	assert.Equal(t, true, responseTransports[0].IsArchived)

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmarks_DefaultArchivedValue(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	mockService.On("GetAllBookmarks", "user123", false).Return([]model.Bookmark{}, nil)

	err := handler.GetBookmarks(c)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmarks_InvalidArchivedParam(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks?archived=invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	// Should default to false when parsing fails
	mockService.On("GetAllBookmarks", "user123", false).Return([]model.Bookmark{}, nil)

	err := handler.GetBookmarks(c)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmarks_MissingAuthentication(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Do not set user in context to simulate missing authentication

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.GetBookmarks(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Equal(t, "User not authenticated", httpErr.Message)
}

func TestBookmarkHandler_GetBookmarks_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	mockService.On("GetAllBookmarks", "user123", false).Return([]model.Bookmark{}, errors.New("database error"))

	err := handler.GetBookmarks(c)

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmarks_EmptyResult(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	mockService.On("GetAllBookmarks", "user123", false).Return([]model.Bookmark{}, nil)

	err := handler.GetBookmarks(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseTransports []transport.BookmarkTransport
	err = json.Unmarshal(rec.Body.Bytes(), &responseTransports)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(responseTransports))

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_ArchiveBookmark_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/bookmarks/bookmark123/archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	existingBookmark := model.Bookmark{
		ID:         "bookmark123",
		URL:        "https://example.com",
		Title:      "Example Site",
		UserID:     "user123",
		IsArchived: false,
	}

	archivedBookmark := model.Bookmark{
		ID:         "bookmark123",
		URL:        "https://example.com",
		Title:      "Example Site",
		UserID:     "user123",
		IsArchived: true,
	}

	mockService.On("GetBookmark", "bookmark123").Return(existingBookmark, nil)
	mockService.On("ArchiveBookmark", "bookmark123").Return(archivedBookmark, nil)

	err := handler.ArchiveBookmark(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "", rec.Body.String())

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_ArchiveBookmark_MissingID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/bookmarks//archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.ArchiveBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "ID is required", httpErr.Message)
}

func TestBookmarkHandler_ArchiveBookmark_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/bookmarks/bookmark123/archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	mockService.On("GetBookmark", "bookmark123").Return(model.Bookmark{}, errors.New("bookmark not found"))

	err := handler.ArchiveBookmark(c)

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

// Additional edge case tests

func TestBookmarkHandler_CreateBookmark_WithEmptyStringURL(t *testing.T) {
	e := echo.New()
	bookmarkJSON := `{"url":""}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.CreateBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "URL is required", httpErr.Message)
}

func TestBookmarkHandler_CreateBookmark_WithWhitespaceURL(t *testing.T) {
	e := echo.New()
	bookmarkJSON := `{"url":"   "}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmark := model.Bookmark{
		ID:         "bookmark123",
		URL:        "   ",
		UserID:     "user123",
		IsArchived: false,
	}

	mockService.On("CreateBookmark", mock.Anything).Return(expectedBookmark, nil)

	err := handler.CreateBookmark(c)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_CreateBookmark_NullUserContext(t *testing.T) {
	e := echo.New()
	bookmarkJSON := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user", nil)

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.CreateBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Equal(t, "User not authenticated", httpErr.Message)
}

func TestBookmarkHandler_CreateBookmark_WrongTypeInContext(t *testing.T) {
	e := echo.New()
	bookmarkJSON := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user", "not-a-jwt-claims-object")

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.CreateBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

func TestBookmarkHandler_CreateBookmark_LongURL(t *testing.T) {
	e := echo.New()
	longURL := "https://example.com/" + strings.Repeat("a", 2000)
	bookmarkJSON := `{"url":"` + longURL + `"}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmark := model.Bookmark{
		ID:         "bookmark123",
		URL:        longURL,
		UserID:     "user123",
		IsArchived: false,
	}

	mockService.On("CreateBookmark", mock.Anything).Return(expectedBookmark, nil)

	err := handler.CreateBookmark(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_CreateBookmark_SpecialCharactersInURL(t *testing.T) {
	e := echo.New()
	specialURL := "https://example.com/path?query=value&foo=bar#fragment"
	bookmarkJSON := `{"url":"` + specialURL + `"}`
	req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmark := model.Bookmark{
		ID:         "bookmark123",
		URL:        specialURL,
		UserID:     "user123",
		IsArchived: false,
	}

	mockService.On("CreateBookmark", mock.Anything).Return(expectedBookmark, nil)

	err := handler.CreateBookmark(c)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmark_NotFoundError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	mockService.On("GetBookmark", "nonexistent").Return(model.Bookmark{}, errors.New("bookmark not found"))

	err := handler.GetBookmark(c)

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmarks_CaseInsensitiveArchived(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks?archived=TRUE", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmarks := []model.Bookmark{
		{
			ID:         "bookmark1",
			URL:        "https://example.com",
			Title:      "Example",
			UserID:     "user123",
			IsArchived: true,
		},
	}

	mockService.On("GetAllBookmarks", "user123", true).Return(expectedBookmarks, nil)

	err := handler.GetBookmarks(c)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmarks_WithNumericArchived(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks?archived=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmarks := []model.Bookmark{}

	mockService.On("GetAllBookmarks", "user123", true).Return(expectedBookmarks, nil)

	err := handler.GetBookmarks(c)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmarks_LargeResultSet(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmarks := make([]model.Bookmark, 1000)
	for i := 0; i < 1000; i++ {
		expectedBookmarks[i] = model.Bookmark{
			ID:         "bookmark" + string(rune(i)),
			URL:        "https://example.com",
			UserID:     "user123",
			IsArchived: false,
		}
	}

	mockService.On("GetAllBookmarks", "user123", false).Return(expectedBookmarks, nil)

	err := handler.GetBookmarks(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseTransports []transport.BookmarkTransport
	err = json.Unmarshal(rec.Body.Bytes(), &responseTransports)
	assert.NoError(t, err)
	assert.Equal(t, 1000, len(responseTransports))

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_GetBookmarks_NullUserContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set null user in context
	c.Set("user", nil)

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.GetBookmarks(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Equal(t, "User not authenticated", httpErr.Message)
}

func TestBookmarkHandler_ArchiveBookmark_AlreadyArchived(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/bookmarks/bookmark123/archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	existingArchivedBookmark := model.Bookmark{
		ID:         "bookmark123",
		URL:        "https://example.com",
		Title:      "Example Site",
		UserID:     "user123",
		IsArchived: true,
	}

	mockService.On("GetBookmark", "bookmark123").Return(existingArchivedBookmark, nil)
	mockService.On("ArchiveBookmark", "bookmark123").Return(existingArchivedBookmark, nil)

	err := handler.ArchiveBookmark(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_ArchiveBookmark_NotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/bookmarks/nonexistent/archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	mockService.On("GetBookmark", "nonexistent").Return(model.Bookmark{}, errors.New("bookmark not found"))

	err := handler.ArchiveBookmark(c)

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_ArchiveBookmark_EmptyID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/bookmarks//archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("")

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.ArchiveBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "ID is required", httpErr.Message)
}

func TestBookmarkHandler_Ping_MultipleRequests(t *testing.T) {
	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	for i := 0; i < 5; i++ {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Ping(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "pong", rec.Body.String())
	}
}

func TestBookmarkHandler_GetBookmark_MissingAuthentication(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks/bookmark123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Do not set user in context to simulate missing authentication

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.GetBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Equal(t, "User not authenticated", httpErr.Message)
}

func TestBookmarkHandler_GetBookmark_Forbidden_DifferentUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bookmarks/bookmark123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// User user456 trying to access user123's bookmark
	c.Set("user", &JWTClaims{UserID: "user456", Email: "other@example.com", Name: "Other User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	bookmarkOwnedByDifferentUser := model.Bookmark{
		ID:         "bookmark123",
		URL:        "https://example.com",
		Title:      "Example Site",
		UserID:     "user123", // Owned by different user
		IsArchived: false,
	}

	mockService.On("GetBookmark", "bookmark123").Return(bookmarkOwnedByDifferentUser, nil)

	err := handler.GetBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpErr.Code)
	assert.Equal(t, "Access denied", httpErr.Message)

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_ArchiveBookmark_MissingAuthentication(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/bookmarks/bookmark123/archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Do not set user in context to simulate missing authentication

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.ArchiveBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Equal(t, "User not authenticated", httpErr.Message)
}

func TestBookmarkHandler_ArchiveBookmark_Forbidden_DifferentUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/bookmarks/bookmark123/archive", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// User user456 trying to archive user123's bookmark
	c.Set("user", &JWTClaims{UserID: "user456", Email: "other@example.com", Name: "Other User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	bookmarkOwnedByDifferentUser := model.Bookmark{
		ID:         "bookmark123",
		URL:        "https://example.com",
		Title:      "Example Site",
		UserID:     "user123", // Owned by different user
		IsArchived: false,
	}

	mockService.On("GetBookmark", "bookmark123").Return(bookmarkOwnedByDifferentUser, nil)

	err := handler.ArchiveBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpErr.Code)
	assert.Equal(t, "Access denied", httpErr.Message)

	mockService.AssertExpectations(t)
}

func TestBookmarkTransportMarshaling(t *testing.T) {
	bookmark := transport.BookmarkTransport{
		ID:         "bookmark123",
		URL:        "https://example.com",
		Title:      "Example",
		UserID:     "user123",
		CreatedAt:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		IsArchived: false,
	}

	jsonBytes, err := json.Marshal(bookmark)
	assert.NoError(t, err)

	var unmarshaled transport.BookmarkTransport
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, bookmark.ID, unmarshaled.ID)
	assert.Equal(t, bookmark.URL, unmarshaled.URL)
	assert.Equal(t, bookmark.Title, unmarshaled.Title)
	assert.Equal(t, bookmark.UserID, unmarshaled.UserID)
	assert.Equal(t, bookmark.IsArchived, unmarshaled.IsArchived)
}

// Benchmark tests
func BenchmarkBookmarkHandler_CreateBookmark(b *testing.B) {
	e := echo.New()
	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmark := model.Bookmark{
		ID:         "benchmark-id",
		URL:        "https://example.com",
		UserID:     "user123",
		IsArchived: false,
	}

	mockService.On("CreateBookmark", mock.Anything).Return(expectedBookmark, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bookmarkJSON := `{"url":"https://example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/bookmarks", strings.NewReader(bookmarkJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

		_ = handler.CreateBookmark(c)
	}
}

func BenchmarkBookmarkHandler_GetBookmark(b *testing.B) {
	e := echo.New()
	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	expectedBookmark := model.Bookmark{
		ID:         "bookmark123",
		URL:        "https://example.com",
		UserID:     "user123",
		IsArchived: false,
	}

	mockService.On("GetBookmark", "bookmark123").Return(expectedBookmark, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/bookmarks/bookmark123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bookmark123")
		c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

		_ = handler.GetBookmark(c)
	}
}

func TestBookmarkHandler_DeleteBookmark_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/bookmarks/bookmark123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	existingBookmark := model.Bookmark{
		ID:         "bookmark123",
		URL:        "https://example.com",
		Title:      "Example Site",
		UserID:     "user123",
		IsArchived: false,
	}

	mockService.On("GetBookmark", "bookmark123").Return(existingBookmark, nil)
	mockService.On("DeleteBookmark", "bookmark123").Return(nil)

	err := handler.DeleteBookmark(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "", rec.Body.String())

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_DeleteBookmark_MissingID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/bookmarks/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.DeleteBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "ID is required", httpErr.Message)
}

func TestBookmarkHandler_DeleteBookmark_EmptyID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/bookmarks/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("")

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.DeleteBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "ID is required", httpErr.Message)
}

func TestBookmarkHandler_DeleteBookmark_MissingAuthentication(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/bookmarks/bookmark123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Do not set user in context to simulate missing authentication

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.DeleteBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Equal(t, "User not authenticated", httpErr.Message)
}

func TestBookmarkHandler_DeleteBookmark_Forbidden_DifferentUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/bookmarks/bookmark123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// User user456 trying to delete user123's bookmark
	c.Set("user", &JWTClaims{UserID: "user456", Email: "other@example.com", Name: "Other User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	bookmarkOwnedByDifferentUser := model.Bookmark{
		ID:         "bookmark123",
		URL:        "https://example.com",
		Title:      "Example Site",
		UserID:     "user123", // Owned by different user
		IsArchived: false,
	}

	mockService.On("GetBookmark", "bookmark123").Return(bookmarkOwnedByDifferentUser, nil)

	err := handler.DeleteBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpErr.Code)
	assert.Equal(t, "Access denied", httpErr.Message)

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_DeleteBookmark_BookmarkNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/bookmarks/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	mockService.On("GetBookmark", "nonexistent").Return(model.Bookmark{}, errors.New("bookmark not found"))

	err := handler.DeleteBookmark(c)

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_DeleteBookmark_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/bookmarks/bookmark123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	existingBookmark := model.Bookmark{
		ID:         "bookmark123",
		URL:        "https://example.com",
		Title:      "Example Site",
		UserID:     "user123",
		IsArchived: false,
	}

	mockService.On("GetBookmark", "bookmark123").Return(existingBookmark, nil)
	mockService.On("DeleteBookmark", "bookmark123").Return(errors.New("database error"))

	err := handler.DeleteBookmark(c)

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_DeleteBookmark_ArchivedBookmark(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/bookmarks/bookmark123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	// Simulate Echo JWT middleware setting token in context
	c.Set("user", &JWTClaims{UserID: "user123", Email: "test@example.com", Name: "Test User"})

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	archivedBookmark := model.Bookmark{
		ID:         "bookmark123",
		URL:        "https://example.com",
		Title:      "Example Site",
		UserID:     "user123",
		IsArchived: true, // Archived bookmark
	}

	mockService.On("GetBookmark", "bookmark123").Return(archivedBookmark, nil)
	mockService.On("DeleteBookmark", "bookmark123").Return(nil)

	err := handler.DeleteBookmark(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "", rec.Body.String())

	mockService.AssertExpectations(t)
}

func TestBookmarkHandler_DeleteBookmark_NullUserContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/bookmarks/bookmark123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bookmark123")

	c.Set("user", nil)

	mockService := new(MockBookmarkService)
	handler := NewBookmarkHandler(mockService)

	err := handler.DeleteBookmark(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Equal(t, "User not authenticated", httpErr.Message)
}
