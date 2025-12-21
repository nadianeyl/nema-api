package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nadianeyl/nema-api/internal/config"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/jsonlog"
	"github.com/nadianeyl/nema-api/internal/middleware"
	"github.com/nadianeyl/nema-api/internal/service"
)

type application struct {
	config     config.Config
	logger     *jsonlog.Logger
	middleware middleware.Middleware
	httpUtil   httputil.HTTPUtil
	services   service.Services
	wg         sync.WaitGroup
}

func NewApp(cfg config.Config, l *jsonlog.Logger, s service.Services) *application {
	hu := httputil.New(l)

	return &application{
		config:     cfg,
		logger:     l,
		middleware: middleware.New(cfg, l, hu, s.Users),
		httpUtil:   hu,
		services:   s,
	}
}

func (app *application) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		app.logger.LogInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.LogInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.LogInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.Env,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.LogInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
