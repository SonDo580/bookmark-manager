package model

type ShortenURLRequest struct {
	URL        string `json:"url" binding:"required,url"`
	ExpSeconds int    `json:"exp_seconds" binding:"gte=0"`
}

type ShortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
