package service

import "github.com/tsongpon/athena/internal/model"

type UserRepository interface {
	CreateUser(user model.User) (model.User, error)
	GetUserByID(id string) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
}

type BookmarkRepository interface {
	CreateBookmark(bookmark model.Bookmark) (model.Bookmark, error)
	GetBookmark(id string) (model.Bookmark, error)
	ListBookmarks(query model.BookmarkQuery) ([]model.Bookmark, error)
	UpdateBookmark(bookmark model.Bookmark) (model.Bookmark, error)
	DeleteBookmark(id string) error
}
