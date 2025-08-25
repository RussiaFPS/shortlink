package model

type RequestGetAPIShortURL struct {
	URL string `json:"url"`
}

type ResponseGetAPIShortURL struct {
	Result string `json:"result"`
}
