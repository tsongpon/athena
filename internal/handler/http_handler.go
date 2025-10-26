package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tsongpon/athena/internal/model"
	"github.com/tsongpon/athena/internal/transport"
)

type HTTPHandler struct {
	bookmarkService BookmarkService
}

func NewHTTPHandler(service BookmarkService) HTTPHandler {
	return HTTPHandler{
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
	b := model.Bookmark{
		URL:        bt.URL,
		UserID:     bt.UserID,
		IsArchived: false,
	}
	if _, err := h.bookmarkService.CreateBookmark(b); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
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
		ID:     bookmarks.ID,
		URL:    bookmarks.URL,
		Title:  bookmarks.Title,
		UserID: bookmarks.UserID,
	}
	return c.JSON(http.StatusOK, t)
}

func (h *HTTPHandler) GetBookmarks(c echo.Context) error {
	userID := c.QueryParam("userid")
	bookmarks, err := h.bookmarkService.GetAllBookmarks(userID)
	if err != nil {
		return err
	}
	ts := make([]transport.BookmarkTransport, len(bookmarks))
	for i, b := range bookmarks {
		ts[i] = transport.BookmarkTransport{
			ID:        b.ID,
			URL:       b.URL,
			Title:     b.Title,
			UserID:    b.UserID,
			CreatedAt: b.CreatedAt,
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
