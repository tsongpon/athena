package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tsongpon/athena/internal/logger"
	"github.com/tsongpon/athena/internal/model"
	"go.uber.org/zap"
)

// BookmarkPostgresRepository implements BookmarkRepository interface using PostgreSQL
type BookmarkPostgresRepository struct {
	db *sql.DB
}

// NewBookmarkPostgresRepository creates a new instance of BookmarkPostgresRepository
func NewBookmarkPostgresRepository(db *sql.DB) *BookmarkPostgresRepository {
	return &BookmarkPostgresRepository{
		db: db,
	}
}

// CreateBookmark creates a new bookmark in PostgreSQL
func (r *BookmarkPostgresRepository) CreateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
	// Generate ID if not provided
	if bookmark.ID == "" {
		bookmark.ID = uuid.New().String()
	}

	// Set creation time if not provided
	if bookmark.CreatedAt.IsZero() {
		bookmark.CreatedAt = time.Now()
	}

	// Set updated time
	bookmark.UpdatedAt = time.Now()

	query := `
		INSERT INTO bookmarks (id, user_id, url, title, is_archived, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query,
		bookmark.ID,
		bookmark.UserID,
		bookmark.URL,
		bookmark.Title,
		bookmark.IsArchived,
		bookmark.CreatedAt,
		bookmark.UpdatedAt,
	)

	if err != nil {
		logger.Error("Failed to create bookmark in database",
			zap.String("bookmark_id", bookmark.ID),
			zap.String("user_id", bookmark.UserID),
			zap.Error(err))
		return model.Bookmark{}, fmt.Errorf("failed to create bookmark: %w", err)
	}

	logger.Debug("Created bookmark in database", zap.String("id", bookmark.ID))
	return bookmark, nil
}

// GetBookmark retrieves a bookmark by its ID from PostgreSQL
func (r *BookmarkPostgresRepository) GetBookmark(id string) (model.Bookmark, error) {
	logger.Debug("Getting bookmark from database", zap.String("id", id))

	query := `
		SELECT id, user_id, url, title, is_archived, created_at, updated_at
		FROM bookmarks
		WHERE id = $1
	`

	var bookmark model.Bookmark
	err := r.db.QueryRow(query, id).Scan(
		&bookmark.ID,
		&bookmark.UserID,
		&bookmark.URL,
		&bookmark.Title,
		&bookmark.IsArchived,
		&bookmark.CreatedAt,
		&bookmark.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return model.Bookmark{}, fmt.Errorf("bookmark with ID %s not found", id)
	}

	if err != nil {
		logger.Error("Failed to get bookmark from database",
			zap.String("id", id),
			zap.Error(err))
		return model.Bookmark{}, fmt.Errorf("failed to get bookmark: %w", err)
	}

	return bookmark, nil
}

// ListBookmarks retrieves all bookmarks based on the query parameters from PostgreSQL
// Returns bookmarks ordered by created date descending (newest first)
// Supports pagination when Page and PageSize are greater than 0
func (r *BookmarkPostgresRepository) ListBookmarks(query model.BookmarkQuery) ([]model.Bookmark, error) {
	// Build the SQL query with optional pagination
	var sqlQuery string
	var args []interface{}

	if query.Page > 0 && query.PageSize > 0 {
		// With pagination
		offset := (query.Page - 1) * query.PageSize
		sqlQuery = `
			SELECT id, user_id, url, title, is_archived, created_at, updated_at
			FROM bookmarks
			WHERE user_id = $1 AND is_archived = $2
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`
		args = []interface{}{query.UserID, query.Archived, query.PageSize, offset}
	} else {
		// Without pagination
		sqlQuery = `
			SELECT id, user_id, url, title, is_archived, created_at, updated_at
			FROM bookmarks
			WHERE user_id = $1 AND is_archived = $2
			ORDER BY created_at DESC
		`
		args = []interface{}{query.UserID, query.Archived}
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		logger.Error("Failed to list bookmarks from database",
			zap.String("user_id", query.UserID),
			zap.Bool("archived", query.Archived),
			zap.Int("page", query.Page),
			zap.Int("page_size", query.PageSize),
			zap.Error(err))
		return nil, fmt.Errorf("failed to list bookmarks: %w", err)
	}
	defer rows.Close()

	var bookmarks []model.Bookmark
	for rows.Next() {
		var bookmark model.Bookmark
		err := rows.Scan(
			&bookmark.ID,
			&bookmark.UserID,
			&bookmark.URL,
			&bookmark.Title,
			&bookmark.IsArchived,
			&bookmark.CreatedAt,
			&bookmark.UpdatedAt,
		)
		if err != nil {
			logger.Error("Failed to scan bookmark row",
				zap.Error(err))
			return nil, fmt.Errorf("failed to scan bookmark: %w", err)
		}
		bookmarks = append(bookmarks, bookmark)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Error iterating bookmark rows",
			zap.Error(err))
		return nil, fmt.Errorf("error iterating bookmarks: %w", err)
	}

	logger.Debug("Listed bookmarks from database",
		zap.String("user_id", query.UserID),
		zap.Int("count", len(bookmarks)),
		zap.Int("page", query.Page),
		zap.Int("page_size", query.PageSize))

	return bookmarks, nil
}

// CountBookmarks returns the total count of bookmarks matching the query
func (r *BookmarkPostgresRepository) CountBookmarks(query model.BookmarkQuery) (int, error) {
	sqlQuery := `
		SELECT COUNT(*)
		FROM bookmarks
		WHERE user_id = $1 AND is_archived = $2
	`

	var count int
	err := r.db.QueryRow(sqlQuery, query.UserID, query.Archived).Scan(&count)
	if err != nil {
		logger.Error("Failed to count bookmarks from database",
			zap.String("user_id", query.UserID),
			zap.Bool("archived", query.Archived),
			zap.Error(err))
		return 0, fmt.Errorf("failed to count bookmarks: %w", err)
	}

	logger.Debug("Counted bookmarks from database",
		zap.String("user_id", query.UserID),
		zap.Int("count", count))

	return count, nil
}

// UpdateBookmark updates an existing bookmark in PostgreSQL
func (r *BookmarkPostgresRepository) UpdateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
	// Set updated time
	bookmark.UpdatedAt = time.Now()

	query := `
		UPDATE bookmarks
		SET user_id = $2, url = $3, title = $4, is_archived = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.Exec(query,
		bookmark.ID,
		bookmark.UserID,
		bookmark.URL,
		bookmark.Title,
		bookmark.IsArchived,
		bookmark.UpdatedAt,
	)

	if err != nil {
		logger.Error("Failed to update bookmark in database",
			zap.String("bookmark_id", bookmark.ID),
			zap.Error(err))
		return model.Bookmark{}, fmt.Errorf("failed to update bookmark: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Failed to get rows affected",
			zap.String("bookmark_id", bookmark.ID),
			zap.Error(err))
		return model.Bookmark{}, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.Bookmark{}, fmt.Errorf("bookmark with ID %s not found", bookmark.ID)
	}

	logger.Debug("Updated bookmark in database", zap.String("id", bookmark.ID))

	// Retrieve the updated bookmark to get the created_at time
	return r.GetBookmark(bookmark.ID)
}

// DeleteBookmark removes a bookmark from PostgreSQL
func (r *BookmarkPostgresRepository) DeleteBookmark(id string) error {
	query := `DELETE FROM bookmarks WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		logger.Error("Failed to delete bookmark from database",
			zap.String("id", id),
			zap.Error(err))
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Failed to get rows affected",
			zap.String("id", id),
			zap.Error(err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("bookmark with ID %s not found", id)
	}

	logger.Debug("Deleted bookmark from database", zap.String("id", id))
	return nil
}
