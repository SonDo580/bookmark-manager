package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	service_mocks "github.com/SonDo580/bookmark-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		ServiceName: "bookmark_service",
		InstanceID:  uuid.New().String(),
		Port:        8080,
	}

	successResp := dto.HealthCheckResponse{
		Message:     constants.MsgOK,
		ServiceName: cfg.ServiceName,
		InstanceID:  cfg.InstanceID,
	}

	testCases := []struct {
		name           string
		setupMocks     func(mockSvc *service_mocks.Health, ctx context.Context)
		expectedStatus int
		expectedErrMsg string
	}{
		{
			name: "success",
			setupMocks: func(mockSvc *service_mocks.Health, ctx context.Context) {
				mockSvc.On("Check", ctx).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "fail: redis down",
			setupMocks: func(mockSvc *service_mocks.Health, ctx context.Context) {
				mockSvc.On("Check", ctx).
					Return(errors.New("redis down")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErrMsg: constants.MsgInternalServerError,
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

			healthHandler := NewHealth(cfg, mockHealthSvc)

			// Execute handler
			healthHandler.HealthCheck(ginCtx)

			// Assert response status
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response body
			if tc.expectedErrMsg != "" {
				var resp dto.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			} else {
				var resp dto.HealthCheckResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, successResp, resp)
			}

			// Verify all mock expectations are met
			mockHealthSvc.AssertExpectations(t)
		})
	}
}
