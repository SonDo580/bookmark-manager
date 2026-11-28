package dto

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type DataResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}
