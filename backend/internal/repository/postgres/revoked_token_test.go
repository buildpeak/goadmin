package postgres

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"reflect"
	"testing"
)

func randToken() string {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 32

	var result string

	max := big.NewInt(int64(len(alphabet)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}

		result += string(alphabet[n.Int64()])
	}

	return result
}

func TestNewRevokedTokenRepo(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	type args struct {
		db Queryer
	}

	tests := []struct {
		name string
		args args
		want *RevokedTokenRepo
	}{
		{
			name: "TestNewRevokedTokenRepo",
			args: args{
				db: conn,
			},
			want: &RevokedTokenRepo{
				db: conn,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewRevokedTokenRepo(tt.args.db); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewRevokedTokenRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRevokedTokenRepo_AddRevokedToken(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	repo := NewRevokedTokenRepo(conn)

	_ = repo.AddRevokedToken(context.Background(), "token1")

	type args struct {
		token string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestRevokedTokenRepo_AddRevokedToken",
			args: args{
				token: randToken(),
			},
			wantErr: false,
		},
		{
			name: "Fail",
			args: args{
				token: "token1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := repo.AddRevokedToken(context.Background(), tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf(
					"RevokedTokenRepo.AddRevokedToken() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestRevokedTokenRepo_IsRevoked(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	type fields struct {
		db Queryer
	}

	type args struct {
		token string
	}

	tests := []struct {
		name    string
		feilds  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Not Revoked",
			feilds: fields{
				db: conn,
			},
			args: args{
				token: "token9",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Revoked",
			feilds: fields{
				db: conn,
			},
			args: args{
				token: "token1",
			},
			want: true,
		},
		{
			name: "Error",
			feilds: fields{
				db: &queryerMock{err: errors.New("error")},
			},
			args: args{
				token: "token1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &RevokedTokenRepo{
				db: tt.feilds.db,
			}

			if tt.want {
				_ = repo.AddRevokedToken(context.Background(), tt.args.token)
			}

			got, err := repo.IsRevoked(context.Background(), tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"RevokedTokenRepo.IsRevoked() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			if got != tt.want {
				t.Errorf(
					"RevokedTokenRepo.IsRevoked() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}
