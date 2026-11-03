package handler

import (
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/gin-gonic/gin"
)

type Bookmark interface {
	HealthCheck(c *gin.Context)
}

type bookmark struct {
	cfg config.Config
}

func NewBookMark(cfg config.Config) Bookmark {
	return &bookmark{cfg: cfg}
}

func (h *bookmark) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":      "OK",
		"service_name": h.cfg.ServiceName,
		"instance_id":  h.cfg.InstanceID,
	})
}
