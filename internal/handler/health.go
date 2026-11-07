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

// HealthCheck godoc
// @Summary      Health check
// @Description  Check status of application and dependencies
// @Tags         health
// @Produce      application/json
// @Success      200  {object}  model.HealthCheckResponse
// @Router       /health-check [get]
func (h *health) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, h.svc.Check())
}
