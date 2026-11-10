package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/SonDo580/bookmark-manager/internal/api"
	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/handler"
	"github.com/SonDo580/bookmark-manager/internal/service"
	redisPkg "github.com/SonDo580/bookmark-manager/pkg/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURL(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		ServiceName: "bookmark_service",
		InstanceID:  uuid.New().String(),
		Port:        8080,
	}

	correctMethod := http.MethodPost
	correctEndpoint := "/v1/links/shorten"

	const testURL = "https://facebook.com"
	const testExpSeconds = 3600

	errMsgInvalidInput := "Invalid Input"

	newJSONRequest := func(method, endpoint string, payload any) *http.Request {
		var buf bytes.Buffer
		if payload != nil {
			err := json.NewEncoder(&buf).Encode(payload)
			require.NoError(t, err)
		}

		req := httptest.NewRequest(method, endpoint, &buf)
		req.Header.Set("Content-Type", "application/json")
		return req
	}

	testCases := []struct {
		name             string
		cfg              *config.Config
		setupTestRequest func() *http.Request
		expectedStatus   int
		expectedErrMsg   string
	}{
		{
			name: "success",
			setupTestRequest: func() *http.Request {
				return newJSONRequest(correctMethod, correctEndpoint, handler.ShortenURLRequest{
					URL:        testURL,
					ExpSeconds: testExpSeconds,
				})
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "fail: invalid request body (negative exp_seconds)",
			setupTestRequest: func() *http.Request {
				return newJSONRequest(correctMethod, correctEndpoint, handler.ShortenURLRequest{
					URL:        testURL,
					ExpSeconds: -10,
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: errMsgInvalidInput,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, mockRedisClient := redisPkg.InitMockRedis(t)
			defer mockRedisClient.Close()

			testAPIEngine := api.NewEngine(cfg, mockRedisClient)
			rec := httptest.NewRecorder()
			req := tc.setupTestRequest()
			testAPIEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			switch tc.expectedStatus {
			case http.StatusCreated:
				var actualResponse handler.ShortenURLResponse
				err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
				require.NoError(t, err)

				assert.NotEmpty(t, actualResponse.Code)
				assert.Equal(t, service.DefaultCodeLength, len(actualResponse.Code))
			case http.StatusBadRequest:
				var actualResponse dto.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedErrMsg, actualResponse.Error)
			}
		})
	}
}
