package middleware

import (
	"net/http"
	"strings"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/pkg/jwtutils"
	"github.com/SonDo580/bookmark-manager/pkg/request"
	"github.com/gin-gonic/gin"
)

type JWTMiddleware interface {
	Auth() gin.HandlerFunc
}

type jwtMiddleware struct {
	jwtValidator jwtutils.JWTValidator
}

func NewJWTMiddleware(jwtValidator jwtutils.JWTValidator) JWTMiddleware {
	return &jwtMiddleware{jwtValidator: jwtValidator}
}

func (m *jwtMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				dto.ErrorResponse{Message: constants.MsgLackAuthHeader},
			)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				dto.ErrorResponse{Message: constants.MsgWrongAuthHeaderFormat},
			)
			return
		}

		tokenStr := parts[1]

		// Validate token
		claims, err := m.jwtValidator.Validate(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				dto.ErrorResponse{Message: constants.MsgInvalidToken},
			)
			return
		}

		// Save data in context
		request.SetClaimsInCtx(c, claims)

		c.Next()
	}
}
