package postgres

import (
	"context"
	"reflect"
	"testing"
)

func TestNewRevokedTokenRepo(t *testing.T) {
	conn, teardown := before(t)
	defer teardown(t)

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
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRevokedTokenRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRevokedTokenRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRevokedTokenRepo_AddRevokedToken(t *testing.T) {
	conn, teardown := before(t)
	defer teardown(t)

	type fields struct {
		db Queryer
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestRevokedTokenRepo_AddRevokedToken",
			fields: fields{
				db: conn,
			},
			args: args{
				ctx:   context.Background(),
				token: "token",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rtr := &RevokedTokenRepo{
				db: tt.fields.db,
			}
			if err := rtr.AddRevokedToken(tt.args.ctx, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("RevokedTokenRepo.AddRevokedToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRevokedTokenRepo_IsRevoked(t *testing.T) {
	conn, teardown := before(t)
	defer teardown(t)

	type fields struct {
		db Queryer
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "TestRevokedTokenRepo_IsRevoked",
			fields: fields{
				db: conn,
			},
			args: args{
				ctx:   context.Background(),
				token: "token",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rtr := &RevokedTokenRepo{
				db: tt.fields.db,
			}
			got, err := rtr.IsRevoked(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RevokedTokenRepo.IsRevoked() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RevokedTokenRepo.IsRevoked() = %v, want %v", got, tt.want)
			}
		})
	}
}
