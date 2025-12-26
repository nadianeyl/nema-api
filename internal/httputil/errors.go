package httputil

import "net/http"

func (h *HTTPUtil) logError(r *http.Request, err error) {
	h.logger.LogError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (h *HTTPUtil) errorResponse(w http.ResponseWriter, r *http.Request, status Status, data any) {
	err := h.WriteJSON(w, status, data, nil, nil)
	if err != nil {
		h.logError(r, err)

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *HTTPUtil) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)

	h.errorResponse(w, r, StatusInternalServerError, nil)
}

func (h *HTTPUtil) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	h.errorResponse(w, r, StatusNotFound, nil)
}

func (h *HTTPUtil) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	data := Envelope{"error": err.Error()}
	h.errorResponse(w, r, StatusBadRequest, data)
}

func (h *HTTPUtil) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	h.errorResponse(w, r, StatusValidationFailed, errors)
}

func (h *HTTPUtil) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	h.errorResponse(w, r, StatusEditConflict, nil)
}

func (h *HTTPUtil) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	h.errorResponse(w, r, StatusInvalidAuthCredentials, nil)
}

func (h *HTTPUtil) InvalidAuthTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	h.errorResponse(w, r, StatusInvalidOrMissingAuthToken, nil)
}

func (h *HTTPUtil) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	h.errorResponse(w, r, StatusAuthRequired, nil)
}

func (h *HTTPUtil) UserNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	h.errorResponse(w, r, StatusUserNotAllowed, nil)
}

func (h *HTTPUtil) InactiveUserAccountResponse(w http.ResponseWriter, r *http.Request) {
	h.errorResponse(w, r, StatusInactiveUserAccount, nil)
}
