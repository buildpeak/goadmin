package domain

import (
	"fmt"
	"strings"

	"goadmin-backend/internal/platform/httperr"
)

// Error categories:
// - Not Found - READ
// - Already Exists (Creation Conflict) - CREATE
// - Conflict (Update Conflict) - UPDATE
// - Forbidden - DELETE

// DomainErrorer is an interface that represents a domain error.
type Error interface {
	error
	GetType() string
	GetMessage() string
}

// baseError is a struct that represents the base error type for the API.
type baseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e baseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e baseError) GetType() string {
	return e.Type
}

func (e baseError) GetMessage() string {
	return e.Message
}

// ToRESTAPIError converts a domain error to a RESTAPIError.
func (e baseError) ToRESTAPIError() *httperr.RESTAPIError {
	return &httperr.RESTAPIError{
		Type:   "/errors/" + e.GetType(),
		Title:  strings.ToTitle(strings.ReplaceAll(e.GetType(), "-", " ")),
		Detail: e.GetMessage(),
		// Status: httperr.GetHTTPStatusCode(e.GetType()),
	}
}

type ResourceNotFoundError struct {
	baseError
	Resource  string `json:"resource"`
	Condition string `json:"condition"`
}

func NewResourceNotFoundError(resource, condition string) *ResourceNotFoundError {
	return &ResourceNotFoundError{
		baseError: baseError{
			Type: strings.ToLower(resource) + "-not-found",
			Message: fmt.Sprintf(
				"%s with condition %s not found",
				resource,
				condition,
			),
		},
		Resource:  resource,
		Condition: condition,
	}
}

type ResourceExistsError struct {
	baseError
	Resource string `json:"resource"`
	Conflict string `json:"conflict"`
}

func NewResourceExistsError(resource, conflict string) *ResourceExistsError {
	return &ResourceExistsError{
		baseError: baseError{
			Type:    strings.ToLower(resource) + "-already-exists",
			Message: fmt.Sprintf("%s with %s already exists", resource, conflict),
		},
		Resource: resource,
	}
}
