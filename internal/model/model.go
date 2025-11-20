package model

// RequestGetAPIShortURL represents a request to shorten a URL.
type RequestGetAPIShortURL struct {
	URL string `json:"url"`
}

// ResponseGetAPIShortURL represents a response containing a shortened URL.
type ResponseGetAPIShortURL struct {
	Result string `json:"result"`
}

// URLStorage represents the storage structure for a URL.
type URLStorage struct {
	UUID        string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"long_url"`
	UserID      string `json:"user_id" db:"user_id"`
	DeletedFlag bool   `json:"is_deleted" db:"is_deleted"`
}

// MultiReq represents a request to shorten multiple URLs.
type MultiReq struct {
	CorrID string `json:"correlation_id"`
	URL    string `json:"original_url"`
}

// MultiResp represents a response containing a shortened URL for a multi-request.
type MultiResp struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}
