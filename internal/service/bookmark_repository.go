package service

import "github.com/tsongpon/athena/internal/model"

type BookmarkRepository interface {
	CreateBookmark(bookmark model.Bookmark) (model.Bookmark, error)
	GetBookmark(id string) (model.Bookmark, error)
	ListBookmarks(userID string) ([]model.Bookmark, error)
	UpdateBookmark(bookmark model.Bookmark) (model.Bookmark, error)
	DeleteBookmark(id string) error
}
