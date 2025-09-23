package model

type UserURL struct {
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

type DellURL struct {
	UserID string
	URLs   []string
}

type RespUserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
