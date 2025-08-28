package model

type RequestGetAPIShortURL struct {
	URL string `json:"url"`
}

type ResponseGetAPIShortURL struct {
	Result string `json:"result"`
}

type URLStorage struct {
	Uuid        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}
