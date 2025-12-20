package api

import (
	"net/http"
)

func (app *Application) routes() http.Handler {
	router := http.NewServeMux()
	m := app.Middleware

	router.Handle("GET /api/v1/healthcheck", http.HandlerFunc(app.healthcheckHandler))

	router.Handle("POST /api/v1/users", http.HandlerFunc(app.registerUserHandler))

	return m.RecoverPanic(m.EnableCORS(m.RequestLogger(router)))
}
