package repository

import (
	"fmt"
	"time"

	"github.com/tsongpon/athena/internal/model"
)

// BookmarkInMemRepository implements BookmarkRepository interface using an in-memory map
type BookmarkInMemRepository struct {
	bookmarks map[string]model.Bookmark
	// mutex     sync.RWMutex
}

// NewBookmarkInMemRepository creates a new instance of BookmarkInMemRepository
func NewBookmarkInMemRepository() BookmarkInMemRepository {
	return BookmarkInMemRepository{
		bookmarks: make(map[string]model.Bookmark),
		// mutex:     sync.RWMutex{},
	}
}

// CreateBookmark creates a new bookmark in the repository
func (r BookmarkInMemRepository) CreateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
	// r.mutex.Lock()
	// defer r.mutex.Unlock()

	// Check if bookmark with this ID already exists
	if _, exists := r.bookmarks[bookmark.ID]; exists {
		return model.Bookmark{}, fmt.Errorf("bookmark with ID %s already exists", bookmark.ID)
	}

	// Set creation time if not provided
	if bookmark.CreatedAt.IsZero() {
		bookmark.CreatedAt = time.Now()
	}

	// Store the bookmark
	r.bookmarks[bookmark.ID] = bookmark

	return bookmark, nil
}

// GetBookmark retrieves a bookmark by its ID
func (r BookmarkInMemRepository) GetBookmark(id string) (model.Bookmark, error) {
	// r.mutex.RLock()
	// defer r.mutex.RUnlock()

	// Check if bookmark exists
	bookmark, exists := r.bookmarks[id]
	if !exists {
		return model.Bookmark{}, fmt.Errorf("bookmark with ID %s not found", id)
	}

	return bookmark, nil
}

// ListBookmarks retrieves all bookmarks for a specific user
func (r BookmarkInMemRepository) ListBookmarks(userID string) ([]model.Bookmark, error) {
	// r.mutex.RLock()
	// defer r.mutex.RUnlock()

	var userBookmarks []model.Bookmark

	for _, bookmark := range r.bookmarks {
		if bookmark.UserID == userID {
			userBookmarks = append(userBookmarks, bookmark)
		}
	}

	return userBookmarks, nil
}

// UpdateBookmark updates an existing bookmark in the repository
func (r BookmarkInMemRepository) UpdateBookmark(bookmark model.Bookmark) (model.Bookmark, error) {
	// r.mutex.Lock()
	// defer r.mutex.Unlock()

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
func (r BookmarkInMemRepository) DeleteBookmark(id string) error {
	// r.mutex.Lock()
	// defer r.mutex.Unlock()

	// Check if bookmark exists
	if _, exists := r.bookmarks[id]; !exists {
		return fmt.Errorf("bookmark with ID %s not found", id)
	}

	// Delete the bookmark
	delete(r.bookmarks, id)

	return nil
}
