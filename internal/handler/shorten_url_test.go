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

	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/service"
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

	responseSuccess := func(code string) *ShortenURLResponse {
		return &ShortenURLResponse{
			Code:    code,
			Message: "Shorten URL generated successfully!",
		}
	}

	testCases := []struct {
		name            string
		input           any
		setupMocks      func(mockSvc *service_mocks.ShortenURL, ctx context.Context)
		expectedStatus  int
		successResponse *ShortenURLResponse
		expectedErrMsg  string
	}{
		{
			name: "success",
			input: ShortenURLRequest{
				URL:        testURL,
				ExpSeconds: testExpSeconds,
			},
			setupMocks: func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {
				mockSvc.On("CreateShortenedLink", ctx, testURL, testExp).
					Return(testCode, nil).
					Once()
			},
			expectedStatus:  http.StatusCreated,
			successResponse: responseSuccess(testCode),
		},
		{
			name:           "fail: malformed input",
			input:          "{invalid-json",
			setupMocks:     func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Invalid Input",
		},
		{
			name: "fail: wrong input field types",
			input: map[string]any{
				"url":         1234,
				"exp_seconds": "NaN",
			},
			setupMocks:     func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Invalid Input",
		},
		{
			name: "fail: missing required url",
			input: ShortenURLRequest{
				ExpSeconds: 3600,
			},
			setupMocks:     func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Invalid Input",
		},
		{
			name: "fail: negative exp_seconds",
			input: ShortenURLRequest{
				URL:        testURL,
				ExpSeconds: -10,
			},
			setupMocks:     func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Invalid Input",
		},
		{
			name: "fail: exceed exp_seconds limit",
			input: ShortenURLRequest{
				URL:        testURL,
				ExpSeconds: 86400 + 1,
			},
			setupMocks:     func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Invalid Input",
		},
		{
			name: "success: omit optional exp_seconds",
			input: ShortenURLRequest{
				URL: testURL,
			},
			setupMocks: func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {
				mockSvc.On("CreateShortenedLink", ctx, testURL, time.Duration(0)).
					Return(testCode, nil).
					Once()
			},
			expectedStatus:  http.StatusCreated,
			successResponse: responseSuccess(testCode),
		},
		{
			name: "fail: service error",
			input: ShortenURLRequest{
				URL:        testURL,
				ExpSeconds: testExpSeconds,
			},
			setupMocks: func(mockSvc *service_mocks.ShortenURL, ctx context.Context) {
				mockSvc.On("CreateShortenedLink", ctx, testURL, testExp).
					Return("", errors.New("fail to store")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErrMsg: "Internal Server Error",
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
			if tc.successResponse != nil {
				assert.Equal(t, tc.expectedErrMsg, "")
				var resp ShortenURLResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, *tc.successResponse, resp)
			}
			if tc.expectedErrMsg != "" {
				assert.Nil(t, tc.successResponse)
				var resp dto.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}

			// Verify all mock expectations are met
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestShortenURL_Redirect(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const testCode = "abc1234"
	const testURL = "https://google.com"

	testCases := []struct {
		name           string
		codeParam      string
		setupMocks     func(ctx context.Context, mockSvc *service_mocks.ShortenURL)
		expectedStatus int
		expectedURL    string
		expectedErrMsg string
	}{
		{
			name:      "success",
			codeParam: testCode,
			setupMocks: func(ctx context.Context, mockSvc *service_mocks.ShortenURL) {
				mockSvc.On("GetURLFromCode", ctx, testCode).
					Return(testURL, nil).
					Once()
			},
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    testURL,
		},
		{
			name:           "fail: empty code",
			codeParam:      "",
			setupMocks:     func(ctx context.Context, mockSvc *service_mocks.ShortenURL) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Invalid Input",
		},
		{
			name:      "fail: URL not found",
			codeParam: testCode,
			setupMocks: func(ctx context.Context, mockSvc *service_mocks.ShortenURL) {
				mockSvc.On("GetURLFromCode", ctx, testCode).
					Return("", service.ErrURLNotFound).
					Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "URL not found",
		},
		{
			name:      "fail: internal server error",
			codeParam: testCode,
			setupMocks: func(ctx context.Context, mockSvc *service_mocks.ShortenURL) {
				mockSvc.On("GetURLFromCode", ctx, testCode).
					Return("", errors.New("any")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErrMsg: "Internal Server Error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)

			endpoint := "/v1/links/redirect/" + tc.codeParam
			ginCtx.Request = httptest.NewRequest(http.MethodGet, endpoint, nil)
			if tc.codeParam != "" {
				ginCtx.Params = gin.Params{{Key: "code", Value: tc.codeParam}}
			}

			mockSvc := service_mocks.NewShortenURL(t)
			tc.setupMocks(ginCtx, mockSvc)
			handler := NewShortenURL(mockSvc)

			handler.Redirect(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedURL != "" {
				assert.Equal(t, tc.expectedURL, rec.Header().Get("Location"))
			}

			if tc.expectedErrMsg != "" {
				var resp dto.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}
