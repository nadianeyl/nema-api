package api

import (
	"errors"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

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
