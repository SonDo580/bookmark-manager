package dto

type ShortenURLRequest struct {
	URL        string `json:"url" binding:"required,url"`
	ExpSeconds int    `json:"exp_seconds" binding:"gt=0,lte=86400"`
}

type ShortenURLResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
