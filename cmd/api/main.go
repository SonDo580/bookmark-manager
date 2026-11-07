package main

import (
	"github.com/SonDo580/bookmark-manager/internal/api"
	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/SonDo580/bookmark-manager/pkg/redis"
)

// @title           Bookmark Manager API
// @version         1.0.0
// @description		API for managing bookmark
// @BasePath  		/
func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg, redisClient)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
