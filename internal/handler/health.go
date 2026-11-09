package handler

import (
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/model"
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
// @Description  Check status of application and dependencies.
// @Tags         health
// @Produce      json
// @Success      200  {object}  model.HealthCheckResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /health-check [get]
func (h *health) HealthCheck(c *gin.Context) {
	response, err := h.svc.Check(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			model.ErrorResponse{Error: "Internal Server Error"},
		)
		return
	}

	c.JSON(http.StatusOK, response)
}
