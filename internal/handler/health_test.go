package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		cfg                *config.Config
		expectedStatus     int
		expectedMessage    string
		expectedService    string
		expectedInstanceID string
	}{
		{
			name: "success",
			cfg: &config.Config{
				ServiceName: "bookmark_service",
				InstanceID:  "f1fef0a1-e50b-4c52-b205-855181037ab4",
				Port:        8080,
			},
			expectedStatus:     http.StatusOK,
			expectedMessage:    "OK",
			expectedService:    "bookmark_service",
			expectedInstanceID: "f1fef0a1-e50b-4c52-b205-855181037ab4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			healthHandler := NewHealth(tc.cfg)
			healthHandler.HealthCheck(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var body map[string]any
			err := json.Unmarshal(rec.Body.Bytes(), &body)
			require.NoError(t, err)

			assert.Equal(t, "OK", body["message"])
			assert.Equal(t, tc.expectedService, body["service_name"])
			assert.Equal(t, tc.expectedInstanceID, body["instance_id"])
		})
	}
}
