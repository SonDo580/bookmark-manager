package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func InitMockRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	mockServer := miniredis.RunT(t)
	mockClient := redis.NewClient(&redis.Options{
		Addr: mockServer.Addr(),
	})
	return mockServer, mockClient
}
