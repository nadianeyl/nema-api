package httputil

import (
	"encoding/json"
	"net/http"
)

func (h *HTTPUtil) WriteJSON(w http.ResponseWriter, status Status, data any, headers http.Header) error {
	res := Response{
		Status:  status.Code,
		Message: status.Message,
		Data:    data,
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return err
	}

	jsonBytes = append(jsonBytes, '\n')

	for key, val := range headers {
		w.Header()[key] = val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status.HTTPStatus)
	w.Write(jsonBytes)

	return nil
}
