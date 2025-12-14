package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	app.logger.LogError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status Status, data any) {
	err := app.writeJSON(w, status, data, nil)
	if err != nil {
		app.logError(r, err)

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	app.errorResponse(w, r, StatusInternalServerError, nil)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, StatusNotFound, nil)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	data := envelope{"error": err.Error()}
	app.errorResponse(w, r, StatusBadRequest, data)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, StatusValidationFailed, errors)
}
