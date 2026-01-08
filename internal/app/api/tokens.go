package api

import (
	"errors"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

func (app *application) createAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req service.CreateAuthTokenRequest

	err := app.httpUtil.ReadJSON(w, r, &req)
	if err != nil {
		app.httpUtil.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	service.ValidateEmail(v, req.Email)
	service.ValidatePasswordPlaintext(v, req.Password)
	if !v.Valid() {
		app.httpUtil.FailedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := app.services.Tokens.CreateAuthToken(&req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound), errors.Is(err, domain.ErrInvalidAuthCredentials):
			app.httpUtil.InvalidCredentialsResponse(w, r)
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
