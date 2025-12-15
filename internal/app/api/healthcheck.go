package api

import (
	"net/http"

	"github.com/nadianeyl/nema-api/internal/config"
	"github.com/nadianeyl/nema-api/internal/httputil"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := httputil.Envelope{
		"service_status": "available",
		"system_info": map[string]string{
			"environment": app.Config.Env,
			"version":     config.Version,
		},
	}

	err := app.HTTPUtil.WriteJSON(w, httputil.StatusSuccess, data, nil)
	if err != nil {
		app.HTTPUtil.ServerErrorResponse(w, r, err)
	}
}
