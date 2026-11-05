package handler

import (
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/gin-gonic/gin"
)

type Health interface {
	HealthCheck(c *gin.Context)
}

type health struct {
	cfg *config.Config
}

func NewHealth(cfg *config.Config) Health {
	return &health{cfg: cfg}
}

func (h *health) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":      "OK",
		"service_name": h.cfg.ServiceName,
		"instance_id":  h.cfg.InstanceID,
	})
}
