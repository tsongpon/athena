package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tsongpon/athena/internal/model"
	"github.com/tsongpon/athena/internal/transport"
)

type BookmarkHandler struct {
	bookmarkService BookmarkService
}

func NewBookmarkHandler(service BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{
		bookmarkService: service,
	}
}

func (h *BookmarkHandler) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (h *BookmarkHandler) CreateBookmark(c echo.Context) error {
	bt := &transport.BookmarkTransport{}
	if err := c.Bind(bt); err != nil {
		return err
	}
	if bt.URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "URL is required")
	}

	// Get authenticated user ID from JWT token
	authenticatedUser, err := getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	b := model.Bookmark{
		URL:    bt.URL,
		UserID: authenticatedUser.UserID, // Use authenticated user's ID from JWT
	}
	createdBookmark, err := h.bookmarkService.CreateBookmark(b)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	responseTransport := transport.BookmarkTransport{
		ID:         createdBookmark.ID,
		URL:        createdBookmark.URL,
		Title:      createdBookmark.Title,
		UserID:     createdBookmark.UserID,
		IsArchived: createdBookmark.IsArchived,
	}
	return c.JSON(http.StatusCreated, responseTransport)
}

func (h *BookmarkHandler) GetBookmark(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ID is required")
	}

	// Get authenticated user ID from JWT token
	authenticatedUser, err := getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	bookmarks, err := h.bookmarkService.GetBookmark(id)
	if err != nil {
		return err
	}

	// Authorization check: ensure user can only access their own bookmarks
	if bookmarks.UserID != authenticatedUser.UserID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	t := transport.BookmarkTransport{
		ID:         bookmarks.ID,
		URL:        bookmarks.URL,
		Title:      bookmarks.Title,
		UserID:     bookmarks.UserID,
		IsArchived: bookmarks.IsArchived,
	}
	return c.JSON(http.StatusOK, t)
}

func (h *BookmarkHandler) GetBookmarks(c echo.Context) error {
	// Get authenticated user ID from JWT token
	authenticatedUser, err := getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	// Use authenticated user's ID instead of query parameter
	userID := authenticatedUser.UserID
	archivedParam := c.QueryParam("archived")
	archived, err := strconv.ParseBool(archivedParam)
	if err != nil {
		archived = false
	}
	bookmarks, err := h.bookmarkService.GetAllBookmarks(userID, archived)
	if err != nil {
		return err
	}
	ts := make([]transport.BookmarkTransport, len(bookmarks))
	for i, b := range bookmarks {
		ts[i] = transport.BookmarkTransport{
			ID:         b.ID,
			URL:        b.URL,
			Title:      b.Title,
			UserID:     b.UserID,
			CreatedAt:  b.CreatedAt,
			IsArchived: b.IsArchived,
		}
	}
	return c.JSON(http.StatusOK, ts)
}

func (h *BookmarkHandler) ArchiveBookmark(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ID is required")
	}

	// Get authenticated user ID from JWT token
	authenticatedUser, err := getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	// Get the bookmark first to verify ownership
	bookmark, err := h.bookmarkService.GetBookmark(id)
	if err != nil {
		return err
	}

	// Authorization check: ensure user can only archive their own bookmarks
	if bookmark.UserID != authenticatedUser.UserID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	if _, err := h.bookmarkService.ArchiveBookmark(id); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
