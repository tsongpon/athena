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
