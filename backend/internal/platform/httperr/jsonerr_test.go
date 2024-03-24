package httperr

import (
	"net/http"
	"reflect"
	"testing"

	"goadmin-backend/internal/domain"
)

func TestRESTAPIError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		Type     string
		Title    string
		Status   int
		Detail   string
		Instance string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := RESTAPIError{
				Type:     tt.fields.Type,
				Title:    tt.fields.Title,
				Status:   tt.fields.Status,
				Detail:   tt.fields.Detail,
				Instance: tt.fields.Instance,
			}

			if got := e.Error(); got != tt.want {
				t.Errorf("RESTAPIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRESTAPIError(t *testing.T) {
	t.Parallel()

	type args struct {
		instance string
		errType  string
		title    string
		status   int
		detail   string
	}

	tests := []struct {
		name string
		args args
		want *RESTAPIError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewRESTAPIError(tt.args.instance, tt.args.errType, tt.args.title, tt.args.status, tt.args.detail); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRESTAPIError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRESTAPIError_Encode(t *testing.T) {
	t.Parallel()

	type fields struct {
		Type     string
		Title    string
		Status   int
		Detail   string
		Instance string
	}

	type args struct {
		res http.ResponseWriter
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := RESTAPIError{
				Type:     tt.fields.Type,
				Title:    tt.fields.Title,
				Status:   tt.fields.Status,
				Detail:   tt.fields.Detail,
				Instance: tt.fields.Instance,
			}

			if err := e.Encode(tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("RESTAPIError.Encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRESTAPIError_GetInstance(t *testing.T) {
	t.Parallel()

	type fields struct {
		Type     string
		Title    string
		Status   int
		Detail   string
		Instance string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := RESTAPIError{
				Type:     tt.fields.Type,
				Title:    tt.fields.Title,
				Status:   tt.fields.Status,
				Detail:   tt.fields.Detail,
				Instance: tt.fields.Instance,
			}

			if got := e.GetInstance(); got != tt.want {
				t.Errorf("RESTAPIError.GetInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRESTAPIError_SetInstance(t *testing.T) {
	t.Parallel()

	type fields struct {
		Type     string
		Title    string
		Status   int
		Detail   string
		Instance string
	}

	type args struct {
		instance string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &RESTAPIError{
				Type:     tt.fields.Type,
				Title:    tt.fields.Title,
				Status:   tt.fields.Status,
				Detail:   tt.fields.Detail,
				Instance: tt.fields.Instance,
			}

			e.SetInstance(tt.args.instance)
		})
	}
}

func TestWithDomainError(t *testing.T) {
	t.Parallel()

	type args struct {
		err    domain.Error
		status int
	}

	tests := []struct {
		name string
		args args
		want *RESTAPIError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := WithDomainError(tt.args.err, tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDomainError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithStatus(t *testing.T) {
	t.Parallel()

	type args struct {
		status int
	}

	tests := []struct {
		name string
		args args
		want *RESTAPIError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := WithStatus(tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationErrorItem_String(t *testing.T) {
	t.Parallel()

	type fields struct {
		Detail  string
		Pointer string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := ValidationErrorItem{
				Detail:  tt.fields.Detail,
				Pointer: tt.fields.Pointer,
			}

			if got := e.String(); got != tt.want {
				t.Errorf("ValidationErrorItem.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		RESTAPIError RESTAPIError
		Errors       []ValidationErrorItem
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := ValidationError{
				RESTAPIError: tt.fields.RESTAPIError,
				Errors:       tt.fields.Errors,
			}

			if got := e.Error(); got != tt.want {
				t.Errorf("ValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	t.Parallel()

	type args struct {
		instance string
		detail   string
		errs     []ValidationErrorItem
	}

	tests := []struct {
		name string
		args args
		want *ValidationError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewValidationError(tt.args.instance, tt.args.detail, tt.args.errs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_Encode(t *testing.T) {
	t.Parallel()

	type fields struct {
		RESTAPIError RESTAPIError
		Errors       []ValidationErrorItem
	}

	type args struct {
		res http.ResponseWriter
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := ValidationError{
				RESTAPIError: tt.fields.RESTAPIError,
				Errors:       tt.fields.Errors,
			}

			if err := e.Encode(tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("ValidationError.Encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSONError(t *testing.T) {
	t.Parallel()

	type args struct {
		res       http.ResponseWriter
		err       error
		code      int
		endpoints []string
	}

	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			JSONError(tt.args.res, tt.args.err, tt.args.code, tt.args.endpoints...)
		})
	}
}
