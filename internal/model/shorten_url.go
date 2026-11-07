package model

type ShortenURLRequest struct {
	URL        string `json:"url"`
	ExpSeconds int    `json:"exp_seconds"`
}

type ShortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
