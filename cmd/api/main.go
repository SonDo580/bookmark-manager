package main

import (
	"github.com/SonDo580/bookmark-manager/internal/api"
	"github.com/SonDo580/bookmark-manager/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
