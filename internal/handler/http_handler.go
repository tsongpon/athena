package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tsongpon/athena/internal/model"
	"github.com/tsongpon/athena/internal/transport"
)

type HTTPHandler struct {
	bookmarkService BookmarkService
}

func NewHTTPHandler(service BookmarkService) *HTTPHandler {
	return &HTTPHandler{
		bookmarkService: service,
	}
}

func (h *HTTPHandler) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (h *HTTPHandler) CreateBookmark(c echo.Context) error {
	bt := &transport.BookmarkTransport{}
	if err := c.Bind(bt); err != nil {
		return err
	}
	if bt.URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "URL is required")
	}
	if bt.UserID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is required")
	}
	b := model.Bookmark{
		URL:    bt.URL,
		UserID: bt.UserID,
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

func (h *HTTPHandler) GetBookmark(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ID is required")
	}
	bookmarks, err := h.bookmarkService.GetBookmark(id)
	if err != nil {
		return err
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

func (h *HTTPHandler) GetBookmarks(c echo.Context) error {
	userID := c.QueryParam("user_id")
	archivedParam := c.QueryParam("archived")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is required")
	}
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

func (h *HTTPHandler) ArchiveBookmark(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ID is required")
	}
	if _, err := h.bookmarkService.ArchiveBookmark(id); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
