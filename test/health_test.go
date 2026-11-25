package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/api"
	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	dbPkg "github.com/SonDo580/bookmark-manager/pkg/db"
	redisPkg "github.com/SonDo580/bookmark-manager/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	cfg := getDefaultConfig()

	const (
		correctMethod   = http.MethodGet
		correctEndpoint = "/health-check"
	)

	testCases := []struct {
		name               string
		setupTestRequest   func() *http.Request
		expectedStatus     int
		shouldValidateBody bool
	}{
		{
			name: "success",
			setupTestRequest: func() *http.Request {
				return httptest.NewRequest(correctMethod, correctEndpoint, nil)
			},
			expectedStatus:     http.StatusOK,
			shouldValidateBody: true,
		},
		{
			name: "wrong endpoint",
			setupTestRequest: func() *http.Request {
				return httptest.NewRequest(correctMethod, correctEndpoint+"x", nil)
			},
			expectedStatus:     http.StatusNotFound,
			shouldValidateBody: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, mockRedisClient := redisPkg.InitMockRedis(t)
			defer mockRedisClient.Close()

			testDB := dbPkg.CreateTestDB(t)

			testAPIEngine := api.NewEngine(cfg, mockRedisClient, testDB)
			rec := httptest.NewRecorder()
			req := tc.setupTestRequest()
			testAPIEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if !tc.shouldValidateBody {
				return
			}

			expectedResp := dto.HealthCheckResponse{
				Message:     constants.MsgOK,
				ServiceName: cfg.ServiceName,
				InstanceID:  cfg.InstanceID,
			}

			var actualResp dto.HealthCheckResponse
			err := json.Unmarshal(rec.Body.Bytes(), &actualResp)
			require.NoError(t, err)

			assert.Equal(t, expectedResp, actualResp)
		})
	}
}
