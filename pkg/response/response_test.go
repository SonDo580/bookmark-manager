package response

import (
	"errors"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type ShortenURLRequest struct {
	URL        string `json:"url" binding:"required,url"`
	ExpSeconds int    `json:"exp_seconds" binding:"gte=0,lte=86400"`
}

func TestInputFieldError(t *testing.T) {
	t.Parallel()

	validate := validator.New()
	validate.SetTagName("binding")

	testCases := []struct {
		name            string
		inputErr        error
		expectedDetails []string
	}{
		{
			name:            "non-validation error",
			inputErr:        errors.New("any"),
			expectedDetails: nil,
		},
		{
			name: "single validation error",
			inputErr: validate.Struct(ShortenURLRequest{
				URL:        "", // fail 'required'
				ExpSeconds: 3600,
			}),
			expectedDetails: []string{
				"URL is invalid (required)",
			},
		},
		{
			name: "multiple validation errors",
			inputErr: validate.Struct(ShortenURLRequest{
				URL:        "invalid-url", // fail 'url'
				ExpSeconds: 86400 + 1,     // fail 'lte'
			}),
			expectedDetails: []string{
				"URL is invalid (url)",
				"ExpSeconds is invalid (lte)",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := InputFieldError(tc.inputErr)

			expectedMessage := dto.MsgInvalidInput
			assert.Equal(t, expectedMessage, resp.Message)

			if tc.expectedDetails == nil {
				assert.Nil(t, resp.Details)
			} else {
				details, ok := resp.Details.([]string)
				assert.True(t, ok, "expected Details to be of type []string")
				assert.ElementsMatch(t, tc.expectedDetails, details)
			}
		})
	}
}
