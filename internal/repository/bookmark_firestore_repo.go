package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/tsongpon/athena/internal/logger"
	"github.com/tsongpon/athena/internal/model"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
)

const bookmarksCollection = "bookmarks"

// BookmarkFirestoreRepository implements BookmarkRepository interface using GCP Firestore
type BookmarkFirestoreRepository struct {
	client *firestore.Client
	ctx    context.Context
}

// NewBookmarkFirestoreRepository creates a new instance of BookmarkFirestoreRepository
func NewBookmarkFirestoreRepository(ctx context.Context, client *firestore.Client) *BookmarkFirestoreRepository {
	return &BookmarkFirestoreRepository{
		client: client,
		ctx:    ctx,
	}
}

// firestoreBookmark is the structure used to store/retrieve bookmarks in Firestore
type firestoreBookmark struct {
	ID           string    `firestore:"id"`
	UserID       string    `firestore:"user_id"`
	URL          string    `firestore:"url"`
	Title        string    `firestore:"title"`
	IsArchived   bool      `firestore:"is_archived"`
	MainImageURL string    `firestore:"main_image_url"`
	CreatedAt    time.Time `firestore:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at"`
}

// toFirestoreBookmark converts model.Bookmark to firestoreBookmark
func toFirestoreBookmark(bookmark model.Bookmark) firestoreBookmark {
	return firestoreBookmark{
		ID:           bookmark.ID,
		UserID:       bookmark.UserID,
		URL:          bookmark.URL,
		Title:        bookmark.Title,
		IsArchived:   bookmark.IsArchived,
		MainImageURL: bookmark.MainImageURL,
		CreatedAt:    bookmark.CreatedAt,
		UpdatedAt:    bookmark.UpdatedAt,
	}
}

// toModelBookmark converts firestoreBookmark to model.Bookmark
func toModelBookmark(fsBookmark firestoreBookmark) model.Bookmark {
	return model.Bookmark{
		ID:           fsBookmark.ID,
		UserID:       fsBookmark.UserID,
		URL:          fsBookmark.URL,
		Title:        fsBookmark.Title,
		IsArchived:   fsBookmark.IsArchived,
		MainImageURL: fsBookmark.MainImageURL,
		CreatedAt:    fsBookmark.CreatedAt,
		UpdatedAt:    fsBookmark.UpdatedAt,
	}
}

// CreateBookmark creates a new bookmark in Firestore
func (r *BookmarkFirestoreRepository) CreateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
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

	// Convert to Firestore structure
	fsBookmark := toFirestoreBookmark(bookmark)

	// Store in Firestore using bookmark ID as document ID
	_, err := r.client.Collection(bookmarksCollection).Doc(bookmark.ID).Set(r.ctx, fsBookmark)
	if err != nil {
		logger.Error("Failed to create bookmark in Firestore",
			zap.String("bookmark_id", bookmark.ID),
			zap.String("user_id", bookmark.UserID),
			zap.Error(err))
		return model.Bookmark{}, fmt.Errorf("failed to create bookmark: %w", err)
	}

	logger.Debug("Created bookmark in Firestore", zap.String("id", bookmark.ID))
	return bookmark, nil
}

// GetBookmark retrieves a bookmark by its ID from Firestore
func (r *BookmarkFirestoreRepository) GetBookmark(id string) (model.Bookmark, error) {
	logger.Debug("Getting bookmark from Firestore", zap.String("id", id))

	docSnap, err := r.client.Collection(bookmarksCollection).Doc(id).Get(r.ctx)
	if err != nil {
		logger.Error("Failed to get bookmark from Firestore",
			zap.String("id", id),
			zap.Error(err))
		return model.Bookmark{}, fmt.Errorf("bookmark with ID %s not found: %w", id, err)
	}

	var fsBookmark firestoreBookmark
	if err := docSnap.DataTo(&fsBookmark); err != nil {
		logger.Error("Failed to parse bookmark data from Firestore",
			zap.String("id", id),
			zap.Error(err))
		return model.Bookmark{}, fmt.Errorf("failed to parse bookmark data: %w", err)
	}

	return toModelBookmark(fsBookmark), nil
}

// ListBookmarks retrieves all bookmarks based on the query parameters from Firestore
// Returns bookmarks ordered by created date descending (newest first)
// Supports pagination when Page and PageSize are greater than 0
func (r *BookmarkFirestoreRepository) ListBookmarks(query model.BookmarkQuery) ([]model.Bookmark, error) {
	// Build Firestore query
	firestoreQuery := r.client.Collection(bookmarksCollection).
		Where("user_id", "==", query.UserID).
		Where("is_archived", "==", query.Archived).
		OrderBy("created_at", firestore.Desc)

	// Apply pagination if specified
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		firestoreQuery = firestoreQuery.Limit(query.PageSize).Offset(offset)
	}

	// Execute query
	iter := firestoreQuery.Documents(r.ctx)
	defer iter.Stop()

	var bookmarks []model.Bookmark
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error("Failed to list bookmarks from Firestore",
				zap.String("user_id", query.UserID),
				zap.Bool("archived", query.Archived),
				zap.Int("page", query.Page),
				zap.Int("page_size", query.PageSize),
				zap.Error(err))
			return nil, fmt.Errorf("failed to list bookmarks: %w", err)
		}

		var fsBookmark firestoreBookmark
		if err := doc.DataTo(&fsBookmark); err != nil {
			logger.Error("Failed to parse bookmark data from Firestore",
				zap.Error(err))
			return nil, fmt.Errorf("failed to parse bookmark data: %w", err)
		}

		bookmarks = append(bookmarks, toModelBookmark(fsBookmark))
	}

	logger.Debug("Listed bookmarks from Firestore",
		zap.String("user_id", query.UserID),
		zap.Int("count", len(bookmarks)),
		zap.Int("page", query.Page),
		zap.Int("page_size", query.PageSize))

	return bookmarks, nil
}

// CountBookmarks returns the total count of bookmarks matching the query
func (r *BookmarkFirestoreRepository) CountBookmarks(query model.BookmarkQuery) (int, error) {
	// Build Firestore query and iterate to count
	firestoreQuery := r.client.Collection(bookmarksCollection).
		Where("user_id", "==", query.UserID).
		Where("is_archived", "==", query.Archived)

	iter := firestoreQuery.Documents(r.ctx)
	defer iter.Stop()

	count := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error("Failed to count bookmarks from Firestore",
				zap.String("user_id", query.UserID),
				zap.Bool("archived", query.Archived),
				zap.Error(err))
			return 0, fmt.Errorf("failed to count bookmarks: %w", err)
		}
		count++
	}

	logger.Debug("Counted bookmarks from Firestore",
		zap.String("user_id", query.UserID),
		zap.Int("count", count))

	return count, nil
}

// UpdateBookmark updates an existing bookmark in Firestore
func (r *BookmarkFirestoreRepository) UpdateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
	// First, get the existing bookmark to preserve CreatedAt
	existing, err := r.GetBookmark(bookmark.ID)
	if err != nil {
		return model.Bookmark{}, err
	}

	// Preserve creation time
	bookmark.CreatedAt = existing.CreatedAt

	// Set updated time
	bookmark.UpdatedAt = time.Now()

	// Convert to Firestore structure
	fsBookmark := toFirestoreBookmark(bookmark)

	// Update in Firestore
	_, err = r.client.Collection(bookmarksCollection).Doc(bookmark.ID).Set(r.ctx, fsBookmark)
	if err != nil {
		logger.Error("Failed to update bookmark in Firestore",
			zap.String("bookmark_id", bookmark.ID),
			zap.Error(err))
		return model.Bookmark{}, fmt.Errorf("failed to update bookmark: %w", err)
	}

	logger.Debug("Updated bookmark in Firestore", zap.String("id", bookmark.ID))
	return bookmark, nil
}

// DeleteBookmark removes a bookmark from Firestore
func (r *BookmarkFirestoreRepository) DeleteBookmark(id string) error {
	// Check if bookmark exists before deleting
	_, err := r.GetBookmark(id)
	if err != nil {
		return err
	}

	// Delete from Firestore
	_, err = r.client.Collection(bookmarksCollection).Doc(id).Delete(r.ctx)
	if err != nil {
		logger.Error("Failed to delete bookmark from Firestore",
			zap.String("id", id),
			zap.Error(err))
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	logger.Debug("Deleted bookmark from Firestore", zap.String("id", id))
	return nil
}
