package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tsongpon/athena/internal/logger"
	"github.com/tsongpon/athena/internal/model"
	"github.com/tsongpon/athena/internal/transport"
	"go.uber.org/zap"
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
		logger.Warn("Failed to bind bookmark request", zap.Error(err))
		return err
	}
	if bt.URL == "" {
		logger.Warn("Create bookmark request missing URL")
		return echo.NewHTTPError(http.StatusBadRequest, "URL is required")
	}

	// Get authenticated user ID from JWT token
	authenticatedUser, err := getAuthenticatedUser(c)
	if err != nil {
		logger.Error("Failed to get authenticated user", zap.Error(err))
		return err
	}

	b := model.Bookmark{
		URL:    bt.URL,
		UserID: authenticatedUser.UserID, // Use authenticated user's ID from JWT
	}
	createdBookmark, err := h.bookmarkService.CreateBookmark(b)
	if err != nil {
		logger.Error("Failed to create bookmark",
			zap.String("user_id", authenticatedUser.UserID),
			zap.String("url", bt.URL),
			zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	logger.Info("Bookmark created via API",
		zap.String("bookmark_id", createdBookmark.ID),
		zap.String("user_id", createdBookmark.UserID),
		zap.String("url", createdBookmark.URL))
	responseTransport := transport.BookmarkTransport{
		ID:             createdBookmark.ID,
		URL:            createdBookmark.URL,
		Title:          createdBookmark.Title,
		UserID:         createdBookmark.UserID,
		MainImageURL:   createdBookmark.MainImageURL,
		ContentSummary: createdBookmark.ContentSummary,
		IsArchived:     createdBookmark.IsArchived,
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

	bookmark, err := h.bookmarkService.GetBookmark(id)
	if err != nil {
		return err
	}

	// Authorization check: ensure user can only access their own bookmarks
	if bookmark.UserID != authenticatedUser.UserID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	t := transport.BookmarkTransport{
		ID:             bookmark.ID,
		URL:            bookmark.URL,
		Title:          bookmark.Title,
		UserID:         bookmark.UserID,
		MainImageURL:   bookmark.MainImageURL,
		ContentSummary: bookmark.ContentSummary,
		IsArchived:     bookmark.IsArchived,
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

	// Check for pagination parameters
	pageParam := c.QueryParam("page")
	pageSizeParam := c.QueryParam("page_size")

	// If pagination parameters are provided, use paginated endpoint
	if pageParam != "" || pageSizeParam != "" {
		page, err := strconv.Atoi(pageParam)
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil || pageSize < 1 {
			pageSize = 20 // Default page size
		}

		response, err := h.bookmarkService.GetBookmarksWithPagination(userID, archived, page, pageSize)
		if err != nil {
			logger.Error("Failed to get paginated bookmarks",
				zap.String("user_id", userID),
				zap.Int("page", page),
				zap.Int("page_size", pageSize),
				zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// Convert bookmarks to transport format
		ts := make([]transport.BookmarkTransport, len(response.Bookmarks))
		for i, b := range response.Bookmarks {
			ts[i] = transport.BookmarkTransport{
				ID:             b.ID,
				URL:            b.URL,
				Title:          b.Title,
				UserID:         b.UserID,
				MainImageURL:   b.MainImageURL,
				ContentSummary: b.ContentSummary,
				CreatedAt:      b.CreatedAt,
				IsArchived:     b.IsArchived,
			}
		}

		// Return paginated response
		paginatedResponse := map[string]interface{}{
			"bookmarks":   ts,
			"total_count": response.TotalCount,
			"page":        response.Page,
			"page_size":   response.PageSize,
			"total_pages": response.TotalPages,
		}

		return c.JSON(http.StatusOK, paginatedResponse)
	}

	// No pagination - return all bookmarks
	bookmarks, err := h.bookmarkService.GetAllBookmarks(userID, archived)
	if err != nil {
		return err
	}
	ts := make([]transport.BookmarkTransport, len(bookmarks))
	for i, b := range bookmarks {
		ts[i] = transport.BookmarkTransport{
			ID:             b.ID,
			URL:            b.URL,
			Title:          b.Title,
			UserID:         b.UserID,
			MainImageURL:   b.MainImageURL,
			ContentSummary: b.ContentSummary,
			CreatedAt:      b.CreatedAt,
			IsArchived:     b.IsArchived,
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

func (h *BookmarkHandler) DeleteBookmark(c echo.Context) error {
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

	// Authorization check: ensure user can only delete their own bookmarks
	if bookmark.UserID != authenticatedUser.UserID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	if err := h.bookmarkService.DeleteBookmark(id); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
