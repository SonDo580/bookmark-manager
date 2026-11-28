package dto

import (
	"github.com/SonDo580/bookmark-manager/internal/model"
)

type RegisterRequest struct {
	DisplayName string `json:"display_name" binding:"required,max=100"`
	Username    string `json:"username" binding:"required,max=100"`
	Password    string `json:"password" binding:"required,min=8,max=72"`
	Email       string `json:"email" binding:"required,email,max=255"`
}

type RegisterResponse struct {
	Message string      `json:"message"`
	Data    *model.User `json:"data"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

type GetSelfInfoResponse struct {
	Data *model.User `json:"data"`
}

type EditSelfInfoRequest struct {
	DisplayName string `json:"display_name" binding:"required,max=100"`
	Email       string `json:"email" binding:"required,email,max=255"`
}

type EditSelfInfoResponse struct {
	Message string `json:"message"`
}
