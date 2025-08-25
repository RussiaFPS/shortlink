package model

type RequestGetAPIShortURL struct {
	Url string `json:"url"`
}

type ResponseGetAPIShortURL struct {
	Result string `json:"result"`
}
