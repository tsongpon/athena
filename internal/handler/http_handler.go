package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tsongpon/athena/internal/service"
)

type HTTPHandler struct {
	bookmarkService service.BookmarkService
}

func NewHTTPHandler(service service.BookmarkService) HTTPHandler {
	return HTTPHandler{
		bookmarkService: service,
	}
}

func (h HTTPHandler) Ping(c echo.Context) error {
	h.bookmarkService.CreateBookmark("123", "site")
	return c.String(http.StatusOK, "pong")
}
