package api

import (
	"fmt"
	"net/http"

	"github.com/SonDo580/bookmark-manager/docs"
	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/internal/handler"
	"github.com/SonDo580/bookmark-manager/internal/repository"
	"github.com/SonDo580/bookmark-manager/internal/service"
	"github.com/SonDo580/bookmark-manager/pkg/jwtutils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "github.com/SonDo580/bookmark-manager/docs"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app          *gin.Engine
	cfg          *config.Config
	redisClient  *redis.Client
	db           *gorm.DB
	jwtGenerator jwtutils.JWTGenerator
	jwtValidator jwtutils.JWTValidator
}

func NewEngine(
	cfg *config.Config,
	redisClient *redis.Client,
	db *gorm.DB,
	jwtGenerator jwtutils.JWTGenerator,
	jwtValidator jwtutils.JWTValidator,
) Engine {
	e := &engine{
		app:          gin.Default(),
		cfg:          cfg,
		redisClient:  redisClient,
		db:           db,
		jwtGenerator: jwtGenerator,
		jwtValidator: jwtValidator,
	}

	e.initRoutes()

	return e
}

type handlers struct {
	healthHandler handler.Health
	linkHandler   handler.ShortenURL
	userHandler   handler.UserHandler
}

func (e *engine) initHandlers() *handlers {
	healthSvc := service.NewHealth(e.cfg, e.redisClient)
	healthHandler := handler.NewHealth(healthSvc)

	urlStorageRepo := repository.NewURLStorage(e.redisClient)
	genCodeSvc := service.NewGenCode()
	shortenURLSvc := service.NewShortenURL(urlStorageRepo, genCodeSvc)
	shortenURLHandler := handler.NewShortenURL(shortenURLSvc)

	userRepo := repository.NewUserRepository(e.db)
	userSvc := service.NewUserService(userRepo, e.jwtGenerator, e.cfg)
	userHandler := handler.NewUserHandler(userSvc)

	return &handlers{
		healthHandler: healthHandler,
		linkHandler:   shortenURLHandler,
		userHandler:   userHandler,
	}
}

func (e *engine) initRoutes() {
	allHandlers := e.initHandlers()

	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	e.app.GET("/health-check", allHandlers.healthHandler.HealthCheck)

	v1Routes := e.app.Group("/v1")
	{
		linkRoutes := v1Routes.Group("/links")
		{
			linkRoutes.POST("/shorten", allHandlers.linkHandler.CreateShortenedLink)
			linkRoutes.GET("/redirect/:code", allHandlers.linkHandler.Redirect)
		}

		userRoutes := v1Routes.Group("/users")
		{
			userRoutes.POST("/register", allHandlers.userHandler.Register)
			userRoutes.POST("/login", allHandlers.userHandler.Login)
		}
	}
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%d", e.cfg.Port))
}

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}
