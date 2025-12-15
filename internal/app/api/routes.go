package api

import (
	"net/http"
)

func (app *Application) routes() http.Handler {
	router := http.NewServeMux()
	m := app.Middleware

	router.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler)

	return m.RecoverPanic(m.EnableCORS(m.RequestLogger(router)))
}
