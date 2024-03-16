package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"goadmin-backend/internal/domain"
)

// password hash of "password"
//
//nolint:gosec // for unit test
const passwordHash = "$2a$15$ioZSvZDYml862cyAk1l.x.AEGq77G1u8ruQQuA25Ic/QMpWsNDG/m"

func TestNewAuthService(t *testing.T) {
	t.Parallel()

	type args struct {
		userRepo         domain.UserRepository
		revokedTokenRepo domain.RevokedTokenRepository
		jwtSecret        []byte
		oauth2Service    GoogleOAuth2Service
	}

	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "Test NewAuthService()",
			args: args{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				oauth2Service:    &GoogleOAuth2ServiceMock{},
			},
			want: &authService{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				oauth2Service:    &GoogleOAuth2ServiceMock{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewAuthService(
				tt.args.userRepo,
				tt.args.revokedTokenRepo,
				tt.args.jwtSecret,
				tt.args.oauth2Service,
			); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewAuthService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authService_Login(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepo         domain.UserRepository
		revokedTokenRepo domain.RevokedTokenRepository
		jwtSecret        []byte
		oauth2Service    GoogleOAuth2Service
	}

	type args struct {
		credentials domain.Credentials
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Login Success",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				oauth2Service:    &GoogleOAuth2ServiceMock{},
			},
			args: args{
				credentials: domain.Credentials{
					Username: "username",
					Password: "password",
				},
			},
			want: true,
		},
		{
			name: "Login Fail",
			fields: fields{
				userRepo:         &UserRepositoryMock{hasError: true},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				oauth2Service:    &GoogleOAuth2ServiceMock{},
			},
			args: args{
				credentials: domain.Credentials{
					Username: "username",
					Password: "password",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := &authService{
				userRepo:         tt.fields.userRepo,
				revokedTokenRepo: tt.fields.revokedTokenRepo,
				jwtSecret:        tt.fields.jwtSecret,
				oauth2Service:    tt.fields.oauth2Service,
			}
			got, err := a.Login(context.Background(), tt.args.credentials)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"authService.Login() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			if _, err := a.VerifyToken(context.Background(), got); (err == nil) != tt.want {
				t.Errorf("authService.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authService_VerifyToken(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepo         domain.UserRepository
		revokedTokenRepo domain.RevokedTokenRepository
		jwtSecret        []byte
		oauth2Service    GoogleOAuth2Service
	}

	tests := []struct {
		name    string
		fields  fields
		want    *domain.User
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				oauth2Service:    &GoogleOAuth2ServiceMock{},
			},
			want: &domain.User{
				ID:       "1",
				Username: "username",
				Password: passwordHash,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := &authService{
				userRepo:         tt.fields.userRepo,
				revokedTokenRepo: tt.fields.revokedTokenRepo,
				jwtSecret:        tt.fields.jwtSecret,
				oauth2Service:    tt.fields.oauth2Service,
			}

			token, _ := a.Login(context.Background(), domain.Credentials{
				Username: "username",
				Password: "password",
			})

			got, err := a.VerifyToken(context.Background(), token)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"authService.VerifyToken() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"authService.VerifyToken() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func Test_authService_Logout(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepo         domain.UserRepository
		revokedTokenRepo domain.RevokedTokenRepository
		jwtSecret        []byte
		oauth2Service    GoogleOAuth2Service
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				oauth2Service:    &GoogleOAuth2ServiceMock{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := &authService{
				userRepo:         tt.fields.userRepo,
				revokedTokenRepo: tt.fields.revokedTokenRepo,
				jwtSecret:        tt.fields.jwtSecret,
				oauth2Service:    tt.fields.oauth2Service,
			}

			token, _ := a.Login(context.Background(), domain.Credentials{
				Username: "username",
				Password: "password",
			})

			if err := a.Logout(context.Background(), token); (err != nil) != tt.wantErr {
				t.Errorf(
					"authService.Logout() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func Test_authService_Register(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepo         domain.UserRepository
		revokedTokenRepo domain.RevokedTokenRepository
		jwtSecret        []byte
		oauth2Service    GoogleOAuth2Service
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
			name: "Success",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				oauth2Service:    &GoogleOAuth2ServiceMock{},
			},
			args: args{
				user: &domain.User{
					Username: "username",
					Password: "password",
				},
			},
			want: &domain.User{
				ID:       "1",
				Username: "username",
				Password: passwordHash,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := &authService{
				userRepo:         tt.fields.userRepo,
				revokedTokenRepo: tt.fields.revokedTokenRepo,
				jwtSecret:        tt.fields.jwtSecret,
				oauth2Service:    tt.fields.oauth2Service,
			}
			got, err := a.Register(context.Background(), tt.args.user)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"authService.Register() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authService.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

var _ domain.UserRepository = &UserRepositoryMock{}

type UserRepositoryMock struct {
	hasError bool
}

func (u *UserRepositoryMock) FindByUsername(
	_ context.Context,
	_ string,
) (*domain.User, error) {
	if u.hasError {
		return nil, domain.NewUserNotFoundError("1")
	}

	return &domain.User{
		ID:       "1",
		Username: "username",
		Password: passwordHash,
	}, nil
}

func (u *UserRepositoryMock) FindByID(
	_ context.Context,
	_ string,
) (*domain.User, error) {
	if u.hasError {
		return nil, domain.NewUserNotFoundError("1")
	}

	return &domain.User{
		ID:       "1",
		Username: "username",
		Password: passwordHash,
	}, nil
}

func (u *UserRepositoryMock) Create(
	_ context.Context,
	_ *domain.User,
) (*domain.User, error) {
	if u.hasError {
		return nil, errors.New("error")
	}

	return &domain.User{
		ID:       "1",
		Username: "username",
		Password: passwordHash,
	}, nil
}

func (u *UserRepositoryMock) Update(
	_ context.Context,
	_ *domain.User,
) (*domain.User, error) {
	if u.hasError {
		return nil, domain.NewUserNotFoundError("1")
	}

	return &domain.User{
		ID:       "1",
		Username: "username",
		Password: passwordHash,
	}, nil
}

func (u *UserRepositoryMock) Delete(_ context.Context, _ string) error {
	if u.hasError {
		return errors.New("error")
	}

	return nil
}

func (u *UserRepositoryMock) SoftDelete(_ context.Context, _ string) error {
	if u.hasError {
		return errors.New("error")
	}

	return nil
}

func (u *UserRepositoryMock) FindAll(
	_ context.Context,
	_ *domain.UserFilter,
) ([]domain.User, error) {
	if u.hasError {
		return nil, errors.New("error")
	}

	return []domain.User{
		{
			ID:       "1",
			Username: "username",
			Password: passwordHash,
		},
	}, nil
}

var _ domain.RevokedTokenRepository = &RevokedTokenRepositoryMock{}

type RevokedTokenRepositoryMock struct {
	hasError  bool
	isRevoked bool
}

func (r *RevokedTokenRepositoryMock) AddRevokedToken(
	_ context.Context,
	_ string,
) error {
	if r.hasError {
		return errors.New("error")
	}

	return nil
}

func (r *RevokedTokenRepositoryMock) IsRevoked(
	_ context.Context,
	_ string,
) (bool, error) {
	if r.hasError {
		return false, errors.New("error")
	}

	return r.isRevoked, nil
}
