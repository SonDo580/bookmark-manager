package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/api"
	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/service"
	dbPkg "github.com/SonDo580/bookmark-manager/pkg/db"
	jwtutils_mocks "github.com/SonDo580/bookmark-manager/pkg/jwtutils/mocks"
	redisPkg "github.com/SonDo580/bookmark-manager/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURL_CreateShortenedLink(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	cfg := &config.Config{}

	const (
		correctMethod   = http.MethodPost
		correctEndpoint = "/v1/links/shorten"

		testURL        = "https://facebook.com"
		testExpSeconds = 3600
	)

	testCases := []struct {
		name             string
		setupTestRequest func(t *testing.T) *http.Request
		expectedStatus   int
		expectedErrMsg   string
	}{
		{
			name: "success",
			setupTestRequest: func(t *testing.T) *http.Request {
				return newJSONRequest(t, correctMethod, correctEndpoint, dto.ShortenURLRequest{
					URL:        testURL,
					ExpSeconds: testExpSeconds,
				})
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "fail: invalid request body (non-positive exp_seconds)",
			setupTestRequest: func(t *testing.T) *http.Request {
				return newJSONRequest(t, correctMethod, correctEndpoint, dto.ShortenURLRequest{
					URL:        testURL,
					ExpSeconds: 0,
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: constants.MsgInvalidInput,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, mockRedisClient := redisPkg.InitMockRedis(t)
			defer mockRedisClient.Close()

			testDB := dbPkg.CreateTestDB(t)

			mockJWTGenerator := jwtutils_mocks.NewJWTGenerator(t)
			mockJWTValidator := jwtutils_mocks.NewJWTValidator(t)

			testAPIEngine := api.NewEngine(cfg, mockRedisClient, testDB, mockJWTGenerator, mockJWTValidator)

			rec := httptest.NewRecorder()
			req := tc.setupTestRequest(t)
			testAPIEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			switch tc.expectedStatus {
			case http.StatusCreated:
				var resp dto.ShortenURLResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.NotEmpty(t, resp.Code)
				assert.Equal(t, service.DefaultCodeLength, len(resp.Code))
			case http.StatusBadRequest:
				var resp dto.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}
		})
	}
}

func TestShortenURL_Redirect(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	cfg := &config.Config{}

	const (
		testCode = "abc1234"
		testURL  = "https://google.com"
	)

	testCases := []struct {
		name           string
		codeParam      string
		expectedStatus int
		expectedURL    string
		expectedErrMsg string
	}{
		{
			name:           "success",
			codeParam:      testCode,
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    testURL,
		},
		{
			name:           "fail: URL not found",
			codeParam:      testCode + "x",
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: constants.MsgURLNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, mockRedisClient := redisPkg.InitMockRedis(t)
			defer mockRedisClient.Close()

			testDB := dbPkg.CreateTestDB(t)

			// Store code-to-URL mapping
			ctx := t.Context()
			err := mockRedisClient.Set(ctx, testCode, testURL, 0).Err()
			require.NoError(t, err)

			mockJWTGenerator := jwtutils_mocks.NewJWTGenerator(t)
			mockJWTValidator := jwtutils_mocks.NewJWTValidator(t)

			testAPIEngine := api.NewEngine(cfg, mockRedisClient, testDB, mockJWTGenerator, mockJWTValidator)

			rec := httptest.NewRecorder()
			endpoint := "/v1/links/redirect/" + tc.codeParam
			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			testAPIEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedURL != "" {
				assert.Equal(t, tc.expectedURL, rec.Header().Get("Location"))
			}

			if tc.expectedErrMsg != "" {
				var resp dto.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}
		})
	}
}
