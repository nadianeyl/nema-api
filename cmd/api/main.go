package main

import (
	"os"

	"github.com/nadianeyl/nema-api/internal/app/api"
	"github.com/nadianeyl/nema-api/internal/config"
	"github.com/nadianeyl/nema-api/internal/db"
	"github.com/nadianeyl/nema-api/internal/jsonlog"
	"github.com/nadianeyl/nema-api/internal/repository"
	"github.com/nadianeyl/nema-api/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Init()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db := db.Init(cfg, logger)
	defer db.Close()

	repositories := repository.NewRepositories(db, logger)
	services := service.NewServices(repositories)
	app := api.NewApp(cfg, logger, services)

	err := app.Serve()
	if err != nil {
		logger.LogFatal(err, nil)
	}
}
