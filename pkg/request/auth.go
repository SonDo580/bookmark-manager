package request

import (
	"errors"
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	CtxKeyClaims = "claims"
)

var (
	ErrMissingClaims = errors.New("missing claims in context")
	ErrInvalidClaims = errors.New("invalid claims in context")
	ErrMissingSub    = errors.New("missing sub in claims")
)

func SetClaimsInCtx(c *gin.Context, claims jwt.MapClaims) {
	c.Set(CtxKeyClaims, claims)
}

func ExtractClaimsFromCtx(c *gin.Context) (jwt.MapClaims, error) {
	raw, exists := c.Get(CtxKeyClaims)
	if !exists {
		return nil, ErrMissingClaims
	}

	claims, ok := raw.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

func ExtractUserIDFromCtx(c *gin.Context) (string, error) {
	claims, err := ExtractClaimsFromCtx(c)
	if err != nil {
		return "", err
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", ErrMissingSub
	}

	return userID, nil
}

func ExtractUserIDFromCtxWithAbort(c *gin.Context) (string, error) {
	userID, err := ExtractUserIDFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			dto.ErrorResponse{Message: constants.MsgInvalidToken},
		)
		return "", err
	}

	return userID, nil
}
