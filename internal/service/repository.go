package service

import "github.com/tsongpon/athena/internal/model"

type UserRepository interface {
	CreateUser(user model.User) (model.User, error)
	GetUserByID(id string) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	GetUserByEmailAndPassword(email, hashedPassword string) (model.User, error)
}

type BookmarkRepository interface {
	CreateBookmark(bookmark model.Bookmark) (model.Bookmark, error)
	GetBookmark(id string) (model.Bookmark, error)
	ListBookmarks(query model.BookmarkQuery) ([]model.Bookmark, error)
	CountBookmarks(query model.BookmarkQuery) (int, error)
	UpdateBookmark(bookmark model.Bookmark) (model.Bookmark, error)
	DeleteBookmark(id string) error
}

type WebRepository interface {
	GetTitle(url string) (string, error)
	GetMainImage(url string) (string, error)
	GetContentSummary(url string) (string, error)
}
