package api

import (
	"net/http"

	"github.com/nadianeyl/nema-api/internal/config"
	"github.com/nadianeyl/nema-api/internal/httputil"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := httputil.Envelope{
		"service_status": "available",
		"system_info": map[string]string{
			"environment": app.config.Env,
			"version":     config.Version,
		},
	}

	err := app.httpUtil.WriteJSON(w, httputil.StatusSuccess, data, nil, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}
