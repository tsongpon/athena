package model

import "time"

type Bookmark struct {
	ID         string
	UserID     string
	URL        string
	Title      string
	IsArchived bool
	CreatedAt  time.Time
}

// BookmarkQuery represents query parameters for listing bookmarks
type BookmarkQuery struct {
	UserID   string
	Archived bool
}
