package service

import (
	"fmt"

	"github.com/tsongpon/athena/internal/model"
)

type BookmarkService struct {
	repository BookmarkRepository
}

func NewBookmarkService(repo BookmarkRepository) BookmarkService {
	return BookmarkService{
		repository: repo,
	}
}

func (s *BookmarkService) CreateBookmark(userID, url string) (model.Bookmark, error) {
	b := model.Bookmark{UserID: userID, URL: url}
	createdBookmark, err := s.repository.CreateBookmark(b)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to create bookmark for URL %s: %w", url, err)
	}

	return createdBookmark, nil
}

func (s *BookmarkService) ArchiveBookmark(userID, bookmarkID string) (model.Bookmark, error) {
	bookmark, err := s.repository.GetBookmark(bookmarkID)
	if err != nil {
		return bookmark, fmt.Errorf("failed to get bookmark with ID %s: %w", bookmarkID, err)
	}

	if bookmark.UserID != userID {
		return bookmark, fmt.Errorf("user %s is not authorized to archive bookmark %s", userID, bookmarkID)
	}

	bookmark.IsArchived = true
	updated, err := s.repository.UpdateBookmark(bookmark)
	if err != nil {
		return bookmark, fmt.Errorf("failed to update bookmark with ID %s: %w", bookmarkID, err)
	}

	return updated, nil
}
