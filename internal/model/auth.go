package model

// UserURL represents a user's URL mapping.
type UserURL struct {
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

// DellURL represents a request to delete a URL.
type DellURL struct {
	UserID string
	URLs   []string
}

// RespUserURL represents a response containing a user's URL mapping.
type RespUserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
