package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	jwtutils_mocks "github.com/SonDo580/bookmark-manager/pkg/jwtutils/mocks"
	"github.com/SonDo580/bookmark-manager/pkg/request"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTMiddleware_Auth(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const (
		testToken  = "<valid_jwt>"
		testUserID = "<user_id>"
	)

	testClaims := jwt.MapClaims{"sub": testUserID}

	testCases := []struct {
		name           string
		authHeader     string
		setupMock      func(validator *jwtutils_mocks.JWTValidator)
		expectedErrMsg string
		expectedClaims jwt.MapClaims
	}{
		{
			name:       "success",
			authHeader: "Bearer " + testToken,
			setupMock: func(validator *jwtutils_mocks.JWTValidator) {
				validator.On("Validate", testToken).
					Return(testClaims, nil).
					Once()
			},
			expectedClaims: testClaims,
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			setupMock:      func(validator *jwtutils_mocks.JWTValidator) {},
			expectedErrMsg: constants.MsgLackAuthHeader,
		},
		{
			name:           "wrong authorization header format",
			authHeader:     testToken, // no Bearer prefix
			setupMock:      func(validator *jwtutils_mocks.JWTValidator) {},
			expectedErrMsg: constants.MsgWrongAuthHeaderFormat,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer " + testToken,
			setupMock: func(validator *jwtutils_mocks.JWTValidator) {
				validator.On("Validate", testToken).
					Return(nil, errors.New("invalid token")).
					Once()
			},
			expectedErrMsg: constants.MsgInvalidToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockJwtValidator := jwtutils_mocks.NewJWTValidator(t)
			tc.setupMock(mockJwtValidator)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.authHeader != "" {
				c.Request.Header.Set("Authorization", tc.authHeader)
			}

			jwtMiddleware := NewJWTMiddleware(mockJwtValidator)
			handler := jwtMiddleware.Auth()
			handler(c)

			if tc.expectedErrMsg != "" {
				assert.True(t, c.IsAborted())
				assert.Equal(t, http.StatusUnauthorized, rec.Code)

				var resp dto.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			} else {
				assert.False(t, c.IsAborted())
				claims, err := request.ExtractClaimsFromCtx(c)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedClaims, claims)
			}

			mockJwtValidator.AssertExpectations(t)
		})
	}
}
