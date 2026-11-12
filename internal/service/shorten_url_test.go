package service

import (
	"context"
	"errors"
	"testing"
	"time"

	repository_mocks "github.com/SonDo580/bookmark-manager/internal/repository/mocks"
	service_mocks "github.com/SonDo580/bookmark-manager/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL_CreateShortenedLink(t *testing.T) {
	t.Parallel()

	const testURL = "https://facebook.com"
	const testCode = "abc1234"
	const testExp = 1 * time.Hour

	errGenCode := errors.New("fail to generate code")
	errStore := errors.New("fail to store")

	testCases := []struct {
		name         string
		setupMocks   func(ctx context.Context, mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode)
		expectedCode string
		expectedErr  error
	}{
		{
			name: "success",
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode) {
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return(testCode, nil).Once()
				mockRepo.On("StoreURL", ctx, testCode, testURL, testExp).Return(true, nil).Once()
			},
			expectedCode: testCode,
			expectedErr:  nil,
		},
		{
			name: "success: retry on collision",
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode) {
				// Collide
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return(testCode, nil).Once()
				mockRepo.On("StoreURL", ctx, testCode, testURL, testExp).Return(false, nil).Once()

				// Retry
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return(testCode, nil).Once()
				mockRepo.On("StoreURL", ctx, testCode, testURL, testExp).Return(true, nil).Once()
			},
			expectedCode: testCode,
			expectedErr:  nil,
		},
		{
			name: "fail: code generation error",
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode) {
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return("", errGenCode).Once()
			},
			expectedCode: "",
			expectedErr:  errGenCode,
		},
		{
			name: "fail: storage error",
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode) {
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return(testCode, nil).Once()
				mockRepo.On("StoreURL", ctx, testCode, testURL, testExp).Return(false, errStore).Once()
			},
			expectedCode: "",
			expectedErr:  errStore,
		},
		{
			name: "fail: exhaust max retries",
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode) {
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return(testCode, nil).Times(MaxRetries)
				mockRepo.On("StoreURL", ctx, testCode, testURL, testExp).Return(false, nil).Times(MaxRetries)
			},
			expectedCode: "",
			expectedErr:  ErrMaxRetriesExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockGenCode := service_mocks.NewGenCode(t)
			mockRepo := repository_mocks.NewURLStorage(t)
			tc.setupMocks(ctx, mockRepo, mockGenCode)

			svc := NewShortenURL(mockRepo, mockGenCode)

			code, err := svc.CreateShortenedLink(ctx, testURL, testExp)

			assert.Equal(t, tc.expectedCode, code)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockGenCode.AssertExpectations(t)
		})
	}
}

func TestShortenURL_GetURLFromCode(t *testing.T) {
	t.Parallel()

	const testURL = "https://facebook.com"
	const testCode = "abc1234"

	errGetURL := errors.New("any")

	testCases := []struct {
		name        string
		code        string
		setupMocks  func(ctx context.Context, mockRepo *repository_mocks.URLStorage)
		expectedURL string
		expectedErr error
	}{
		{
			name: "success",
			code: testCode,
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.URLStorage) {
				mockRepo.On("GetURL", ctx, testCode).
					Return(testURL, nil).Once()
			},
			expectedURL: testURL,
			expectedErr: nil,
		},
		{
			name: "fail: URL not found",
			code: testCode,
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.URLStorage) {
				mockRepo.On("GetURL", ctx, testCode).
					Return("", ErrURLNotFound).Once()
			},
			expectedURL: "",
			expectedErr: ErrURLNotFound,
		},
		{
			name: "fail: other error",
			code: testCode,
			setupMocks: func(ctx context.Context, mockRepo *repository_mocks.URLStorage) {
				mockRepo.On("GetURL", ctx, testCode).
					Return("", errGetURL).Once()
			},
			expectedURL: "",
			expectedErr: errGetURL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockGenCode := service_mocks.NewGenCode(t)
			mockRepo := repository_mocks.NewURLStorage(t)
			tc.setupMocks(ctx, mockRepo)

			svc := NewShortenURL(mockRepo, mockGenCode)

			url, err := svc.GetURLFromCode(ctx, tc.code)

			assert.Equal(t, tc.expectedURL, url)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockGenCode.AssertExpectations(t)
		})
	}
}
