package handler

import "github.com/tsongpon/athena/internal/model"

type BookmarkService interface {
	CreateBookmark(b model.Bookmark) (model.Bookmark, error)
	GetBookmark(id string) (model.Bookmark, error)
	DeleteBookmark(id string) error
	GetAllBookmarks(userID string) ([]model.Bookmark, error)
	ArchiveBookmark(id string) (model.Bookmark, error)
}
