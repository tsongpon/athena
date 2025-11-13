package transport

import "time"

type BookmarkTransport struct {
	ID           string    `json:"id"`
	URL          string    `json:"url"`
	Title        string    `json:"title"`
	UserID       string    `json:"user_id"`
	MainImageURL string    `json:"main_image_url"`
	CreatedAt    time.Time `json:"created_at"`
	IsArchived   bool      `json:"is_archived"`
}
