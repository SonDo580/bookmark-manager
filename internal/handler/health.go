package handler

import (
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/service"
	"github.com/gin-gonic/gin"
)

type Health interface {
	HealthCheck(c *gin.Context)
}

type health struct {
	svc service.Health
}

func NewHealth(cfg *config.Config) Health {
	healthSvc := service.NewHealth(cfg)
	return &health{svc: healthSvc}
}

func (h *health) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, h.svc.Check())
}
