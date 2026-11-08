package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/SonDo580/bookmark-manager/internal/model"
	service_mocks "github.com/SonDo580/bookmark-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	successResponse := model.HealthCheckResponse{
		Message:     model.HealthStatusOK,
		ServiceName: "bookmark_service",
		InstanceID:  uuid.New().String(),
	}

	testCases := []struct {
		name             string
		setupMocks       func(mockSvc *service_mocks.Health)
		expectedStatus   int
		expectedResponse any
	}{
		{
			name: "success",
			setupMocks: func(mockSvc *service_mocks.Health) {
				mockSvc.On("Check", mock.Anything).
					Return(successResponse, nil).
					Once()
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: successResponse,
		},
		{
			name: "fail: redis down",
			setupMocks: func(mockSvc *service_mocks.Health) {
				mockSvc.On("Check", mock.Anything).
					Return(model.HealthCheckResponse{}, errors.New("redis down")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: model.ErrorResponse{
				Error: "Internal Server Error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockHealthSvc := service_mocks.NewHealth(t)
			tc.setupMocks(mockHealthSvc)

			healthHandler := NewHealth(mockHealthSvc)

			req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request = req

			healthHandler.HealthCheck(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			switch expectedResponse := tc.expectedResponse.(type) {
			case model.HealthCheckResponse:
				var actualResponse model.HealthCheckResponse
				err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
				require.NoError(t, err)
				assert.Equal(t, expectedResponse, actualResponse)
			case model.ErrorResponse:
				var actualResponse model.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
				require.NoError(t, err)
				assert.Equal(t, expectedResponse, actualResponse)
			}

			mockHealthSvc.AssertExpectations(t)
		})
	}
}
