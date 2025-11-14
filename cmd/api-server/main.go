package main

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tsongpon/athena/internal/handler"
	"github.com/tsongpon/athena/internal/logger"
	"github.com/tsongpon/athena/internal/repository"
	"github.com/tsongpon/athena/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.Initialize(); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting Athena API server")

	// Determine storage type from environment variable
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "memory" // Default to in-memory
	}

	var bookmarkRepo service.BookmarkRepository
	var userRepo service.UserRepository

	switch storageType {
	case "firestore":
		// Firestore configuration from environment variables
		ctx := context.Background()
		projectID := os.Getenv("GCP_PROJECT_ID")
		databaseID := os.Getenv("GCP_FIRESTORE_DATABASE_ID")
		if projectID == "" {
			logger.Fatal("GCP_PROJECT_ID environment variable is required for Firestore")
		}

		client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
		if err != nil {
			logger.Fatal("Failed to create Firestore client", zap.Error(err))
		}
		defer client.Close()

		bookmarkRepo = repository.NewBookmarkFirestoreRepository(ctx, client)
		userRepo = repository.NewUserFirestoreRepository(ctx, client)
		logger.Info("Using Firestore storage for bookmarks and users", zap.String("project_id", projectID))

	default:
		bookmarkRepo = repository.NewBookmarkInMemRepository()
		userRepo = repository.NewUserInMemRepository()
		logger.Info("Using in-memory storage for bookmarks and users")
	}

	webRepo := repository.NewWebRepository()
	bookmarkService := service.NewBookmarkService(bookmarkRepo, userRepo, webRepo)
	userService := service.NewUserService(userRepo)

	bookmarkHandler := handler.NewBookmarkHandler(bookmarkService)
	authHandler := handler.NewAuthHandler(userService)

	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
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
	e.DELETE("/bookmarks/:id", bookmarkHandler.DeleteBookmark, echojwt.WithConfig(jwtConfig))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	logger.Info("Server starting on port", zap.String("port", port))
	if err := e.Start(":" + port); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
