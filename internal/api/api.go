package api

import (
	"fmt"
	"net/http"

	"github.com/SonDo580/bookmark-manager/docs"
	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/handler"
	"github.com/SonDo580/bookmark-manager/internal/repository"
	"github.com/SonDo580/bookmark-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/SonDo580/bookmark-manager/docs"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app         *gin.Engine
	cfg         *config.Config
	redisClient *redis.Client
}

func NewEngine(cfg *config.Config, redisClient *redis.Client) Engine {
	e := &engine{app: gin.Default(), cfg: cfg, redisClient: redisClient}
	e.initRoutes()

	return e
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%d", e.cfg.Port))
}

func (e *engine) initRoutes() {
	healthSvc := service.NewHealth(e.cfg, e.redisClient)
	healthHandler := handler.NewHealth(healthSvc)

	urlStorageRepo := repository.NewURLStorage(e.redisClient)
	genCodeSvc := service.NewGenCode()
	shortenURLSvc := service.NewShortenURL(urlStorageRepo, genCodeSvc)
	shortenURLHandler := handler.NewShortenURL(shortenURLSvc)

	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	e.app.GET("/health-check", healthHandler.HealthCheck)

	e.app.POST("/v1/links/shorten", shortenURLHandler.CreateShortenedLink)
	e.app.GET("/v1/links/redirect/:code", shortenURLHandler.Redirect)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}
