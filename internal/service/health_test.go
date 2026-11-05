package service

import (
	"testing"
	"uuid"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		cfg  *config.Config
	}{
		{
			name: "success",
			cfg: &config.Config{
				ServiceName: "bookmark_service",
				InstanceID:  uuid.New().String(),
				Port:        8080,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			healthSvc := NewHealth(tc.cfg)
			actualResponse := healthSvc.Check()

			expectedResponse := model.HealthCheckResponse{
				Message:     model.HealthStatusOK,
				ServiceName: tc.cfg.ServiceName,
				InstanceID:  tc.cfg.InstanceID,
			}

			assert.Equal(t, expectedResponse, actualResponse)
		})
	}
}
