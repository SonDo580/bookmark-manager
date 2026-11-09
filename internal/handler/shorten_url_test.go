package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SonDo580/bookmark-manager/internal/model"
	service_mocks "github.com/SonDo580/bookmark-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURL_CreateShortenedLink(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const testURL = "https://facebook.com"
	const testCode = "abc1234"
	const testExpSeconds = 3600
	const testExp = testExpSeconds * time.Second

	responseInvalidInput := model.ErrorResponse{
		Error: "Invalid Input",
	}

	responseSuccess := func(code string) model.ShortenURLResponse {
		return model.ShortenURLResponse{
			Code:    code,
			Message: "Shorten URL generated successfully!",
		}
	}

	testCases := []struct {
		name             string
		input            any
		setupMocks       func(mockSvc *service_mocks.ShortenURL, ctx context.Context)
		expectedStatus   int
		expectedResponse any
	}{
		{
			name: "success",
			input: model.ShortenURLRequest{
				URL:        testURL,
				ExpSeconds: testExpSeconds,
			},
			setupMocks: func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {
				mockSvc.On("CreateShortenedLink", ctx, testURL, testExp).
					Return(testCode, nil).
					Once()
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: responseSuccess(testCode),
		},
		{
			name:             "fail: malformed input",
			input:            "{invalid-json",
			setupMocks:       func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: responseInvalidInput,
		},
		{
			name: "fail: wrong input field types",
			input: map[string]any{
				"url":         1234,
				"exp_seconds": "NaN",
			},
			setupMocks:       func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: responseInvalidInput,
		},
		{
			name: "fail: missing required url",
			input: model.ShortenURLRequest{
				ExpSeconds: 3600,
			},
			setupMocks:       func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: responseInvalidInput,
		},
		{
			name: "fail: negative exp_seconds",
			input: model.ShortenURLRequest{
				URL:        testURL,
				ExpSeconds: -10,
			},
			setupMocks:       func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: responseInvalidInput,
		},
		{
			name: "fail: exceed exp_seconds limit",
			input: model.ShortenURLRequest{
				URL:        testURL,
				ExpSeconds: 86400 + 1,
			},
			setupMocks:       func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: responseInvalidInput,
		},
		{
			name: "success: omit optional exp_seconds",
			input: model.ShortenURLRequest{
				URL: testURL,
			},
			setupMocks: func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {
				mockSvc.On("CreateShortenedLink", ctx, testURL, time.Duration(0)).
					Return(testCode, nil).
					Once()
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: responseSuccess(testCode),
		},
		{
			name: "fail: service error",
			input: model.ShortenURLRequest{
				URL:        testURL,
				ExpSeconds: testExpSeconds,
			},
			setupMocks: func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {
				mockSvc.On("CreateShortenedLink", ctx, testURL, testExp).
					Return("", errors.New("fail to store")).
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

			// Serialize request body
			var inputBytes []byte
			var err error
			switch input := tc.input.(type) {
			case string:
				inputBytes = []byte(input)
			default:
				inputBytes, err = json.Marshal(input)
				require.NoError(t, err)
			}

			// Create request & Gin test context
			req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(inputBytes))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request = req

			// Setup mocks & Create handler
			mockSvc := service_mocks.NewShortenURL(t)
			tc.setupMocks(mockSvc, ginCtx)

			handler := NewShortenURL(mockSvc)

			// Execute handler
			handler.CreateShortenedLink(ginCtx)

			// Assert response status
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response body
			switch expectedResponse := tc.expectedResponse.(type) {
			case model.ShortenURLResponse:
				var actualResponse model.ShortenURLResponse
				err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
				require.NoError(t, err)
				assert.Equal(t, expectedResponse, actualResponse)
			case model.ErrorResponse:
				var actualResponse model.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
				require.NoError(t, err)
				assert.Equal(t, expectedResponse, actualResponse)
			}

			// Verify all mock expectations are met
			mockSvc.AssertExpectations(t)
		})
	}
}
