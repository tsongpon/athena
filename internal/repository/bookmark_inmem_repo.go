package repository

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tsongpon/athena/internal/logger"
	"github.com/tsongpon/athena/internal/model"
	"go.uber.org/zap"
)

// BookmarkInMemRepository implements BookmarkRepository interface using an in-memory map
type BookmarkInMemRepository struct {
	bookmarks map[string]model.Bookmark
	mutex     sync.RWMutex
}

// NewBookmarkInMemRepository creates a new instance of BookmarkInMemRepository
func NewBookmarkInMemRepository() *BookmarkInMemRepository {
	return &BookmarkInMemRepository{
		bookmarks: make(map[string]model.Bookmark),
		mutex:     sync.RWMutex{},
	}
}

// CreateBookmark creates a new bookmark in the repository
func (r *BookmarkInMemRepository) CreateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	bookmark.ID = uuid.New().String()
	// Set creation time if not provided
	if bookmark.CreatedAt.IsZero() {
		bookmark.CreatedAt = time.Now()
	}

	// Store the bookmark
	r.bookmarks[bookmark.ID] = bookmark

	return bookmark, nil
}

// GetBookmark retrieves a bookmark by its ID
func (r *BookmarkInMemRepository) GetBookmark(id string) (model.Bookmark, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Check if bookmark exists
	logger.Debug("Getting bookmark", zap.String("id", id))
	bookmark, exists := r.bookmarks[id]
	if !exists {
		return model.Bookmark{}, fmt.Errorf("bookmark with ID %s not found", id)
	}

	return bookmark, nil
}

// ListBookmarks retrieves all bookmarks based on the query parameters
// Returns bookmarks ordered by created date descending (newest first)
// Supports pagination when Page and PageSize are greater than 0
func (r *BookmarkInMemRepository) ListBookmarks(query model.BookmarkQuery) ([]model.Bookmark, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var userBookmarks []model.Bookmark

	for _, bookmark := range r.bookmarks {
		if bookmark.UserID == query.UserID && bookmark.IsArchived == query.Archived {
			userBookmarks = append(userBookmarks, bookmark)
		}
	}

	// Sort by created date descending (newest first)
	sort.Slice(userBookmarks, func(i, j int) bool {
		return userBookmarks[i].CreatedAt.After(userBookmarks[j].CreatedAt)
	})

	// Apply pagination if specified
	if query.Page > 0 && query.PageSize > 0 {
		start := (query.Page - 1) * query.PageSize
		end := start + query.PageSize

		// Handle out of bounds
		if start >= len(userBookmarks) {
			return []model.Bookmark{}, nil
		}
		if end > len(userBookmarks) {
			end = len(userBookmarks)
		}

		return userBookmarks[start:end], nil
	}

	return userBookmarks, nil
}

// CountBookmarks returns the total count of bookmarks matching the query
func (r *BookmarkInMemRepository) CountBookmarks(query model.BookmarkQuery) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, bookmark := range r.bookmarks {
		if bookmark.UserID == query.UserID && bookmark.IsArchived == query.Archived {
			count++
		}
	}

	return count, nil
}

// UpdateBookmark updates an existing bookmark in the repository
func (r *BookmarkInMemRepository) UpdateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if bookmark exists
	existing, exists := r.bookmarks[bookmark.ID]
	if !exists {
		return model.Bookmark{}, fmt.Errorf("bookmark with ID %s not found", bookmark.ID)
	}

	// Preserve creation time from existing bookmark
	bookmark.CreatedAt = existing.CreatedAt

	// Update the bookmark
	r.bookmarks[bookmark.ID] = bookmark

	return bookmark, nil
}

// DeleteBookmark removes a bookmark from the repository
func (r *BookmarkInMemRepository) DeleteBookmark(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if bookmark exists
	if _, exists := r.bookmarks[id]; !exists {
		return fmt.Errorf("bookmark with ID %s not found", id)
	}

	// Delete the bookmark
	delete(r.bookmarks, id)

	return nil
}
