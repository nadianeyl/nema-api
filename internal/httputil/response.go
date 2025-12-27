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
	StatusSuccess         = NewStatus(http.StatusOK, "200000", "Success")
	StatusRetrieveSuccess = NewStatus(http.StatusOK, "200001", "Resource retrieved successfully")
	StatusUpdateSuccess   = NewStatus(http.StatusOK, "200002", "Resource updated successfully")
	StatusDeleteSuccess   = NewStatus(http.StatusOK, "200003", "Resource deleted successfully")
	StatusCreated         = NewStatus(http.StatusCreated, "201000", "Resource created successfully")
	StatusAccepted        = NewStatus(http.StatusAccepted, "202000", "Request accepted and is being processed")

	// 4xx
	StatusBadRequest                = NewStatus(http.StatusBadRequest, "400000", "Bad request")
	StatusAuthRequired              = NewStatus(http.StatusUnauthorized, "401000", "Authentication required")
	StatusInvalidAuthCredentials    = NewStatus(http.StatusUnauthorized, "401001", "Invalid authentication credentials")
	StatusInvalidOrMissingAuthToken = NewStatus(http.StatusUnauthorized, "401002", "Invalid or missing authentication token")
	StatusUserNotAllowed            = NewStatus(http.StatusForbidden, "403000", "Access denied")
	StatusInactiveUserAccount       = NewStatus(http.StatusForbidden, "403001", "User account not activated")
	StatusNotFound                  = NewStatus(http.StatusNotFound, "404000", "Resource not found")
	StatusEditConflict              = NewStatus(http.StatusConflict, "409000", "Edit conflict detected")
	StatusValidationFailed          = NewStatus(http.StatusUnprocessableEntity, "422000", "Validation failed")

	// 5xx
	StatusInternalServerError = NewStatus(http.StatusInternalServerError, "500000", "Server encountered a problem")
)
