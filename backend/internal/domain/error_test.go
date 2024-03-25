package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"goadmin-backend/internal/platform/httperr"
)

func TestBaseError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		Type    string
		Message string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "base error",
			fields: fields{
				Type:    "base",
				Message: "base error",
			},
			want: "base: base error",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := baseError{
				Type:    tt.fields.Type,
				Message: tt.fields.Message,
			}

			if got := e.Error(); got != tt.want {
				t.Errorf("BaseError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseError_As(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		createError func() error
	}{
		{
			name: "wrap BaseError",
			createError: func() error {
				return fmt.Errorf("new base error: %w", &baseError{
					Type:    "test",
					Message: "test error",
				})
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.createError()

			var be *baseError
			if !errors.As(got, &be) {
				t.Errorf("errors.As check failed; %v %v", got, be)
			}
		})
	}
}

func TestBaseError_ToRESTAPIError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		baseErr baseError
		want    *httperr.RESTAPIError
	}{
		{
			name: "Success",
			baseErr: baseError{
				Type:    "test",
				Message: "test error",
			},
			want: &httperr.RESTAPIError{
				Type:   "/errors/test",
				Title:  "TEST",
				Detail: "test error",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.baseErr.ToRESTAPIError(); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("WithDomainError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceNotFoundError(t *testing.T) {
	t.Parallel()

	type args struct {
		resource  string
		condition string
	}

	tests := []struct {
		name string
		args args
		want *ResourceNotFoundError
	}{
		{
			name: "Good",
			args: args{
				resource:  "User",
				condition: "id=1",
			},
			want: &ResourceNotFoundError{
				baseError: baseError{
					Type:    "user-not-found",
					Message: "User with condition id=1 not found",
				},
				Resource:  "User",
				Condition: "id=1",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewResourceNotFoundError(tt.args.resource, tt.args.condition); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceNotFoundError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		BaseError baseError
		Resource  string
		Condition string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "resource not found",
			fields: fields{
				BaseError: baseError{
					Type:    "resource_not_found",
					Message: "resource with condition id=1 not found",
				},
				Resource:  "resource",
				Condition: "id=1",
			},
			want: "resource_not_found: resource with condition id=1 not found",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := ResourceNotFoundError{
				baseError: tt.fields.BaseError,
				Resource:  tt.fields.Resource,
				Condition: tt.fields.Condition,
			}

			b, _ := json.Marshal(e)
			t.Log(string(b))

			if got := e.Error(); got != tt.want {
				t.Errorf("ResourceNotFoundError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceNotFoundError_As(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "wrap BaseError",
			err:  NewResourceNotFoundError("Role", "id=1"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var r *ResourceNotFoundError
			if !errors.As(tt.err, &r) {
				t.Errorf("errors.As check failed; %v %v", tt.err, r)
			}

			var er Error
			if !errors.As(tt.err, &er) {
				t.Errorf("errors.As check failed; %v %v", tt.err, er)
			}
		})
	}
}

func TestResourceExistsError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		BaseError baseError
		Resource  string
		Conflict  string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "resource exists",
			fields: fields{
				BaseError: baseError{
					Type:    "user_already_exists",
					Message: "User with id=1 already exists",
				}, Resource: "resource",
				Conflict: "id=1",
			},
			want: "user_already_exists: User with id=1 already exists",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := ResourceExistsError{
				baseError: tt.fields.BaseError,
				Resource:  tt.fields.Resource,
				Conflict:  tt.fields.Conflict,
			}

			if got := e.Error(); got != tt.want {
				t.Errorf("ResourceExistsError.Detail() = %v, want %v", got, tt.want)
			}
		})
	}
}
