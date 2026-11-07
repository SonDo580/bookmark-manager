package service

import (
	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/model"
)

//go:generate mockery --name=Health --filename=health.go --outpkg=mocks
type Health interface {
	Check() model.HealthCheckResponse
}

type health struct {
	cfg *config.Config
}

func NewHealth(cfg *config.Config) Health {
	return &health{cfg: cfg}
}

func (s *health) Check() model.HealthCheckResponse {
	return model.HealthCheckResponse{
		Message:     model.HealthStatusOK,
		ServiceName: s.cfg.ServiceName,
		InstanceID:  s.cfg.InstanceID,
	}
}
