package model

import "time"

type Bookmark struct {
	ID             string
	UserID         string
	URL            string
	Title          string
	IsArchived     bool
	MainImageURL   string
	ContentSummary string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// BookmarkQuery represents query parameters for listing bookmarks
type BookmarkQuery struct {
	UserID   string
	Archived bool
	Page     int // Page number (1-based), 0 means no pagination
	PageSize int // Number of items per page, 0 means no pagination
}

// BookmarkListResponse represents paginated response for listing bookmarks
type BookmarkListResponse struct {
	Bookmarks  []Bookmark `json:"bookmarks"`
	TotalCount int        `json:"total_count"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}
