package handler

import (
	"errors"
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/service"
	"github.com/SonDo580/bookmark-manager/pkg/request"
	"github.com/SonDo580/bookmark-manager/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetSelfInfo(c *gin.Context)
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

// Login godoc
// @Summary      Login
// @Description  Return token if credential is correct.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Input required"
// @Success      200      {object}  dto.LoginResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /v1/users/login [post]
func (h *userHandler) Login(c *gin.Context) {
	input := &dto.LoginRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			response.InputFieldError(err),
		)
		return
	}

	tokenStr, err := h.svc.Login(c, input.Username, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWrongCredentials):
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				dto.ErrorResponse{Message: constants.MsgWrongCredentials},
			)
		default:
			log.Error().Err(err).Msg("Login failed")

			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				dto.ErrorResponse{Message: constants.MsgInternalServerError},
			)
		}
		return
	}

	c.JSON(
		http.StatusOK,
		dto.LoginResponse{
			Message: "Logged in successfully!",
			Data:    tokenStr,
		},
	)
}

// GetSelfInfo godoc
// @Summary      Get your information
// @Description  Get your information.
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200      {object}  dto.GetSelfInfoResponse
// @Failure      401      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /v1/self/info [get]
func (h *userHandler) GetSelfInfo(c *gin.Context) {
	userID, err := request.ExtractUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			dto.ErrorResponse{Message: constants.MsgInvalidToken},
		)
		return
	}

	user, err := h.svc.GetSelfInfoByID(c, userID)
	if err != nil {
		log.Error().Err(err).Msgf("GetSelfInfo failed for %s", userID)

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			dto.ErrorResponse{Message: constants.MsgInternalServerError},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dto.GetSelfInfoResponse{Data: user},
	)
}
