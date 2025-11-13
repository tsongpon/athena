package service

import (
	"fmt"

	"github.com/tsongpon/athena/internal/logger"
	"github.com/tsongpon/athena/internal/model"
	"go.uber.org/zap"
)

// bookmarkService is the concrete implementation of BookmarkService interface
type BookmarkService struct {
	bookmarkRepository BookmarkRepository
	webRepository      WebRepository
}

// NewBookmarkService creates a new instance of BookmarkService
func NewBookmarkService(bookmarkRepo BookmarkRepository, webrepo WebRepository) *BookmarkService {
	return &BookmarkService{
		bookmarkRepository: bookmarkRepo,
		webRepository:      webrepo,
	}
}

func (s *BookmarkService) CreateBookmark(b model.Bookmark) (model.Bookmark, error) {
	if b.ID != "" {
		return model.Bookmark{}, fmt.Errorf("bookmark ID must be empty")
	}

	webTitle, err := s.webRepository.GetTitle(b.URL)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to fetch title for URL %s: %w", b.URL, err)
	}
	mainImageURL, err := s.webRepository.GetMainImage(b.URL)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to fetch main image URL for URL %s: %w", b.URL, err)
	}
	b.Title = webTitle
	b.MainImageURL = mainImageURL
	b.IsArchived = false
	createdBookmark, err := s.bookmarkRepository.CreateBookmark(b)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to create bookmark for URL %s: %w", b.URL, err)
	}
	logger.Info("Created bookmark",
		zap.String("id", createdBookmark.ID),
		zap.String("user_id", createdBookmark.UserID),
		zap.String("url", createdBookmark.URL),
		zap.String("title", createdBookmark.Title))

	return createdBookmark, nil
}

func (s *BookmarkService) ArchiveBookmark(id string) (model.Bookmark, error) {
	b, err := s.bookmarkRepository.GetBookmark(id)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to get bookmark with ID %s: %w", id, err)
	}
	b.IsArchived = true
	updated, err := s.bookmarkRepository.UpdateBookmark(b)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to update bookmark with ID %s: %w", b.ID, err)
	}

	return updated, nil
}

func (s *BookmarkService) GetBookmark(id string) (model.Bookmark, error) {
	if id == "" {
		return model.Bookmark{}, fmt.Errorf("id is required")
	}
	bookmarks, err := s.bookmarkRepository.GetBookmark(id)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to get bookmarks for id %s: %w", id, err)
	}

	return bookmarks, nil
}

func (s *BookmarkService) GetAllBookmarks(userID string, archived bool) ([]model.Bookmark, error) {
	query := model.BookmarkQuery{
		UserID:   userID,
		Archived: archived,
	}
	bookmarks, err := s.bookmarkRepository.ListBookmarks(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all bookmarks: %w", err)
	}

	return bookmarks, nil
}

// GetBookmarksWithPagination retrieves bookmarks with pagination support
func (s *BookmarkService) GetBookmarksWithPagination(userID string, archived bool, page, pageSize int) (model.BookmarkListResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20 // Default page size
	}
	if pageSize > 100 {
		pageSize = 100 // Maximum page size
	}

	query := model.BookmarkQuery{
		UserID:   userID,
		Archived: archived,
		Page:     page,
		PageSize: pageSize,
	}

	// Get paginated bookmarks
	bookmarks, err := s.bookmarkRepository.ListBookmarks(query)
	if err != nil {
		return model.BookmarkListResponse{}, fmt.Errorf("failed to get bookmarks: %w", err)
	}

	// Get total count
	totalCount, err := s.bookmarkRepository.CountBookmarks(query)
	if err != nil {
		return model.BookmarkListResponse{}, fmt.Errorf("failed to count bookmarks: %w", err)
	}

	// Calculate total pages
	totalPages := (totalCount + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	response := model.BookmarkListResponse{
		Bookmarks:  bookmarks,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return response, nil
}

func (s *BookmarkService) DeleteBookmark(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	err := s.bookmarkRepository.DeleteBookmark(id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark with ID %s: %w", id, err)
	}

	return nil
}
