package service

import (
	"context"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=Health --filename=health.go --outpkg=service_mocks
type Health interface {
	Check(ctx context.Context) error
}

type health struct {
	redisClient *redis.Client
}

func NewHealth(redisClient *redis.Client) Health {
	return &health{redisClient: redisClient}
}

func (s *health) Check(ctx context.Context) error {
	// Check Redis connection
	return s.redisClient.Ping(ctx).Err()

	// TODO: check DB connection
}
