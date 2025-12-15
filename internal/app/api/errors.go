package api

import "net/http"

func (app *Application) logError(r *http.Request, err error) {
	app.Logger.LogError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status Status, data any) {
	err := app.writeJSON(w, status, data, nil)
	if err != nil {
		app.logError(r, err)

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	app.errorResponse(w, r, StatusInternalServerError, nil)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, StatusNotFound, nil)
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	data := Envelope{"error": err.Error()}
	app.errorResponse(w, r, StatusBadRequest, data)
}

func (app *Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, StatusValidationFailed, errors)
}
