package handler

import (
	"errors"
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/service"
	"github.com/SonDo580/bookmark-manager/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserHandler interface {
	Register(c *gin.Context)
}

type userHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return &userHandler{svc: svc}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterRequest  true  "Input required"
// @Success      201      {object}  dto.RegisterResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      409      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /v1/users/register [post]
func (h *userHandler) Register(c *gin.Context) {
	input := &dto.RegisterRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			response.InputFieldError(err),
		)
		return
	}

	user, err := h.svc.Register(c, *input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDuplicateUsername):
			c.AbortWithStatusJSON(
				http.StatusConflict,
				dto.ErrorResponse{Message: constants.MsgUsernameAlreadyExists},
			)
		case errors.Is(err, service.ErrDuplicateEmail):
			c.AbortWithStatusJSON(
				http.StatusConflict,
				dto.ErrorResponse{Message: constants.MsgEmailAlreadyExists},
			)
		default:
			log.Error().Err(err).Msg("Register failed")

			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				dto.ErrorResponse{Message: constants.MsgInternalServerError},
			)
		}

		return
	}

	c.JSON(
		http.StatusCreated,
		dto.RegisterResponse{
			Message: "Register user successfully!",
			Data:    user,
		},
	)
}
