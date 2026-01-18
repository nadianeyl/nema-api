package api

import (
	"errors"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var req service.RegisterUserRequest

	err := app.httpUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.httpUtil.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if service.ValidateRegisterUserReq(v, &req); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.services.Users.Register(r.Context(), &req, &app.wg)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateEmail):
			v.AddError("email", "email already exists")
			app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		default:
			app.httpUtil.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.httpUtil.WriteJSON(w, httputil.StatusAccepted, res, nil, nil)
	if err != nil {
		app.httpUtil.ServerErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req service.ActivateUserRequest

	err := app.httpUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.httpUtil.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if service.ValidateTokenPlaintext(v, req.TokenPlaintext); !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.services.Users.Activate(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidOrExpiredToken):
			v.AddError("token", domain.ErrInvalidOrExpiredToken.Error())
			app.httpUtil.FailedValidationResponse(w, r, v.Errors)
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
