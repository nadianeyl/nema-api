package httputil

import (
	"net/http"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type Envelope map[string]any

type Response struct {
	Status   string           `json:"status"`
	Message  string           `json:"message"`
	Data     any              `json:"data,omitempty"`
	Metadata *domain.Metadata `json:"metadata,omitempty"`
}

type Status struct {
	HTTPStatus int
	Code       string
	Message    string
}

func NewStatus(httpStatus int, code, message string) Status {
	return Status{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

var (
	// 2xx
	StatusSuccess       = NewStatus(http.StatusOK, "200000", "Success")
	StatusDeleteSuccess = NewStatus(http.StatusOK, "200001", "Resource successfully deleted")
	StatusCreated       = NewStatus(http.StatusCreated, "201000", "Resource successfully created")
	StatusAccepted      = NewStatus(http.StatusAccepted, "202000", "Request accepted and is being processed")

	// 4xx
	StatusBadRequest                = NewStatus(http.StatusBadRequest, "400000", "Bad request")
	StatusInvalidAuthCredentials    = NewStatus(http.StatusUnauthorized, "401000", "Invalid authentication credentials")
	StatusInvalidOrMissingAuthToken = NewStatus(http.StatusUnauthorized, "401001", "Invalid or missing authentication token")
	StatusAuthRequired              = NewStatus(http.StatusUnauthorized, "401002", "You must be authenticated to access this resource")
	StatusUserNotAllowed            = NewStatus(http.StatusForbidden, "403000", "You are not allowed to access this resource")
	StatusInactiveUserAccount       = NewStatus(http.StatusForbidden, "403001", "Your user account must be activated to access this resource")
	StatusNotFound                  = NewStatus(http.StatusNotFound, "404000", "Resource could not be found")
	StatusEditConflict              = NewStatus(http.StatusConflict, "409000", "Unable to update the record due to an edit conflict, please try again")
	StatusValidationFailed          = NewStatus(http.StatusUnprocessableEntity, "422000", "Validation failed")

	// 5xx
	StatusInternalServerError = NewStatus(http.StatusInternalServerError, "500000", "Server encountered a problem and could not process your request")
)
