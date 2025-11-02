package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/tsongpon/athena/internal/logger"
	"go.uber.org/zap"
)

// PostgresConfig holds PostgreSQL connection configuration
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresConnection creates a new PostgreSQL database connection
func NewPostgresConnection(config PostgresConfig) (*sql.DB, error) {
	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	logger.Info("Connecting to PostgreSQL database",
		zap.String("host", config.Host),
		zap.String("port", config.Port),
		zap.String("dbname", config.DBName))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Failed to open database connection", zap.Error(err))
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database")
	return db, nil
}

// RunMigrations executes database migrations
func RunMigrations(db *sql.DB) error {
	logger.Info("Running database migrations")

	migration := `
		-- Create bookmarks table
		CREATE TABLE IF NOT EXISTS bookmarks (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			url TEXT NOT NULL,
			title TEXT NOT NULL,
			is_archived BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);

		-- Create index on user_id for efficient queries
		CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);

		-- Create index on user_id and is_archived for efficient filtering
		CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id_archived ON bookmarks(user_id, is_archived);

		-- Create index on created_at for sorting
		CREATE INDEX IF NOT EXISTS idx_bookmarks_created_at ON bookmarks(created_at DESC);
	`

	_, err := db.Exec(migration)
	if err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}
