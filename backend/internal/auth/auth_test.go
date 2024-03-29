package auth

import (
	"context"
	"errors"
	"reflect"
	"strings"
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
		idTokenValidator GoogleIDTokenValidator
		audience         string
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
				idTokenValidator: &GoogleIDTokenValidatorMock{},
				audience:         "audience",
			},
			want: &authService{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
				audience:         "audience",
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
				tt.args.idTokenValidator,
				tt.args.audience,
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
		jwtSecret        interface{}
		idTokenValidator GoogleIDTokenValidator
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
				idTokenValidator: &GoogleIDTokenValidatorMock{},
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
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			args: args{
				credentials: domain.Credentials{
					Username: "username",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Credentials",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			args: args{
				credentials: domain.Credentials{
					Username: "username",
					Password: "invalid_password",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Secret",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        "invalid_secret", // non-bytes
				idTokenValidator: &GoogleIDTokenValidatorMock{},
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
				idTokenValidator: tt.fields.idTokenValidator,
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

			if got == nil && tt.wantErr {
				return
			}

			if _, err := a.VerifyToken(context.Background(), got.AccessToken); (err == nil) != tt.want {
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
		jwtSecret        interface{}
		idTokenValidator GoogleIDTokenValidator
	}

	type args struct {
		tokenString string
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
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			want: &domain.User{
				ID:       "1",
				Username: "username",
				Password: passwordHash,
			},
		},
		{
			name: "Get Revived Token Error",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{hasError: true},
				jwtSecret:        []byte("secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			wantErr: true,
		},
		{
			name: "Revived Token",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{isRevoked: true},
				jwtSecret:        []byte("secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			wantErr: true,
		},
		{
			name: "Find user error",
			fields: fields{
				userRepo:         &UserRepositoryMock{hasError: true},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Secret",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("invalid_secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imd1b2pmOTlAZ21haWwuY29tIiwiZXhwIjoxNzExMTY5MTM5fQ.iI7yE4lJ6eSnDKiBGrRQeTdYg_2zRj3cwtlfIEpLd8k",
			},
			wantErr: true,
		},
		{
			name: "Invalid Token",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imd1b2pmOTlAZ21haWwuY29tIiwiZXhwIjoxNzExMTY5MTM5fQ.iI7yE4lJ6eSnDKiBGrRQeTdYg_2zRj3cwtlfIEpLd8k",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := &authService{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			}

			token, _ := a.Login(context.Background(), domain.Credentials{
				Username: "username",
				Password: "password",
			})

			a.userRepo = tt.fields.userRepo
			a.revokedTokenRepo = tt.fields.revokedTokenRepo
			a.jwtSecret = tt.fields.jwtSecret
			a.idTokenValidator = tt.fields.idTokenValidator

			tokenString := token.AccessToken
			if tt.args.tokenString != "" {
				tokenString = tt.args.tokenString
			}

			got, err := a.VerifyToken(context.Background(), tokenString)

			t.Log(got, err)

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
		idTokenValidator GoogleIDTokenValidator
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
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
		},
		{
			name: "Add RevokedToken Error",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{hasError: true},
				jwtSecret:        []byte("secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
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
				idTokenValidator: tt.fields.idTokenValidator,
			}

			token, _ := a.Login(context.Background(), domain.Credentials{
				Username: "username",
				Password: "password",
			})

			if err := a.Logout(context.Background(), token.AccessToken); (err != nil) != tt.wantErr {
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
		idTokenValidator GoogleIDTokenValidator
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
				idTokenValidator: &GoogleIDTokenValidatorMock{},
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
		{
			name: "Create User Error",
			fields: fields{
				userRepo:         &UserRepositoryMock{hasError: true},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			args: args{
				user: &domain.User{
					Username: "username",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "Too Long Password",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("secret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			args: args{
				user: &domain.User{
					Username: "username",
					Password: strings.Repeat("z", 80),
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
				idTokenValidator: tt.fields.idTokenValidator,
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
		return nil, domain.NewResourceNotFoundError("User", "username=username")
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
		return nil, domain.NewResourceNotFoundError("User", "id=1")
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
		return nil, domain.NewResourceNotFoundError("User", "id=1")
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
