package model

type RequestGetAPIShortURL struct {
	URL string `json:"url"`
}

type ResponseGetAPIShortURL struct {
	Result string `json:"result"`
}

type URLStorage struct {
	UUID        string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"long_url"`
	UserID      string `json:"user_id" db:"user_id"`
	DeletedFlag bool   `json:"is_deleted" db:"is_deleted"`
}

type MultiReq struct {
	CorrID string `json:"correlation_id"`
	URL    string `json:"original_url"`
}

type MultiResp struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}
