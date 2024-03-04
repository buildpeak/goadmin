package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/pb33f/libopenapi"
	openapivalidator "github.com/pb33f/libopenapi-validator"
)

type multiErrors []error

func (m multiErrors) Error() string {
	errStrs := make([]string, len(m))

	for i, err := range m {
		errStrs[i] = err.Error()
	}

	return strings.Join(errStrs, "; ")
}

type OpenAPIValidator struct {
	validator openapivalidator.Validator
}

func NewOpenAPIValidator(oaiPath string) (*OpenAPIValidator, error) {
	if oaiPath == "" {
		oaiPath = "openapi.yaml"
	}

	spec, err := os.ReadFile(oaiPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read openapi.yaml: %w", err)
	}

	document, err := libopenapi.NewDocument(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to parse openapi.yaml: %w", err)
	}

	validator, errs := openapivalidator.NewValidator(document)
	if len(errs) > 0 {
		errStrs := make([]string, len(errs))

		for i, err := range errs {
			errStrs[i] = err.Error()
		}

		return nil, fmt.Errorf("failed to create validator: %w", multiErrors(errs))
	}

	return &OpenAPIValidator{
		validator: validator,
	}, nil
}

// Middleware returns a middleware that validates incoming requests against the OpenAPI 3+ document.
func (v *OpenAPIValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		valid, errs := v.validator.ValidateHttpRequest(req)
		if !valid {
			errStrs := make([]string, len(errs))

			for i, err := range errs {
				errStrs[i] = err.Error()
			}

			http.Error(res, "request validation failed: "+strings.Join(errStrs, "; "), http.StatusBadRequest)

			return
		}

		next.ServeHTTP(res, req)
	})
}
