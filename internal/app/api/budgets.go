package api

import (
	"errors"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

func (app *application) createBudgetHandler(w http.ResponseWriter, r *http.Request) {
	var req service.CreateBudgetRequest

	err := app.httpUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.httpUtil.BadRequestResponse(w, r, err)
		return
	}

	req.UserID = httputil.ContextGetUserID(r)

	v := validator.New()
	if service.ValidateCreateBudgetReq(v, &req); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.services.Budgets.Create(&req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateRecord):
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

func (app *application) getBudgetDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.httpUtil.ReadIDParam(r)
	if err != nil {
		app.httpUtil.NotFoundResponse(w, r)
		return
	}

	var req service.GetBudgetDetailRequest

	req.ID = *id
	req.UserID = httputil.ContextGetUserID(r)

	res, err := app.services.Budgets.GetByID(&req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			app.httpUtil.NotFoundResponse(w, r)
		default:
			app.httpUtil.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.httpUtil.WriteJSON(w, httputil.StatusRetrieveSuccess, res, nil, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}
