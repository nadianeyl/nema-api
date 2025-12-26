package api

import (
	"errors"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

func (app *application) listAccountsHandler(w http.ResponseWriter, r *http.Request) {
	var req service.ListAccountsRequest

	v := validator.New()
	qs := r.URL.Query()
	req.Class = domain.AccountClass(app.httpUtil.ReadString(qs, "class", ""))
	req.Limit = app.httpUtil.ReadInt(qs, "limit", 10, v)
	req.Page = app.httpUtil.ReadInt(qs, "page", 1, v)
	req.UserID = httputil.ContextGetUserID(r)

	if service.ValidateListAccountsReq(v, &req); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, metadata, err := app.services.Accounts.List(&req)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
		return
	}

	err = app.httpUtil.WriteJSON(w, httputil.StatusSuccess, res, &metadata, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}

func (app *application) addAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req service.AddAccountRequest

	err := app.httpUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.httpUtil.BadRequestResponse(w, r, err)
		return
	}

	req.UserID = httputil.ContextGetUserID(r)

	v := validator.New()
	if service.ValidateAddAccountReq(v, &req); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.services.Accounts.Add(&req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInputValue):
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

func (app *application) updateAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.httpUtil.ReadIDParam(r)
	if err != nil {
		app.httpUtil.NotFoundResponse(w, r)
		return
	}

	var req service.UpdateAccountRequest

	err = app.httpUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.httpUtil.BadRequestResponse(w, r, err)
		return
	}

	req.ID = *id
	req.UserID = httputil.ContextGetUserID(r)

	v := validator.New()
	if service.ValidateUpdateAccountReq(v, &req); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.services.Accounts.Update(&req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			app.httpUtil.NotFoundResponse(w, r)
		case errors.Is(err, domain.ErrUserNotAllowed):
			app.httpUtil.UserNotAllowedResponse(w, r)
		case errors.Is(err, domain.ErrEditConflict):
			app.httpUtil.EditConflictResponse(w, r)
		default:
			app.httpUtil.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.httpUtil.WriteJSON(w, httputil.StatusSuccess, res, nil, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}

func (app *application) deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.httpUtil.ReadIDParam(r)
	if err != nil {
		app.httpUtil.NotFoundResponse(w, r)
		return
	}

	var req service.DeleteAccountRequest

	req.ID = *id
	req.UserID = httputil.ContextGetUserID(r)

	err = app.services.Accounts.Delete(&req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			app.httpUtil.NotFoundResponse(w, r)
		case errors.Is(err, domain.ErrUserNotAllowed):
			app.httpUtil.UserNotAllowedResponse(w, r)
		default:
			app.httpUtil.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.httpUtil.WriteJSON(w, httputil.StatusDeleteSuccess, nil, nil, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}
