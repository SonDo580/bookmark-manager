package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/model"
	"github.com/SonDo580/bookmark-manager/internal/service"
	service_mocks "github.com/SonDo580/bookmark-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const testUsername = "shadow"
	const testEmail = "shadow@x.com"
	const testDisplayName = "Shadow"
	const testPassword = "examplePass"

	validReq := dto.RegisterRequest{
		Username:    testUsername,
		Email:       testEmail,
		DisplayName: testDisplayName,
		Password:    testPassword,
	}

	expectedUser := &model.User{
		Username:    testUsername,
		Email:       testEmail,
		DisplayName: testDisplayName,
	}

	testCases := []struct {
		name            string
		input           any
		setupMocks      func(ctx context.Context, mockSvc *service_mocks.UserService)
		expectedStatus  int
		successResponse *dto.RegisterResponse
		expectedErrMsg  string
	}{
		{
			name:  "success",
			input: validReq,
			setupMocks: func(ctx context.Context, mockSvc *service_mocks.UserService) {
				mockSvc.On("Register", ctx, validReq).
					Return(expectedUser, nil).
					Once()
			},
			expectedStatus: http.StatusCreated,
			successResponse: &dto.RegisterResponse{
				Message: "Register user successfully!",
				Data:    expectedUser,
			},
		},
		{
			name:           "fail: missing required fields",
			input:          dto.RegisterRequest{},
			setupMocks:     func(ctx context.Context, mockSvc *service_mocks.UserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: constants.MsgInvalidInput,
		},
		{
			name: "fail: password length < 8",
			input: dto.RegisterRequest{
				Username:    testUsername,
				Email:       testEmail,
				DisplayName: testDisplayName,
				Password:    strings.Repeat("a", 7),
			},
			setupMocks:     func(ctx context.Context, mockSvc *service_mocks.UserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: constants.MsgInvalidInput,
		},
		{
			name: "fail: invalid email",
			input: dto.RegisterRequest{
				Username:    testUsername,
				Email:       "not-an-email",
				DisplayName: testDisplayName,
				Password:    testPassword,
			},
			setupMocks:     func(ctx context.Context, mockSvc *service_mocks.UserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: constants.MsgInvalidInput,
		},
		// More invalid input:
		// - invalid JSON.
		// - wrong field types.
		// - other constraint violations.

		{
			name:  "fail: duplicate username",
			input: validReq,
			setupMocks: func(ctx context.Context, mockSvc *service_mocks.UserService) {
				mockSvc.On("Register", ctx, validReq).
					Return(nil, service.ErrDuplicateUsername).
					Once()
			},
			expectedStatus: http.StatusConflict,
			expectedErrMsg: constants.MsgUsernameAlreadyExists,
		},
		{
			name:  "fail: duplicate email",
			input: validReq,
			setupMocks: func(ctx context.Context, mockSvc *service_mocks.UserService) {
				mockSvc.On("Register", ctx, validReq).
					Return(nil, service.ErrDuplicateEmail).
					Once()
			},
			expectedStatus: http.StatusConflict,
			expectedErrMsg: constants.MsgEmailAlreadyExists,
		},
		{
			name:  "fail: service error",
			input: validReq,
			setupMocks: func(ctx context.Context, mockSvc *service_mocks.UserService) {
				mockSvc.On("Register", ctx, validReq).
					Return(nil, errors.New("any")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErrMsg: constants.MsgInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var inputBytes []byte
			var err error
			switch input := tc.input.(type) {
			case string:
				inputBytes = []byte(input)
			default:
				inputBytes, err = json.Marshal(input)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(inputBytes))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request = req

			mockSvc := service_mocks.NewUserService(t)
			tc.setupMocks(ginCtx, mockSvc)

			handler := NewUserHandler(mockSvc)

			handler.Register(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.successResponse != nil {
				assert.Equal(t, tc.expectedErrMsg, "")
				var resp dto.RegisterResponse
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

			mockSvc.AssertExpectations(t)
		})
	}
}
