// Copyright 2016 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package postgres

import (
	"errors"
	"reflect"
	"testing"
)

func Test_txError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		cause error
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Success",
			fields: fields{
				cause: errors.New("error"),
			},
			want: "error",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &txError{
				cause: tt.fields.cause,
			}

			if got := e.Error(); got != tt.want {
				t.Errorf("txError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txError_Cause(t *testing.T) {
	t.Parallel()

	type fields struct {
		cause error
	}

	err := errors.New("error")

	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "Success",
			fields: fields{
				cause: err,
			},
			want: err,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &txError{
				cause: tt.fields.cause,
			}

			if got := e.Cause(); got != tt.want {
				t.Errorf("txError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txError_Unwrap(t *testing.T) {
	t.Parallel()

	err := errors.New("error")

	type fields struct {
		cause error
	}

	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "Success",
			fields: fields{
				cause: err,
			},
			want: err,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &txError{
				cause: tt.fields.cause,
			}

			if got := e.Unwrap(); got != tt.want {
				t.Errorf("txError.Unwrap() error = %v, wantErr %v", got, tt.want)
			}
		})
	}
}

func Test_newAmbiguousCommitError(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want *AmbiguousCommitError
	}{
		{
			name: "Success",
			args: args{
				err: errors.New("db error"),
			},
			want: &AmbiguousCommitError{
				txError{
					cause: errors.New("db error"),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			//nolint: govet // unit test
			if got := newAmbiguousCommitError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAmbiguousCommitError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newMaxRetriesExceededError(t *testing.T) {
	t.Parallel()

	type args struct {
		err        error
		maxRetries int
	}

	tests := []struct {
		name string
		args args
		want *MaxRetriesExceededError
	}{
		{
			name: "Success",
			args: args{
				err:        errors.New("error"),
				maxRetries: 3,
			},
			want: &MaxRetriesExceededError{
				txError: txError{
					cause: errors.New("error"),
				},
				msg: "retrying txn failed after 3 attempts. original error: error.",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			//nolint: govet //unit test
			if got := newMaxRetriesExceededError(tt.args.err, tt.args.maxRetries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMaxRetriesExceededError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxRetriesExceededError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		txError txError
		msg     string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Success",
			fields: fields{
				txError: txError{
					cause: errors.New("error"),
				},
				msg: "retrying txn failed after 3 attempts. original error: error.",
			},
			want: "retrying txn failed after 3 attempts. original error: error.",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &MaxRetriesExceededError{
				txError: tt.fields.txError,
				msg:     tt.fields.msg,
			}

			if got := e.Error(); got != tt.want {
				t.Errorf("MaxRetriesExceededError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newTxnRestartError(t *testing.T) {
	t.Parallel()

	type args struct {
		err      error
		retryErr error
	}

	tests := []struct {
		name string
		args args
		want *TxnRestartError
	}{
		{
			name: "Success",
			args: args{
				err:      errors.New("error"),
				retryErr: errors.New("retry error"),
			},
			want: &TxnRestartError{
				txError: txError{
					cause: errors.New("error"),
				},
				retryCause: errors.New("retry error"),
				msg:        "restarting txn failed. ROLLBACK TO SAVEPOINT encountered error: error. Original error: retry error.",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			//nolint:govet //unit test
			if got := newTxnRestartError(tt.args.err, tt.args.retryErr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTxnRestartError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTxnRestartError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		txError    txError
		retryCause error
		msg        string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Success",
			fields: fields{
				txError: txError{
					cause: errors.New("error"),
				},
				retryCause: errors.New("retry error"),
				msg:        "test error",
			},
			want: "test error",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &TxnRestartError{
				txError:    tt.fields.txError,
				retryCause: tt.fields.retryCause,
				msg:        tt.fields.msg,
			}

			if got := e.Error(); got != tt.want {
				t.Errorf("TxnRestartError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTxnRestartError_RetryCause(t *testing.T) {
	t.Parallel()

	type fields struct {
		txError    txError
		retryCause error
		msg        string
	}

	retryErr := errors.New("retry error")

	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "Success",
			fields: fields{
				txError: txError{
					cause: errors.New("error"),
				},
				retryCause: retryErr,
				msg:        "test error",
			},
			want: retryErr,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &TxnRestartError{
				txError:    tt.fields.txError,
				retryCause: tt.fields.retryCause,
				msg:        tt.fields.msg,
			}

			if got := e.RetryCause(); got != tt.want {
				t.Errorf("TxnRestartError.RetryCause() error = %v, wantErr %v", got, tt.want)
			}
		})
	}
}
