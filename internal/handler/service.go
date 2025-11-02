package handler

import "github.com/tsongpon/athena/internal/model"

type UserService interface {
	AuthenticateUser(email, password string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
}

type BookmarkService interface {
	CreateBookmark(b model.Bookmark) (model.Bookmark, error)
	GetBookmark(id string) (model.Bookmark, error)
	DeleteBookmark(id string) error
	GetAllBookmarks(userID string, archived bool) ([]model.Bookmark, error)
	ArchiveBookmark(id string) (model.Bookmark, error)
}
