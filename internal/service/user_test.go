package service

import (
	"context"
	"errors"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/model"
	"github.com/SonDo580/bookmark-manager/internal/repository"
	repository_mocks "github.com/SonDo580/bookmark-manager/internal/repository/mocks"
	jwtutils_mocks "github.com/SonDo580/bookmark-manager/pkg/jwtutils/mocks"
	"github.com/SonDo580/bookmark-manager/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_Register(t *testing.T) {
	t.Parallel()

	const (
		testUsername    = "shadow"
		testEmail       = "shadow@x.com"
		testDisplayName = "Shadow"
		testPassword    = "examplePass"
	)

	input := dto.RegisterRequest{
		Username:    testUsername,
		Email:       testEmail,
		DisplayName: testDisplayName,
		Password:    testPassword,
	}
	mockUser := mock.AnythingOfType("*model.User")
	errCreateUser := errors.New("fail to create user")

	testCases := []struct {
		name               string
		input              dto.RegisterRequest
		setupMocks         func(ctx context.Context, mockRepo *repository_mocks.UserRepository)
		shouldValidateUser bool
		expectedErr        error
	}{
		{
			name:  "success",
			input: input,
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.UserRepository) {
				mockRepo.On("CreateUser", ctx, mockUser).
					Return(nil).
					Once()
			},
			shouldValidateUser: true,
			expectedErr:        nil,
		},
		{
			name:  "fail: duplicate username",
			input: input,
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.UserRepository) {
				mockRepo.On("CreateUser", ctx, mockUser).
					Return(repository.ErrDuplicateUsername).
					Once()
			},
			expectedErr: ErrDuplicateUsername,
		},
		{
			name:  "fail: duplicate email",
			input: input,
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.UserRepository) {
				mockRepo.On("CreateUser", ctx, mockUser).
					Return(repository.ErrDuplicateEmail).
					Once()
			},
			expectedErr: ErrDuplicateEmail,
		},
		{
			name:  "fail: repository error",
			input: input,
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.UserRepository) {
				mockRepo.On("CreateUser", ctx, mockUser).
					Return(errCreateUser).
					Once()
			},
			expectedErr: errCreateUser,
		},
	}

	validateUser := func(t *testing.T, result *model.User, input dto.RegisterRequest) {
		t.Helper()

		require.NotNil(t, result)

		assert.Equal(t, input.Username, result.Username)
		assert.Equal(t, input.Email, result.Email)
		assert.Equal(t, input.DisplayName, result.DisplayName)

		assert.NotEmpty(t, result.Password)
		assert.NotEqual(t, input.Password, result.Password)
		ok, err := utils.CompareHashAndPassword(result.Password, input.Password)
		require.NoError(t, err)
		assert.True(t, ok)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockRepo := repository_mocks.NewUserRepository(t)
			mockJWTGenerator := jwtutils_mocks.NewJWTGenerator(t)
			tc.setupMocks(ctx, mockRepo)

			svc := NewUserService(mockRepo, mockJWTGenerator, &config.Config{})

			result, err := svc.Register(ctx, tc.input)

			if tc.expectedErr != nil {
				assert.Nil(t, result)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				if tc.shouldValidateUser {
					validateUser(t, result, tc.input)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
