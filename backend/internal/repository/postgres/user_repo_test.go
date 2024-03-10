package postgres

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"goadmin-backend/internal/domain"
)

func strToTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05-07:00", s)
	localTime, _ := time.LoadLocation("Local")

	return t.In(localTime)
}

func before(t *testing.T) (*pgxpool.Pool, func(t *testing.T)) {
	t.Helper()

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testPgConnStr)
	if err != nil {
		t.Fatal(err)
	}

	return pool, func(t *testing.T) {
		t.Helper()

		pool.Close()
	}
}

func TestNewUserRepo(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	type args struct {
		db *pgxpool.Pool
	}

	tests := []struct {
		name string
		args args
		want *UserRepo
	}{
		{
			name: "success",
			args: args{
				db: conn,
			},
			want: &UserRepo{
				db: conn,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewUserRepo(tt.args.db); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewUserRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UserRepo_FindAll(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	type fields struct {
		db *pgxpool.Pool
	}

	type args struct {
		filter *domain.UserFilter
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.User
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				db: conn,
			},
			args: args{
				filter: &domain.UserFilter{
					CreatedBetween: [2]time.Time{
						{},
						strToTime("2024-02-01T00:00:01+00:00"),
					},
				},
			},
			want: []domain.User{
				{
					ID:        "1",
					Username:  "johndoe",
					Email:     "johndoe@goadmin.com",
					Password:  "test!only",
					FirstName: "John",
					LastName:  "Doe",
					Active:    true,
					Deleted:   false,
					CreatedAt: strToTime("2024-02-01T00:00:00+00:00"),
					UpdatedAt: strToTime("2024-02-01T00:00:00+00:00"),
				},
				{
					ID:        "2",
					Username:  "janedoe",
					Email:     "janedoe@goadmin.com",
					Password:  "test!only",
					FirstName: "Jane",
					LastName:  "Doe",
					Active:    true,
					Deleted:   false,
					CreatedAt: strToTime("2024-02-01T00:00:00+00:00"),
					UpdatedAt: strToTime("2024-02-01T00:00:00+00:00"),
				},
			},
			wantErr: false,
		},
		{
			name: "success with filter",
			fields: fields{
				db: conn,
			},
			args: args{
				filter: &domain.UserFilter{
					FirstName: "John",
				},
			},
			want: []domain.User{
				{
					ID:        "1",
					Username:  "johndoe",
					Email:     "johndoe@goadmin.com",
					Password:  "test!only",
					FirstName: "John",
					LastName:  "Doe",
					Active:    true,
					Deleted:   false,
					CreatedAt: strToTime("2024-02-01T00:00:00+00:00"),
					UpdatedAt: strToTime("2024-02-01T00:00:00+00:00"),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &UserRepo{
				db: tt.fields.db,
			}

			got, err := r.FindAll(context.Background(), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"UserRepo.FindAll() name = %s error = %v, wantErr %v",
					tt.name,
					err,
					tt.wantErr,
				)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UserRepo_FindByID(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	type fields struct {
		db *pgxpool.Pool
	}

	type args struct {
		id string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				db: conn,
			},
			args: args{
				id: "1",
			},
			want: &domain.User{
				ID:        "1",
				Username:  "johndoe",
				Email:     "johndoe@goadmin.com",
				Password:  "test!only",
				FirstName: "John",
				LastName:  "Doe",
				Active:    true,
				Deleted:   false,
				CreatedAt: strToTime("2024-02-01T00:00:00+00:00"),
				UpdatedAt: strToTime("2024-02-01T00:00:00+00:00"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &UserRepo{
				db: tt.fields.db,
			}

			got, err := r.FindByID(context.Background(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"UserRepo.FindByID() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UserRepo_FindByUsername(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	type fields struct {
		db *pgxpool.Pool
	}

	type args struct {
		username string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				db: conn,
			},
			args: args{
				username: "janedoe",
			},
			want: &domain.User{
				ID:        "2",
				Username:  "janedoe",
				Email:     "janedoe@goadmin.com",
				Password:  "test!only",
				FirstName: "Jane",
				LastName:  "Doe",
				Active:    true,
				Deleted:   false,
				CreatedAt: strToTime("2024-02-01T00:00:00+00:00"),
				UpdatedAt: strToTime("2024-02-01T00:00:00+00:00"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &UserRepo{
				db: tt.fields.db,
			}

			got, err := r.FindByUsername(context.Background(), tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"UserRepo.FindByUsername() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"UserRepo.FindByUsername() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func Test_UserRepo_Create(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	type fields struct {
		db *pgxpool.Pool
	}

	type args struct {
		user *domain.User
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				db: conn,
			},
			args: args{
				user: &domain.User{
					ID:        "3",
					Username:  "jackc",
					Email:     "jackc@goadmin.com",
					Password:  "test!only",
					FirstName: "Jack",
					LastName:  "C",
				},
			},
			want: &domain.User{
				ID:        "3",
				Username:  "jackc",
				Email:     "jackc@goadmin.com",
				Password:  "test!only",
				FirstName: "Jack",
				LastName:  "C",
				Active:    true,
				Deleted:   false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &UserRepo{
				db: tt.fields.db,
			}

			got, err := r.Create(context.Background(), tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"UserRepo.Create() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			if !got.CreatedAt.IsZero() {
				tt.want.CreatedAt = got.CreatedAt
			}

			if !got.UpdatedAt.IsZero() {
				tt.want.UpdatedAt = got.UpdatedAt
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
