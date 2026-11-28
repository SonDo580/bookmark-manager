package request

import (
	"net/http"

	"github.com/SonDo580/bookmark-manager/pkg/response"
	"github.com/gin-gonic/gin"
)

func bindWithAbort[T any](c *gin.Context, shouldBindFunc func(*T) error) (*T, error) {
	var input T
	err := shouldBindFunc(&input)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			response.InputFieldError(err),
		)
		return nil, err
	}

	return &input, nil
}

func BindJSONWithAbort[T any](c *gin.Context) (*T, error) {
	return bindWithAbort(c, func(input *T) error {
		return c.ShouldBindJSON(input)
	})
}

func BindURIWithAbort[T any](c *gin.Context) (*T, error) {
	return bindWithAbort(c, func(input *T) error {
		return c.ShouldBindUri(input)
	})
}

func BindHeaderWithAbort[T any](c *gin.Context) (*T, error) {
	return bindWithAbort(c, func(input *T) error {
		return c.ShouldBindHeader(input)
	})
}

func BindQueryWithAbort[T any](c *gin.Context) (*T, error) {
	return bindWithAbort(c, func(input *T) error {
		return c.ShouldBindQuery(input)
	})
}
