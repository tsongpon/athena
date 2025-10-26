package service

import (
	"fmt"
	"log"

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

func (s *BookmarkService) CreateBookmark(b model.Bookmark) (model.Bookmark, error) {
	if b.ID != "" {
		return model.Bookmark{}, fmt.Errorf("bookmark ID must be empty")
	}
	//TODO: fetch title and set it to bookmark.Title
	createdBookmark, err := s.repository.CreateBookmark(b)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to create bookmark for URL %s: %w", b.URL, err)
	}
	log.Printf("Created bookmark with ID %s", createdBookmark.ID)

	return createdBookmark, nil
}

func (s *BookmarkService) ArchiveBookmark(id string) (model.Bookmark, error) {
	b, err := s.repository.GetBookmark(id)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to get bookmark with ID %s: %w", id, err)
	}
	b.IsArchived = true
	updated, err := s.repository.UpdateBookmark(b)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to update bookmark with ID %s: %w", b.ID, err)
	}

	return updated, nil
}

func (s *BookmarkService) GetBookmark(id string) (model.Bookmark, error) {
	if id == "" {
		return model.Bookmark{}, fmt.Errorf("id is required")
	}
	bookmarks, err := s.repository.GetBookmark(id)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("failed to get bookmarks for id %s: %w", id, err)
	}

	return bookmarks, nil
}

func (s *BookmarkService) GetAllBookmarks(userID string) ([]model.Bookmark, error) {
	bookmarks, err := s.repository.ListBookmarks(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all bookmarks: %w", err)
	}

	return bookmarks, nil
}

func (s *BookmarkService) DeleteBookmark(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	err := s.repository.DeleteBookmark(id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark with ID %s: %w", id, err)
	}

	return nil
}
