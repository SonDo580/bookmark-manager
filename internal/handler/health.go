package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/dto"
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
// @Success      200  {object}  dto.HealthCheckResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /health-check [get]
func (h *health) HealthCheck(c *gin.Context) {
	resp, err := h.svc.Check(c)
	if err != nil {
		log.Println(fmt.Errorf("HealthCheck failed: %w", err).Error())

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			dto.ErrorResponse{Message: "Internal Server Error"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}
