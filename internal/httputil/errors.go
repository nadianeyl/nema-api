package httputil

import "net/http"

func (h *HTTPUtil) logError(r *http.Request, err error) {
	h.Logger.LogError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (h *HTTPUtil) errorResponse(w http.ResponseWriter, r *http.Request, status Status, data any) {
	err := h.WriteJSON(w, status, data, nil)
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
