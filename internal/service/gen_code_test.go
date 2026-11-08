package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenCode_GenerateCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		expectedLen int
		expectedErr error
	}{
		{
			name:        "success",
			expectedLen: 10,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := NewGenCode()

			code, err := svc.GenerateCode(tc.expectedLen)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedLen, len(code))
		})
	}
}
