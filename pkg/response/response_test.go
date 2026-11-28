package response

import (
	"errors"
	"testing"

	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type exampleInput struct {
	URL string `json:"url" binding:"required,url"`
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
			inputErr: validate.Struct(exampleInput{
				URL: "",
			}),
			expectedDetails: []string{
				"URL is invalid (required)",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := InputFieldError(tc.inputErr)

			expectedMessage := constants.MsgInvalidInput
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
