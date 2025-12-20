package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nadianeyl/nema-api/internal/config"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/jsonlog"
	"github.com/nadianeyl/nema-api/internal/middleware"
	"github.com/nadianeyl/nema-api/internal/service"
)

type Application struct {
	Config     config.Config
	Logger     *jsonlog.Logger
	Middleware middleware.Middleware
	HTTPUtil   httputil.HTTPUtil
	Services   service.Services
}

func NewApp(cfg config.Config, logger *jsonlog.Logger, services service.Services) *Application {
	hu := httputil.New(logger)

	return &Application{
		Config:     cfg,
		Logger:     logger,
		Middleware: middleware.New(cfg, logger, hu),
		HTTPUtil:   hu,
		Services:   services,
	}
}

func (app *Application) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.Logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		app.Logger.LogInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.Logger.LogInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.Config.Env,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.Logger.LogInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
