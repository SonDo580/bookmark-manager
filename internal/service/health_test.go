package service

import (
	"testing"
	"uuid"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	redisPkg "github.com/SonDo580/bookmark-manager/pkg/redis"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	defaultCfg := &config.Config{
		ServiceName: "bookmark_service",
		InstanceID:  uuid.New().String(),
		Port:        8080,
	}

	testCases := []struct {
		name          string
		cfg           *config.Config
		setupMocks    func(mockRedisServer *miniredis.Miniredis)
		errIsExpected bool
	}{
		{
			name:          "success",
			cfg:           defaultCfg,
			setupMocks:    func(mockRedisServer *miniredis.Miniredis) {},
			errIsExpected: false,
		},
		{
			name: "fail: redis down",
			cfg:  defaultCfg,
			setupMocks: func(mockRedisServer *miniredis.Miniredis) {
				// Shut down the Redis server
				mockRedisServer.Close()
			},
			errIsExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRedisServer, mockRedisClient := redisPkg.InitMockRedis(t)
			defer mockRedisClient.Close()
			tc.setupMocks(mockRedisServer)

			healthSvc := NewHealth(tc.cfg, mockRedisClient)

			actualResponse, err := healthSvc.Check(t.Context())

			if tc.errIsExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				expectedResponse := dto.HealthCheckResponse{
					Message:     dto.MsgOK,
					ServiceName: tc.cfg.ServiceName,
					InstanceID:  tc.cfg.InstanceID,
				}
				assert.Equal(t, expectedResponse, actualResponse)
			}
		})
	}
}
