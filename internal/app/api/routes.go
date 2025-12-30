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

	router.Handle("GET /api/v1/accounts", m.RequireActivatedUser(app.listAccountsHandler))
	router.Handle("POST /api/v1/accounts", m.RequireActivatedUser(app.addAccountHandler))
	router.Handle("PATCH /api/v1/accounts/{id}", m.RequireActivatedUser(app.updateAccountHandler))
	router.Handle("DELETE /api/v1/accounts/{id}", m.RequireActivatedUser(app.deleteAccountHandler))

	router.Handle("GET /api/v1/categories", m.RequireActivatedUser(app.listCategoriesHandler))

	router.Handle("POST /api/v1/transactions", m.RequireActivatedUser(app.addTransactionHandler))

	return m.RecoverPanic(m.EnableCORS(m.RequestLogger(m.Authenticate(router))))
}
