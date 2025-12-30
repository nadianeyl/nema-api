package main

import (
	"os"

	_ "github.com/lib/pq"

	"github.com/nadianeyl/nema-api/internal/app/api"
	"github.com/nadianeyl/nema-api/internal/config"
	"github.com/nadianeyl/nema-api/internal/db"
	"github.com/nadianeyl/nema-api/internal/jsonlog"
	"github.com/nadianeyl/nema-api/internal/mailer"
	"github.com/nadianeyl/nema-api/internal/repository"
	"github.com/nadianeyl/nema-api/internal/service"
)

func main() {
	cfg := config.Init()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	mailer := mailer.New(cfg, logger)
	db := db.Init(cfg, logger)
	defer db.Close()

	repositories := repository.NewRepositories(db)
	txProvider := repository.NewTxProvider(db)
	services := service.NewServices(repositories, txProvider, mailer, logger)
	app := api.NewApp(cfg, logger, services)

	err := app.Serve()
	if err != nil {
		logger.LogFatal(err, nil)
	}
}
