package request

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	CtxKeyClaims = "claims"
)

var (
	ErrMissingCtxClaims = errors.New("missing claims in context")
	ErrInvalidCtxClaims = errors.New("invalid claims in context")
)

func SetClaims(c *gin.Context, claims jwt.MapClaims) {
	c.Set(CtxKeyClaims, claims)
}

func ExtractClaims(c *gin.Context) (jwt.MapClaims, error) {
	raw, exists := c.Get(CtxKeyClaims)
	if !exists {
		return nil, ErrMissingCtxClaims
	}

	claims, ok := raw.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidCtxClaims
	}

	return claims, nil
}

func ExtractUserID(c *gin.Context) (string, error) {
	claims, err := ExtractClaims(c)
	if err != nil {
		return "", err
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", ErrInvalidCtxClaims
	}

	return userID, nil
}
