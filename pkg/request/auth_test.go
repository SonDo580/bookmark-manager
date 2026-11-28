package request

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractUserIDFromCtxWithAbort(t *testing.T) {
	t.Parallel()

	const testUserID = "<user_id>"

	testCases := []struct {
		name        string
		setupCtx    func(c *gin.Context)
		expectedID  string
		expectedErr error
	}{
		{
			name: "success",
			setupCtx: func(c *gin.Context) {
				SetClaimsInCtx(c, jwt.MapClaims{"sub": testUserID})
			},
			expectedID: testUserID,
		},
		{
			name:        "missing claims",
			setupCtx:    func(c *gin.Context) {},
			expectedErr: ErrMissingClaims,
		},
		{
			name: "invalid claims type",
			setupCtx: func(c *gin.Context) {
				c.Set(CtxKeyClaims, "not-map-claims")
			},
			expectedErr: ErrInvalidClaims,
		},
		{
			name: "missing sub",
			setupCtx: func(c *gin.Context) {
				SetClaimsInCtx(c, jwt.MapClaims{})
			},
			expectedErr: ErrMissingSub,
		},
		{
			name: "empty sub",
			setupCtx: func(c *gin.Context) {
				SetClaimsInCtx(c, jwt.MapClaims{"sub": ""})
			},
			expectedErr: ErrMissingSub,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, rec := newTestGinCtx(t)
			tc.setupCtx(c)

			actualID, err := ExtractUserIDFromCtxWithAbort(c)
			assert.ErrorIs(t, err, tc.expectedErr)

			if tc.expectedErr != nil {
				assert.True(t, c.IsAborted())
				assert.Equal(t, http.StatusUnauthorized, rec.Code)

				var actualResp dto.ErrorResponse
				unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &actualResp)
				require.NoError(t, unmarshalErr)
				assert.Equal(t, constants.MsgInvalidToken, actualResp.Message)
			} else {
				assert.False(t, c.IsAborted())
				assert.Equal(t, tc.expectedID, actualID)
			}
		})
	}
}
