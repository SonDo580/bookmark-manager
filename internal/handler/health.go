package handler

import (
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/service"
	"github.com/gin-gonic/gin"
)

type Health interface {
	HealthCheck(c *gin.Context)
}

type health struct {
	svc service.Health
}

func NewHealth(svc service.Health) Health {
	return &health{svc: svc}
}

func (h *health) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, h.svc.Check())
}
