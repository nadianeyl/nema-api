package api

import (
	"net/http"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

func (app *application) listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	var req service.ListCategoriesRequest

	v := validator.New()
	qs := r.URL.Query()
	req.TransactionType = domain.TransactionType(app.httpUtil.ReadString(qs, "transaction_type", ""))
	req.Limit = app.httpUtil.ReadInt(qs, "limit", 10, v)
	req.Page = app.httpUtil.ReadInt(qs, "page", 1, v)
	req.UserID = httputil.ContextGetUserID(r)

	if service.ValidateListCategoriesReq(v, &req); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, metadata, err := app.services.Categories.List(&req)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
		return
	}

	err = app.httpUtil.WriteJSON(w, httputil.StatusRetrieveSuccess, res, &metadata, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}
