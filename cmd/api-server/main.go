package main

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tsongpon/athena/internal/handler"
	"github.com/tsongpon/athena/internal/repository"
	"github.com/tsongpon/athena/internal/service"
)

func main() {
	bookmarkRepo := repository.NewBookmarkInMemRepository()
	bookmarkService := service.NewBookmarkService(bookmarkRepo)

	userRepo := repository.NewUserInMemRepository()
	userService := service.NewUserService(userRepo)

	bookmarkHandler := handler.NewBookmarkHandler(bookmarkService)
	authHandler := handler.NewAuthHandler(userService)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = handler.DefaultJWTSecret
	}

	// JWT middleware config
	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handler.JWTClaims)
		},
		SigningKey: []byte(jwtSecret),
		// Important: Echo JWT v4 uses "user" as the default context key
		// The middleware will extract claims from the token and store them in context
		ContextKey: "user",
	}

	// Routes
	e.GET("/ping", bookmarkHandler.Ping)

	// Authentication routes
	e.POST("/users", authHandler.CreateUser)
	e.POST("/login", authHandler.Login)

	// Bookmark routes (all protected with JWT)
	e.POST("/bookmarks", bookmarkHandler.CreateBookmark, echojwt.WithConfig(jwtConfig))
	e.GET("/bookmarks/:id", bookmarkHandler.GetBookmark, echojwt.WithConfig(jwtConfig))
	e.GET("/bookmarks", bookmarkHandler.GetBookmarks, echojwt.WithConfig(jwtConfig))
	e.POST("/bookmarks/:id/archive", bookmarkHandler.ArchiveBookmark, echojwt.WithConfig(jwtConfig))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
