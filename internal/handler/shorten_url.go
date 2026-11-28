package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/service"
	"github.com/SonDo580/bookmark-manager/pkg/request"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ShortenURL interface {
	CreateShortenedLink(c *gin.Context)
	Redirect(c *gin.Context)
}

type shortenURL struct {
	svc service.ShortenURL
}

func NewShortenURL(svc service.ShortenURL) ShortenURL {
	return &shortenURL{svc: svc}
}

// CreateShortenedLink godoc
// @Summary      Generate a shortened URL
// @Description  Generate a shortened URL from the original URL.
// @Tags         link
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ShortenURLRequest true "Input required"
// @Success      201      {object}  dto.ShortenURLResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /v1/links/shorten [post]
func (h *shortenURL) CreateShortenedLink(c *gin.Context) {
	input, err := request.BindJSONWithAbort[dto.ShortenURLRequest](c)
	if err != nil {
		return
	}

	exp := time.Duration(input.ExpSeconds) * time.Second
	code, err := h.svc.CreateShortenedLink(c, input.URL, exp)
	if err != nil {
		log.Error().Err(err).Msg("CreateShortenedLink failed")

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			dto.ErrorResponse{Message: constants.MsgInternalServerError},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		dto.ShortenURLResponse{
			Code:    code,
			Message: "Shorten URL generated successfully!",
		},
	)
}

type redirectPathParams struct {
	Code string `uri:"code" binding:"required"`
}

// Redirect godoc
// @Summary      Redirect to original URL
// @Description  Redirect to original URL.
// @Tags         link
// @Param        code   path		string true "Shortened code"
// @Success      301    {string}	string "Moved permanently"
// @Failure      400    {object}  	dto.ErrorResponse
// @Failure      500    {object}  	dto.ErrorResponse
// @Router       /v1/links/redirect/{code} [get]
func (h *shortenURL) Redirect(c *gin.Context) {
	input, err := request.BindURIWithAbort[redirectPathParams](c)
	if err != nil {
		return
	}
	code := input.Code

	url, err := h.svc.GetURLFromCode(c, code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				dto.ErrorResponse{Message: constants.MsgURLNotFound},
			)
			return
		}

		log.Error().Err(err).Msg("Redirect failed")

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			dto.ErrorResponse{Message: constants.MsgInternalServerError},
		)
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}
