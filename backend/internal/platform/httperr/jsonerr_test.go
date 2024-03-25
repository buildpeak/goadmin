package httperr

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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
	}{
		{
			name: "Success",
			fields: fields{
				Type:     "Type",
				Title:    "Title",
				Status:   400,
				Detail:   "Detail",
				Instance: "Instance",
			},
			want: "Title - Detail: Type",
		},
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
		{
			name: "Success",
			args: args{
				instance: "instance",
				errType:  "errType",
				title:    "title",
				status:   400,
				detail:   "detail",
			},
			want: &RESTAPIError{
				Type:     "errType",
				Title:    "title",
				Status:   400,
				Detail:   "detail",
				Instance: "instance",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewRESTAPIError(
				tt.args.instance, tt.args.errType, tt.args.title, tt.args.status, tt.args.detail,
			); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewRESTAPIError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRESTAPIError_WriteJSON(t *testing.T) {
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
		{
			name: "Success",
			fields: fields{
				Type:     "Type",
				Title:    "Title",
				Status:   400,
				Detail:   "Detail",
				Instance: "Instance",
			},
			args: args{
				res: httptest.NewRecorder(),
			},
		},
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

			if err := e.WriteJSON(tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf(
					"RESTAPIError.WriteJSON() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
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
		{
			name: "Success",
			fields: fields{
				Type:     "Type",
				Title:    "Title",
				Status:   400,
				Detail:   "Detail",
				Instance: "Instance",
			},
			want: "Instance",
		},
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
				t.Errorf(
					"RESTAPIError.GetInstance() = %v, want %v",
					got,
					tt.want,
				)
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
		want   string
	}{
		{
			name: "Success",
			fields: fields{
				Type:     "Type",
				Title:    "Title",
				Status:   400,
				Detail:   "Detail",
				Instance: "Instance",
			},
			args: args{
				instance: "NewInstance",
			},
			want: "NewInstance",
		},
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

			if got := e.Instance; got != tt.want {
				t.Errorf(
					"RESTAPIError.SetInstance() = %v, want %v",
					got,
					tt.want,
				)
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
		{
			name: "Bad Request",
			args: args{
				status: 400,
			},
			want: &RESTAPIError{
				Type:   "/errors/bad-request",
				Title:  "Bad Request",
				Status: 400,
				Detail: "The request was invalid or cannot be served",
			},
		},
		{
			name: "Unauthorized",
			args: args{
				status: 401,
			},
			want: &RESTAPIError{
				Type:   "/errors/unauthorized",
				Title:  "Unauthorized",
				Status: 401,
				Detail: "You are not authorized to perform this action",
			},
		},
		{
			name: "Forbidden",
			args: args{
				status: 403,
			},
			want: &RESTAPIError{
				Type:   "/errors/forbidden",
				Title:  "Forbidden",
				Status: 403,
				Detail: "You are not allowed to access this resource",
			},
		},
		{
			name: "Not Found",
			args: args{
				status: 404,
			},
			want: &RESTAPIError{
				Type:   "/errors/not-found",
				Title:  "Not Found",
				Status: 404,
				Detail: "The requested resource was not found",
			},
		},
		{
			name: "Conflict",
			args: args{
				status: 409,
			},
			want: &RESTAPIError{
				Type:   "/errors/conflict",
				Title:  "Conflict",
				Status: 409,
				Detail: "A conflict occurred while processing the request",
			},
		},
		{
			name: "Internal Server Error",
			args: args{
				status: 500,
			},
			want: &RESTAPIError{
				Type:   "/errors/internal-server-error",
				Title:  "Internal Server Error",
				Status: 500,
				Detail: "An internal server error occurred",
			},
		},
		{
			name: "Unknown Error",
			args: args{
				status: 444,
			},
			want: &RESTAPIError{
				Type:   "/errors/unknown-error",
				Title:  "Unknown Error",
				Status: 444,
				Detail: "An unknown error occurred",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := WithStatus(tt.args.status); !reflect.DeepEqual(
				got,
				tt.want,
			) {
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
		{
			name: "Success",
			fields: fields{
				Detail:  "Detail",
				Pointer: "Pointer",
			},
			want: "Detail - Pointer",
		},
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
				t.Errorf(
					"ValidationErrorItem.String() = %v, want %v",
					got,
					tt.want,
				)
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
		{
			name: "Success",
			fields: fields{
				RESTAPIError: RESTAPIError{
					Type:   "/errors/validation-error",
					Title:  "Validation Error",
					Status: http.StatusUnprocessableEntity,
					Detail: "Detail",
				},
				Errors: []ValidationErrorItem{
					{
						Detail:  "Detail",
						Pointer: "Pointer",
					},
				},
			},
			want: "Validation Error - Detail: /errors/validation-error, errors: Detail - Pointer",
		},
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
		{
			name: "Success",
			args: args{
				instance: "instance",
				detail:   "detail",
				errs: []ValidationErrorItem{
					{
						Detail:  "Detail",
						Pointer: "Pointer",
					},
				},
			},
			want: &ValidationError{
				RESTAPIError: RESTAPIError{
					Type:     "/errors/validation-error",
					Title:    "Validation Error",
					Status:   http.StatusUnprocessableEntity,
					Detail:   "detail",
					Instance: "instance",
				},
				Errors: []ValidationErrorItem{
					{
						Detail:  "Detail",
						Pointer: "Pointer",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewValidationError(tt.args.instance, tt.args.detail, tt.args.errs); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_WriteJSON(t *testing.T) {
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
		{
			name: "Success",
			fields: fields{
				RESTAPIError: RESTAPIError{
					Type:     "/errors/validation-error",
					Title:    "Validation Error",
					Status:   http.StatusUnprocessableEntity,
					Detail:   "Detail",
					Instance: "Instance",
				},
				Errors: []ValidationErrorItem{
					{
						Detail:  "Detail",
						Pointer: "Pointer",
					},
				},
			},
			args: args{
				res: httptest.NewRecorder(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := ValidationError{
				RESTAPIError: tt.fields.RESTAPIError,
				Errors:       tt.fields.Errors,
			}

			if err := e.WriteJSON(tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf(
					"ValidationError.WriteJSON() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
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

	selfRefErr := &castableError{}
	selfRefErr.Self = selfRefErr

	tests := []struct {
		name string
		args args
	}{
		{
			name: "REST API Error",
			args: args{
				res:  httptest.NewRecorder(),
				err:  WithStatus(400),
				code: 400,
				endpoints: []string{
					"/users",
				},
			},
		},
		{
			name: "Validation Error",
			args: args{
				res: httptest.NewRecorder(),
				err: NewValidationError(
					"instance",
					"detail",
					[]ValidationErrorItem{
						{
							Detail:  "Detail",
							Pointer: "Pointer",
						},
					},
				),
				code: 422,
				endpoints: []string{
					"/users",
				},
			},
		},
		{
			name: "Castable Error",
			args: args{
				res:  httptest.NewRecorder(),
				err:  castableError{},
				code: 409,
				endpoints: []string{
					"/users",
				},
			},
		},
		{
			name: "400 Bad Request",
			args: args{
				res:  httptest.NewRecorder(),
				err:  errors.New("I don't want to tell"),
				code: 400,
				endpoints: []string{
					"/users",
				},
			},
		},
		{
			name: "500 Internal Server Error",
			args: args{
				res:  httptest.NewRecorder(),
				err:  errors.New("I don't want to tell"),
				code: 500,
				endpoints: []string{
					"/users",
				},
			},
		},
		{
			name: "encoding json error",
			args: args{
				res:  httptest.NewRecorder(),
				err:  selfRefErr,
				code: 500,
				endpoints: []string{
					"/users",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			JSONError(
				tt.args.res,
				tt.args.err,
				tt.args.code,
				tt.args.endpoints...,
			)

			t.Logf(
				"Response: %d %+v",
				//nolint:forcetypeassert // no need to check
				tt.args.res.(*httptest.ResponseRecorder).Code,
				//nolint:forcetypeassert // no need to check
				tt.args.res.(*httptest.ResponseRecorder).Body.String(),
			)
		})
	}
}

type castableError struct {
	error

	Message string

	Self *castableError
}

func (c castableError) ToRESTAPIError() *RESTAPIError {
	return &RESTAPIError{
		Type:   "/errors/castable-error",
		Title:  "Castable Error",
		Detail: "A castable error occurred",
	}
}
