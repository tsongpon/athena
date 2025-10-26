package transport

import "time"

type BookmarkTransport struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
