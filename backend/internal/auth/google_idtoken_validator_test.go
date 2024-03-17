package auth

import (
	"context"
	"errors"
	"goadmin-backend/internal/domain"
	"reflect"
	"testing"

	"google.golang.org/api/idtoken"
)

func Test_authService_VerifyGoogleIDToken(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepo         domain.UserRepository
		revokedTokenRepo domain.RevokedTokenRepository
		jwtSecret        []byte
		idTokenValidator GoogleIDTokenValidator
	}

	type args struct {
		idToken  string
		audience string
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
				jwtSecret:        []byte("jwtSecret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			args: args{
				idToken:  "idToken",
				audience: "audience",
			},
			want: &domain.User{
				ID:       "1",
				Username: "username",
				Password: passwordHash,
			},
		},
		{
			name: "Error",
			fields: fields{
				userRepo:         &UserRepositoryMock{},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("jwtSecret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{
					hasError: true,
				},
			},
			args: args{
				idToken:  "idToken",
				audience: "audience",
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
			got, err := a.ValidateGoogleIDToken(
				context.Background(),
				tt.args.idToken,
				tt.args.audience,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"authService.VerifyGoogleIDToken() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"authService.VerifyGoogleIDToken() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

var _ GoogleIDTokenValidator = &GoogleIDTokenValidatorMock{}

type GoogleIDTokenValidatorMock struct {
	hasError bool
}

func (g *GoogleIDTokenValidatorMock) Validate(
	_ context.Context,
	_ string,
	_ string,
) (*idtoken.Payload, error) {
	if g.hasError {
		return nil, errors.New("error")
	}

	return &idtoken.Payload{}, nil
}
