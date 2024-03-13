package httperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

var (
	// ErrUnauthorized is returned when a user is not authorized to perform an action.
	ErrUnauthorized = &RESTAPIError{
		Type:   "/errors/unauthorized",
		Title:  "Unauthorized",
		Status: http.StatusUnauthorized,
		Detail: "You are not authorized to perform this action",
	}

	// ErrForbidden is returned when a user is forbidden from performing an action.
	ErrForbidden = &RESTAPIError{
		Type:   "/errors/forbidden",
		Title:  "Forbidden",
		Status: http.StatusForbidden,
		Detail: "You are forbidden from performing this action",
	}

	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = &RESTAPIError{
		Type:   "/errors/not_found",
		Title:  "Not Found",
		Status: http.StatusNotFound,
		Detail: "The requested resource was not found",
	}

	// ErrInternalServerError is returned when an internal server error occurs.
	ErrInternalServerError = &RESTAPIError{
		Type:   "/errors/internal_server_error",
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
		Detail: "An internal server error occurred",
	}
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

func NewRESTAPIError(
	instance string,
	errType string,
	title string,
	status int,
	detail string,
) *RESTAPIError {
	return &RESTAPIError{
		Type:     errType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
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

func lookupError(code int) *RESTAPIError {
	errMap := map[int]*RESTAPIError{
		http.StatusUnauthorized: ErrUnauthorized,
		http.StatusForbidden:    ErrForbidden,
		http.StatusNotFound:     ErrNotFound,
	}

	if err, ok := errMap[code]; ok {
		return err
	}

	return ErrInternalServerError
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
	case code != http.StatusInternalServerError:
		errEnc = encoder.Encode(lookupError(code))
	default:
		slog.Error("error not recognized", slog.Any("err", err))

		errEnc = encoder.Encode(ErrInternalServerError)
	}

	if errEnc != nil {
		slog.Error("error encoding json", slog.Any("err", errEnc))
	}
}
