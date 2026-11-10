package service

import (
	"context"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=Health --filename=health.go --outpkg=service_mocks
type Health interface {
	Check(ctx context.Context) (dto.HealthCheckResponse, error)
}

type health struct {
	cfg         *config.Config
	redisClient *redis.Client
}

func NewHealth(cfg *config.Config, redisClient *redis.Client) Health {
	return &health{cfg: cfg, redisClient: redisClient}
}

func (s *health) Check(ctx context.Context) (dto.HealthCheckResponse, error) {
	// Check Redis connection
	if err := s.redisClient.Ping(ctx).Err(); err != nil {
		return dto.HealthCheckResponse{}, err
	}

	return dto.HealthCheckResponse{
		Message:     dto.HealthStatusOK,
		ServiceName: s.cfg.ServiceName,
		InstanceID:  s.cfg.InstanceID,
	}, nil
}
