package postgres

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"goadmin-backend/internal/domain"
	"goadmin-backend/internal/platform/random"
)

var testUsers = []domain.User{
	{
		ID:        "1",
		Username:  "johndoe",
		Email:     "johndoe@goadmin.com",
		Password:  "test!only",
		FirstName: "John",
		LastName:  "Doe",
		Active:    true,
		Picture:   "https://example.com/johndoe.jpg",
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
		Picture:   "https://example.com/janedoe.jpg",
		Deleted:   false,
		CreatedAt: strToTime("2024-02-01T00:00:00+00:00"),
		UpdatedAt: strToTime("2024-02-01T00:00:00+00:00"),
	},
}

func ptr[T any](v T) *T {
	return &v
}

func randomUser() *domain.User {
	return &domain.User{
		Username:  random.String(10),
		Email:     random.String(10) + "@goadmin.com",
		Password:  random.String(10),
		FirstName: random.String(10),
		LastName:  random.String(10),
	}
}

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

	userRepo := NewUserRepo(conn)

	t.Cleanup(func() {
		teardown(t)
	})

	type args struct {
		filter *domain.UserFilter
	}

	tests := []struct {
		name    string
		args    args
		want    []domain.User
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				filter: &domain.UserFilter{
					CreatedBetween: [2]time.Time{
						{},
						strToTime("2024-02-01T00:00:01+00:00"),
					},
				},
			},
			want:    testUsers,
			wantErr: false,
		},
		{
			name: "success with filter",
			args: args{
				filter: &domain.UserFilter{
					FirstName: "John",
					Email:     "johndoe@goadmin.com",
				},
			},
			want: testUsers[:1],
		},
		{
			name: "success with filter",
			args: args{
				filter: &domain.UserFilter{
					LastName: "Doe",
					Email:    "johndoe@goadmin.com",
					Active:   ptr(true),
					Deleted:  ptr(false),
				},
			},
			want: testUsers[:1],
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := userRepo.FindAll(context.Background(), tt.args.filter)
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

	userRepo := NewUserRepo(conn)

	t.Cleanup(func() {
		teardown(t)
	})

	type args struct {
		id string
	}

	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				id: "1",
			},
			want: &testUsers[0],
		},
		{
			name: "fail",
			args: args{
				id: "aaaaaa",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := userRepo.FindByID(context.Background(), tt.args.id)
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

	userRepo := NewUserRepo(conn)

	t.Cleanup(func() {
		teardown(t)
	})

	type args struct {
		username string
	}

	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				username: "janedoe",
			},
			want: &testUsers[1],
		},
		{
			name: "fail",
			args: args{
				username: "aaaaaa",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := userRepo.FindByUsername(context.Background(), tt.args.username)
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

	userRepo := NewUserRepo(conn)

	t.Cleanup(func() {
		teardown(t)
	})

	type args struct {
		user *domain.User
	}

	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				user: &domain.User{
					Username:  "jackc",
					Email:     "jackc@goadmin.com",
					Password:  "test!only",
					FirstName: "Jack",
					LastName:  "C",
				},
			},
			want: &domain.User{
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

			got, err := userRepo.Create(context.Background(), tt.args.user)
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

			if got.ID != "" {
				tt.want.ID = got.ID
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepo_Update(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)

	userRepo := NewUserRepo(conn)

	testUsr, err := userRepo.Create(context.Background(), randomUser())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		teardown(t)
		userRepo.Delete(context.Background(), testUsr.ID)
	})

	type args struct {
		user *domain.User
	}

	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				user: &domain.User{
					ID:        testUsr.ID,
					Username:  "johndoe_updated",
					Email:     "johndoe_updated@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Picture:   "https://example.com/johndoe_updated.jpg",
				},
			},
			want: &domain.User{
				ID:        testUsr.ID,
				Username:  "johndoe_updated",
				Email:     "johndoe_updated@example.com",
				Password:  testUsr.Password,
				FirstName: "John",
				LastName:  "Doe",
				Active:    testUsr.Active,
				Picture:   "https://example.com/johndoe_updated.jpg",
				Deleted:   testUsr.Deleted,
				CreatedAt: testUsr.CreatedAt,
				UpdatedAt: testUsr.UpdatedAt,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := userRepo.Update(context.Background(), tt.args.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !got.CreatedAt.IsZero() {
				tt.want.CreatedAt = got.CreatedAt
			}

			if !got.UpdatedAt.IsZero() {
				tt.want.UpdatedAt = got.UpdatedAt
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepo_SoftDelete(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)

	userRepo := NewUserRepo(conn)

	testUsr, err := userRepo.Create(context.Background(), randomUser())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		teardown(t)
		userRepo.Delete(context.Background(), testUsr.ID)
	})

	type args struct {
		usrID string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				usrID: testUsr.ID,
			},
		},
		{
			name: "Fail",
			args: args{
				usrID: "aaaaaa",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := userRepo.SoftDelete(context.Background(), tt.args.usrID); (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.SoftDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepo_Delete(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)

	userRepo := NewUserRepo(conn)

	testUsr, err := userRepo.Create(context.Background(), randomUser())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		teardown(t)
		userRepo.Delete(context.Background(), testUsr.ID)
	})

	type args struct {
		id string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				id: testUsr.ID,
			},
		},
		{
			name: "Fail",
			args: args{
				id: "aaaaaa",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := userRepo.Delete(context.Background(), tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
