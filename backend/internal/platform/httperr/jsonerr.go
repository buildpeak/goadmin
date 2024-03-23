package httperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"goadmin-backend/internal/domain"
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
		Type:   "/errors/not-found",
		Title:  "Not Found",
		Status: http.StatusNotFound,
		Detail: "The requested resource was not found",
	}

	// ErrBadRequest is returned when a bad request is made.
	ErrBadRequest = &RESTAPIError{
		Type:   "/errors/bad-request",
		Title:  "Bad Request",
		Status: http.StatusBadRequest,
		Detail: "The request was invalid or cannot be served",
	}

	// ErrInternalServerError is returned when an internal server error occurs.
	ErrInternalServerError = &RESTAPIError{
		Type:   "/errors/internal-server-error",
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
		Detail: "An internal server error occurred",
	}
)

type jsonError interface {
	error
	Encode(res http.ResponseWriter) error
	GetInstance() string
	SetInstance(instance string)
}

// RESTAPIError is a struct that represents the base error type for the API.
// https://www.rfc-editor.org/rfc/rfc9457.html#name-the-problem-details-json-ob
type RESTAPIError struct {
	// Type is a URI reference that identifies the problem type.
	Type string `json:"type"`

	// Title is a short, human-readable summary of the problem type.
	Title string `json:"title"`

	// Status is the HTTP status code generated by the origin server for this occurrence of the problem.
	Status int `json:"status"`

	// Detail is a human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail"`

	// Instance is a URI reference that identifies the specific occurrence of the problem.
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

func (e RESTAPIError) Encode(res http.ResponseWriter) error {
	return json.NewEncoder(res).Encode(e) //nolint:wrapcheck // no need to wrap
}

func (e RESTAPIError) GetInstance() string {
	return e.Instance
}

func (e *RESTAPIError) SetInstance(instance string) {
	e.Instance = instance
}

func WithDomainError(err domain.Error, status int) *RESTAPIError {
	return &RESTAPIError{
		Type:   "/errors/" + err.GetType(),
		Title:  strings.ToTitle(strings.ReplaceAll(err.GetType(), "_", " ")),
		Status: status,
		Detail: err.GetMessage(),
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
			Type:     "/errors/validation-error",
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
		http.StatusBadRequest:   ErrBadRequest,
	}

	if err, ok := errMap[code]; ok {
		return err
	}

	return ErrInternalServerError
}

func JSONError(res http.ResponseWriter, err error, code int, endpoints ...string) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)

	var jsnErr jsonError

	var restapiError *RESTAPIError

	var validationError *ValidationError

	var domainError domain.Error

	instance := strings.Join(endpoints, ",")

	switch {
	case errors.As(err, &restapiError):
		jsnErr = restapiError
	case errors.As(err, &validationError):
		jsnErr = validationError
	case errors.As(err, &domainError):
		jsnErr = WithDomainError(domainError, code)
	case code != http.StatusInternalServerError:
		jsnErr = lookupError(code)
	default:
		slog.Error("error not recognized", slog.Any("err", err))

		jsnErr = ErrInternalServerError
	}

	// update instance
	if jsnErr.GetInstance() == "" {
		jsnErr.SetInstance(instance)
	}

	if err := jsnErr.Encode(res); err != nil {
		slog.Error("error encoding json", slog.Any("err", err))
	}
}
