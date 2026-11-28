package handler

import (
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/constants"
	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/SonDo580/bookmark-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Health interface {
	HealthCheck(c *gin.Context)
}

type health struct {
	cfg *config.Config
	svc service.Health
}

func NewHealth(cfg *config.Config, svc service.Health) Health {
	return &health{cfg: cfg, svc: svc}
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
	err := h.svc.Check(c)
	if err != nil {
		log.Error().Err(err).Msg("HealthCheck failed")

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			dto.ErrorResponse{Message: constants.MsgInternalServerError},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dto.HealthCheckResponse{
			Message:     constants.MsgOK,
			ServiceName: h.cfg.ServiceName,
			InstanceID:  h.cfg.InstanceID,
		},
	)
}
