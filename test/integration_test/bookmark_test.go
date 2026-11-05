package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/api"
	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		cfg                *config.Config
		setupTestRequest   func() *http.Request
		expectedStatus     int
		shouldValidateBody bool
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
			setupTestRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/health-check", nil)
			},
			expectedStatus:     http.StatusOK,
			shouldValidateBody: true,
			expectedMessage:    "OK",
			expectedService:    "bookmark_service",
			expectedInstanceID: "f1fef0a1-e50b-4c52-b205-855181037ab4",
		},
		{
			name: "wrong endpoint",
			cfg: &config.Config{
				ServiceName: "bookmark_service",
				InstanceID:  "f1fef0a1-e50b-4c52-b205-855181037ab4",
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

			testAPIEngine := api.NewEngine(tc.cfg)
			rec := httptest.NewRecorder()
			req := tc.setupTestRequest()
			testAPIEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if !tc.shouldValidateBody {
				return
			}

			var body map[string]any
			err := json.Unmarshal(rec.Body.Bytes(), &body)
			require.NoError(t, err)

			assert.Equal(t, "OK", body["message"])
			assert.Equal(t, tc.expectedService, body["service_name"])
			assert.Equal(t, tc.expectedInstanceID, body["instance_id"])
		})
	}
}
