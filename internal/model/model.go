package model

type RequestGetAPIShortURL struct {
	URL string `json:"url"`
}

type ResponseGetAPIShortURL struct {
	Result string `json:"result"`
}

type URLStorage struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
