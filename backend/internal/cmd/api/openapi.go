package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/pb33f/libopenapi"
	openapivalidator "github.com/pb33f/libopenapi-validator"

	"goadmin-backend/internal/platform/httperr"
)

type OpenAPIValidator struct {
	validator openapivalidator.Validator
	logger    *slog.Logger
}

func NewOpenAPIValidator(
	oaiPath string,
	logger *slog.Logger,
) (*OpenAPIValidator, error) {
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
		return nil, fmt.Errorf(
			"failed to create validator: %w",
			errors.Join(errs...),
		)
	}

	return &OpenAPIValidator{
		validator: validator,
		logger:    logger,
	}, nil
}

// Middleware returns a middleware that validates incoming requests against the OpenAPI 3+ document.
func (v *OpenAPIValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		valid, validationErrs := v.validator.ValidateHttpRequest(req)
		if !valid {
			var errs []error

			validationErrItems := make([]httperr.ValidationErrorItem, 0)

			for _, err := range validationErrs {
				errs = append(errs, err)

				for _, schemaErr := range err.SchemaValidationErrors {
					validationErrItems = append(
						validationErrItems,
						httperr.ValidationErrorItem{
							Detail:  schemaErr.Reason,
							Pointer: schemaErr.Location,
						},
					)
				}
			}

			v.logger.Error(
				"validation error",
				slog.Any("err", errors.Join(errs...)),
			)

			validationError := httperr.NewValidationError(
				req.URL.Path,
				validationErrItems,
			)

			httperr.JSONError(
				res,
				validationError,
				http.StatusUnprocessableEntity,
			)

			return
		}

		next.ServeHTTP(res, req)
	})
}
