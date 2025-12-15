package api

import "net/http"

func (app *Application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler)

	return app.recoverPanic(app.enableCORS(app.requestLogger(router)))
}
