package service

import (
	"fmt"
	"os"
	"sync"

	"github.com/tsongpon/athena/internal/logger"
	"github.com/tsongpon/athena/internal/model"
	"go.uber.org/zap"
)

// bookmarkService is the concrete implementation of BookmarkService interface
type BookmarkService struct {
	bookmarkRepository BookmarkRepository
	userRepository     UserRepository
	webRepository      WebRepository
}

// NewBookmarkService creates a new instance of BookmarkService
func NewBookmarkService(bookmarkRepo BookmarkRepository, userRepo UserRepository, webrepo WebRepository) *BookmarkService {
	return &BookmarkService{
		bookmarkRepository: bookmarkRepo,
		userRepository:     userRepo,
		webRepository:      webrepo,
	}
}

func (s *BookmarkService) CreateBookmark(b model.Bookmark) (model.Bookmark, error) {
	if b.ID != "" {
		return model.Bookmark{}, fmt.Errorf("bookmark ID must be empty")
	}
	var wg sync.WaitGroup
	var title string
	var imageURL string

	wg.Go(func() {
		var e error
		title, e = s.webRepository.GetTitle(b.URL)
		if e != nil {
			logger.Error("failed to fetch title for URL", zap.String("url", b.URL), zap.Error(e))
		}
	})

	wg.Go(func() {
		var e error
		imageURL, e = s.webRepository.GetMainImage(b.URL)
		if e != nil {
			logger.Error("failed to fetch main image URL for URL", zap.String("url", b.URL), zap.Error(e))
		}
	})

	var content string
	user, err := s.userRepository.GetUserByID(b.UserID)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to fetch user for ID %s: %w", b.UserID, err)
	}
	if user.Tier == "paid" {
		llmSummaryContent := os.Getenv("LLM_SUMMARY_CONTENT")
		if llmSummaryContent == "true" {
			wg.Go(func() {
				logger.Info("LLM content summary is enabled")
				content, err = s.webRepository.GetContentSummary(b.URL)
				if err != nil {
					logger.Error("failed to fetch content summary for URL", zap.String("url", b.URL), zap.Error(err))
				}
			})
		}
	}
	wg.Wait()
	b.Title = title
	b.ContentSummary = content
	b.MainImageURL = imageURL
	b.IsArchived = false
	createdBookmark, err := s.bookmarkRepository.CreateBookmark(b)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to create bookmark for URL %s: %w", b.URL, err)
	}
	logger.Info("Created bookmark",
		zap.String("id", createdBookmark.ID),
		zap.String("user_id", createdBookmark.UserID),
		zap.String("url", createdBookmark.URL),
		zap.String("title", createdBookmark.Title),
		zap.String("content_summary", createdBookmark.ContentSummary),
		zap.String("main_image_url", createdBookmark.MainImageURL),
		zap.Bool("is_archived", createdBookmark.IsArchived))

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
