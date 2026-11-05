package service

import (
	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/model"
)

type Health interface {
	Check() model.HealthCheckResp
}

type health struct {
	cfg *config.Config
}

func NewHealth(cfg *config.Config) Health {
	return &health{cfg: cfg}
}

func (s *health) Check() model.HealthCheckResp {
	return model.HealthCheckResp{
		Message:     "OK",
		ServiceName: s.cfg.ServiceName,
		InstanceID:  s.cfg.InstanceID,
	}
}
