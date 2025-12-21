package api

import (
	"errors"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/repository"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

func (app *Application) createAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req service.CreateAuthTokenRequest

	err := app.HTTPUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.HTTPUtil.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	service.ValidateEmail(v, req.Email)
	service.ValidatePasswordPlaintext(v, req.Password)
	if !v.Valid() {
		app.HTTPUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.Services.Tokens.CreateAuthToken(&req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound), errors.Is(err, repository.ErrInvalidAuthCredentials):
			app.HTTPUtil.InvalidCredentialsResponse(w, r)
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
