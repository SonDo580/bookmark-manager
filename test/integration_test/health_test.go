package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/SonDo580/bookmark-manager/internal/api"
	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	redisPkg "github.com/SonDo580/bookmark-manager/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name               string
		cfg                *config.Config
		setupTestRequest   func() *http.Request
		expectedStatus     int
		shouldValidateBody bool
	}{
		{
			name: "success",
			cfg: &config.Config{
				ServiceName: "bookmark_service",
				InstanceID:  uuid.New().String(),
				Port:        8080,
			},
			setupTestRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/health-check", nil)
			},
			expectedStatus:     http.StatusOK,
			shouldValidateBody: true,
		},
		{
			name: "wrong endpoint",
			cfg: &config.Config{
				ServiceName: "bookmark_service",
				InstanceID:  uuid.New().String(),
				Port:        8080,
			},
			setupTestRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/health-check-123", nil)
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

			testAPIEngine := api.NewEngine(tc.cfg, mockRedisClient)
			rec := httptest.NewRecorder()
			req := tc.setupTestRequest()
			testAPIEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if !tc.shouldValidateBody {
				return
			}

			expectedResponse := dto.HealthCheckResponse{
				Message:     dto.HealthStatusOK,
				ServiceName: tc.cfg.ServiceName,
				InstanceID:  tc.cfg.InstanceID,
			}

			var actualResponse dto.HealthCheckResponse
			err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
			require.NoError(t, err)

			assert.Equal(t, expectedResponse, actualResponse)
		})
	}
}
