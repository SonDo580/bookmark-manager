package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/SonDo580/bookmark-manager/internal/model"
	"github.com/SonDo580/bookmark-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		expectedStatus   int
		expectedResponse model.HealthCheckResponse
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
			expectedResponse: model.HealthCheckResponse{
				Message:     model.HealthStatusOK,
				ServiceName: "bookmark_service",
				InstanceID:  uuid.New().String(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHealthSvc := mocks.NewHealth(t)
			mockHealthSvc.On("Check").Return(tc.expectedResponse).Once()

			healthHandler := NewHealth(mockHealthSvc)

			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			healthHandler.HealthCheck(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var actualResponse model.HealthCheckResponse
			err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResponse, actualResponse)
		})
	}
}
