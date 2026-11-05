package api

import (
	"fmt"
	"net/http"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/handler"
	"github.com/SonDo580/bookmark-manager/internal/service"
	"github.com/gin-gonic/gin"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *config.Config
}

func NewEngine(cfg *config.Config) Engine {
	e := &engine{app: gin.Default(), cfg: cfg}
	e.initRoutes()

	return e
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%d", e.cfg.Port))
}

func (e *engine) initRoutes() {
	healthSvc := service.NewHealth(e.cfg)
	healthHandler := handler.NewHealth(healthSvc)

	e.app.GET("/health-check", healthHandler.HealthCheck)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}
