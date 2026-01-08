package api

import (
	"errors"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

func (app *application) createBudgetItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.httpUtil.ReadIDParam(r)
	if err != nil {
		app.httpUtil.NotFoundResponse(w, r)
		return
	}

	var req service.CreateBudgetItemRequest

	err = app.httpUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.httpUtil.BadRequestResponse(w, r, err)
		return
	}

	req.BudgetID = *id

	v := validator.New()
	if service.ValidateCreateBudgetItemReq(v, &req); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.services.BudgetItems.Create(&req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound), errors.Is(err, domain.ErrDuplicateRecord):
			app.httpUtil.BadRequestResponse(w, r, err)
		default:
			app.httpUtil.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.httpUtil.WriteJSON(w, httputil.StatusCreated, res, nil, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}
