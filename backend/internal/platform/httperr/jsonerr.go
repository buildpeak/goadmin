package httperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// RESTAPIError is a struct that represents the base error type for the API.
// https://www.rfc-editor.org/rfc/rfc9457.html#name-the-problem-details-json-ob
type RESTAPIError struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func (e RESTAPIError) Error() string {
	return fmt.Sprintf("%s - %s: %s", e.Title, e.Detail, e.Type)
}

type ValidationErrorItem struct {
	Detail  string `json:"detail"`
	Pointer string `json:"pointer"`
}

func (e ValidationErrorItem) String() string {
	return fmt.Sprintf("%s - %s", e.Detail, e.Pointer)
}

type ValidationError struct {
	RESTAPIError
	Errors []ValidationErrorItem `json:"errors"`
}

func (e ValidationError) Error() string {
	errStrings := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		errStrings[i] = err.String()
	}

	return fmt.Sprintf(
		"%s - %s: %s, errors: %s",
		e.Title,
		e.Detail,
		e.Type,
		strings.Join(errStrings, ", "),
	)
}

func NewValidationError(
	instance string,
	errs []ValidationErrorItem,
) *ValidationError {
	return &ValidationError{
		RESTAPIError: RESTAPIError{
			Type:     "/errors/validation_error",
			Title:    "Validation Error",
			Status:   http.StatusUnprocessableEntity,
			Detail:   "One or more validation errors occurred",
			Instance: instance,
		},
		Errors: errs,
	}
}

func JSONError(res http.ResponseWriter, err error, code int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)

	encoder := json.NewEncoder(res)

	var restapiError *RESTAPIError

	var validationError *ValidationError

	var errEnc error

	switch {
	case errors.As(err, &restapiError):
		errEnc = encoder.Encode(restapiError)
	case errors.As(err, &validationError):
		errEnc = encoder.Encode(validationError)
	default:
		slog.Error("error not recognized", slog.Any("err", err))

		errEnc = encoder.Encode(RESTAPIError{
			Type:   "/errors/internal_server_error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "An internal server error occurred",
		})
	}

	if errEnc != nil {
		slog.Error("error encoding json", slog.Any("err", errEnc))
	}
}
