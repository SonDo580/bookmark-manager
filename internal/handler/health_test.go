package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/SonDo580/bookmark-manager/internal/dto"
	service_mocks "github.com/SonDo580/bookmark-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	successResponse := dto.HealthCheckResponse{
		Message:     dto.MsgOK,
		ServiceName: "bookmark_service",
		InstanceID:  uuid.New().String(),
	}

	testCases := []struct {
		name             string
		setupMocks       func(mockSvc *service_mocks.Health, ctx context.Context)
		expectedStatus   int
		expectedResponse any
	}{
		{
			name: "success",
			setupMocks: func(mockSvc *service_mocks.Health, ctx context.Context) {
				mockSvc.On("Check", ctx).
					Return(successResponse, nil).
					Once()
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: successResponse,
		},
		{
			name: "fail: redis down",
			setupMocks: func(mockSvc *service_mocks.Health, ctx context.Context) {
				mockSvc.On("Check", ctx).
					Return(dto.HealthCheckResponse{}, errors.New("redis down")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: dto.ErrorResponse{
				Message: dto.MsgInternalServerError,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create request & Gin test context
			req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request = req

			// Setup mocks & Create handler
			mockHealthSvc := service_mocks.NewHealth(t)
			tc.setupMocks(mockHealthSvc, ginCtx)

			healthHandler := NewHealth(mockHealthSvc)

			// Execute handler
			healthHandler.HealthCheck(ginCtx)

			// Assert response status
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response body
			switch expectedResponse := tc.expectedResponse.(type) {
			case dto.HealthCheckResponse:
				var actualResponse dto.HealthCheckResponse
				err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
				require.NoError(t, err)
				assert.Equal(t, expectedResponse, actualResponse)
			case dto.ErrorResponse:
				var actualResponse dto.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
				require.NoError(t, err)
				assert.Equal(t, expectedResponse, actualResponse)
			}

			// Verify all mock expectations are met
			mockHealthSvc.AssertExpectations(t)
		})
	}
}
