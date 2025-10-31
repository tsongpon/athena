package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tsongpon/athena/internal/handler"
	"github.com/tsongpon/athena/internal/repository"
	"github.com/tsongpon/athena/internal/service"
)

func main() {
	bookmarkRepo := repository.NewBookmarkInMemRepository()
	bookmarkService := service.NewBookmarkService(bookmarkRepo)
	httpHandler := handler.NewHTTPHandler(bookmarkService)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/ping", httpHandler.Ping)
	e.POST("/bookmarks", httpHandler.CreateBookmark)
	e.GET("/bookmarks/:id", httpHandler.GetBookmark)
	e.GET("/bookmarks", httpHandler.GetBookmarks)
	e.POST("/bookmarks/:id/archive", httpHandler.ArchiveBookmark)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
