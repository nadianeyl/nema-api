package api

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	router := http.NewServeMux()
	m := app.middleware

	router.Handle("GET /api/v1/healthcheck", http.HandlerFunc(app.healthcheckHandler))

	router.Handle("POST /api/v1/users", http.HandlerFunc(app.registerUserHandler))
	router.Handle("PUT /api/v1/users/activated", http.HandlerFunc(app.activateUserHandler))
	router.Handle("POST /api/v1/users/authentication", http.HandlerFunc(app.createAuthTokenHandler))

	return m.RecoverPanic(m.EnableCORS(m.RequestLogger(m.Authenticate(router))))
}
