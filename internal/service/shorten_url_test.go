package service

import (
	"errors"
	"testing"
	"time"

	repository_mocks "github.com/SonDo580/bookmark-manager/internal/repository/mocks"
	service_mocks "github.com/SonDo580/bookmark-manager/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		setupMocks   func(*repository_mocks.URLStorage, *service_mocks.GenCode)
		expectedCode string
		expectedErr  error
	}{
		{
			name: "success",
			setupMocks: func(mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode) {
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return(testCode, nil).Once()
				mockRepo.On("StoreURL", mock.Anything, testCode, testURL, testExp).Return(true, nil).Once()
			},
			expectedCode: testCode,
			expectedErr:  nil,
		},
		{
			name: "success: retry on collision",
			setupMocks: func(mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode) {
				// Collide
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return(testCode, nil).Once()
				mockRepo.On("StoreURL", mock.Anything, testCode, testURL, testExp).Return(false, nil).Once()

				// Retry
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return(testCode, nil).Once()
				mockRepo.On("StoreURL", mock.Anything, testCode, testURL, testExp).Return(true, nil).Once()
			},
			expectedCode: testCode,
			expectedErr:  nil,
		},
		{
			name: "fail: code generation error",
			setupMocks: func(mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode) {
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return("", errGenCode).Once()
			},
			expectedCode: "",
			expectedErr:  errGenCode,
		},
		{
			name: "fail: storage error",
			setupMocks: func(mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode) {
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return(testCode, nil).Once()
				mockRepo.On("StoreURL", mock.Anything, testCode, testURL, testExp).Return(false, errStore).Once()
			},
			expectedCode: "",
			expectedErr:  errStore,
		},
		{
			name: "fail: exhaust max retries",
			setupMocks: func(mockRepo *repository_mocks.URLStorage, mockGenCode *service_mocks.GenCode) {
				mockGenCode.On("GenerateCode", DefaultCodeLength).Return(testCode, nil).Times(MaxRetries)
				mockRepo.On("StoreURL", mock.Anything, testCode, testURL, testExp).Return(false, nil).Times(MaxRetries)
			},
			expectedCode: "",
			expectedErr:  ErrMaxRetriesExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockGenCode := service_mocks.NewGenCode(t)
			mockRepo := repository_mocks.NewURLStorage(t)
			tc.setupMocks(mockRepo, mockGenCode)

			svc := NewShortenURL(mockRepo, mockGenCode)

			ctx := t.Context()
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
