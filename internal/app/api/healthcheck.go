package api

import (
	"net/http"

	"github.com/nadianeyl/nema-api/internal/config"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := Envelope{
		"service_status": "available",
		"system_info": map[string]string{
			"environment": app.Config.Env,
			"version":     config.Version,
		},
	}

	err := app.writeJSON(w, StatusSuccess, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
