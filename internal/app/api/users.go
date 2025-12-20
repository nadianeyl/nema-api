package api

import (
	"errors"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/repository"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var req service.RegisterUserRequest

	err := app.HTTPUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.HTTPUtil.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if service.ValidateRegisterUserReq(v, &req); !v.Valid() {
		app.HTTPUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.Services.Users.Register(&req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			v.AddError("email", "email already exists")
			app.HTTPUtil.FailedValidationResponse(w, r, v.Errors)
		default:
			app.HTTPUtil.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.HTTPUtil.WriteJSON(w, httputil.StatusCreated, res, nil)
	if err != nil {
		app.HTTPUtil.ServerErrorResponse(w, r, err)
	}
}
