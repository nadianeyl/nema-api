package api

import (
	"errors"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

func (app *application) listTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	var req service.ListTransactionsRequest

	v := validator.New()
	qs := r.URL.Query()
	req.Type = domain.TransactionType(app.httpUtil.ReadString(qs, "type", ""))
	req.CategoryID = app.httpUtil.ReadString(qs, "category_id", "")
	req.AccountID = app.httpUtil.ReadString(qs, "account_id", "")
	req.StartDate = app.httpUtil.ReadString(qs, "start_date", "")
	req.EndDate = app.httpUtil.ReadString(qs, "end_date", "")
	req.Title = app.httpUtil.ReadString(qs, "title", "")
	req.Limit = app.httpUtil.ReadInt(qs, "limit", 10, v)
	req.Page = app.httpUtil.ReadInt(qs, "page", 1, v)
	req.Sort = app.httpUtil.ReadString(qs, "sort", "-date")
	req.SortSafelist = []string{"date", "amount", "-date", "-amount"}
	req.UserID = httputil.ContextGetUserID(r)

	if service.ValidateListTransactionsReq(v, &req); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, metadata, err := app.services.Transactions.List(&req)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
		return
	}

	err = app.httpUtil.WriteJSON(w, httputil.StatusRetrieveSuccess, res, &metadata, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}

func (app *application) addTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var req service.AddTransactionRequest

	err := app.httpUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.httpUtil.BadRequestResponse(w, r, err)
		return
	}

	req.UserID = httputil.ContextGetUserID(r)

	v := validator.New()
	if service.ValidateAddTransactionReq(v, &req); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.services.Transactions.Add(&req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInputValue):
			app.httpUtil.BadRequestResponse(w, r, err)
		case errors.Is(err, domain.ErrUserNotAllowed):
			app.httpUtil.UserNotAllowedResponse(w, r)
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

func (app *application) getTransactionDetailHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.httpUtil.ReadIDParam(r)
	if err != nil {
		app.httpUtil.NotFoundResponse(w, r)
		return
	}

	var req service.GetTransactionDetailRequest

	req.ID = *id
	req.UserID = httputil.ContextGetUserID(r)

	res, err := app.services.Transactions.GetDetailByID(&req)
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

	err = app.httpUtil.WriteJSON(w, httputil.StatusRetrieveSuccess, res, nil, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}

func (app *application) updateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.httpUtil.ReadIDParam(r)
	if err != nil {
		app.httpUtil.NotFoundResponse(w, r)
		return
	}

	var req service.UpdateTransactionRequest

	err = app.httpUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.httpUtil.BadRequestResponse(w, r, err)
		return
	}

	req.ID = *id
	req.UserID = httputil.ContextGetUserID(r)

	v := validator.New()
	if service.ValidateUpdateTransactionReq(v, &req); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.services.Transactions.Update(&req, v)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			app.httpUtil.NotFoundResponse(w, r)
		case errors.Is(err, domain.ErrInvalidInputValue):
			app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		case errors.Is(err, domain.ErrUserNotAllowed):
			app.httpUtil.UserNotAllowedResponse(w, r)
		case errors.Is(err, domain.ErrEditConflict):
			app.httpUtil.EditConflictResponse(w, r)
		default:
			app.httpUtil.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.httpUtil.WriteJSON(w, httputil.StatusUpdateSuccess, res, nil, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}

func (app *application) deleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.httpUtil.ReadIDParam(r)
	if err != nil {
		app.httpUtil.NotFoundResponse(w, r)
		return
	}

	var req service.DeleteTransactionRequest

	req.ID = *id
	req.UserID = httputil.ContextGetUserID(r)

	err = app.services.Transactions.Delete(&req)
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
