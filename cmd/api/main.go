package main

import (
	"os"

	"github.com/nadianeyl/nema-api/internal/app/api"
	"github.com/nadianeyl/nema-api/internal/config"
	"github.com/nadianeyl/nema-api/internal/db"
	"github.com/nadianeyl/nema-api/internal/jsonlog"

	_ "github.com/lib/pq"
)

func main() {
	var cfg config.Config
	config.Init(&cfg)

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db.Init(cfg, logger)

	app := api.New(cfg, logger)

	err := app.Serve()
	if err != nil {
		logger.LogFatal(err, nil)
	}
}
