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
	cfg := config.Init()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	_ = db.Init(cfg, logger)

	app := api.NewApp(cfg, logger)
	err := app.Serve()
	if err != nil {
		logger.LogFatal(err, nil)
	}
}
