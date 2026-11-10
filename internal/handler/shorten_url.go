package handler

import (
	"net/http"
	"time"

	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/service"
	"github.com/gin-gonic/gin"
)

type ShortenURL interface {
	CreateShortenedLink(c *gin.Context)
}

type shortenURL struct {
	svc service.ShortenURL
}

func NewShortenURL(svc service.ShortenURL) ShortenURL {
	return &shortenURL{svc: svc}
}

type ShortenURLRequest struct {
	URL        string `json:"url" binding:"required,url"`
	ExpSeconds int    `json:"exp_seconds" binding:"gte=0,lte=86400"`
}

type ShortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CreateShortenedLink godoc
// @Summary      Generate a shortened URL
// @Description  Generate a shortened URL from the original URL.
// @Tags         link
// @Accept       json
// @Produce      json
// @Param        request  body      ShortenURLRequest true "Input required"
// @Success      201      {object}  ShortenURLResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /v1/links/shorten [post]
func (h *shortenURL) CreateShortenedLink(c *gin.Context) {
	input := &ShortenURLRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			dto.ErrorResponse{Error: "Invalid Input"},
		)
		return
	}

	exp := time.Duration(input.ExpSeconds) * time.Second
	code, err := h.svc.CreateShortenedLink(c, input.URL, exp)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			dto.ErrorResponse{Error: "Internal Server Error"},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		ShortenURLResponse{
			Code:    code,
			Message: "Shorten URL generated successfully!",
		},
	)
}
