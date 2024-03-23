package auth

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/api/idtoken"

	"goadmin-backend/internal/domain"
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
		want    bool
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
			want: true,
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
		{
			name: "Get user error",
			fields: fields{
				userRepo:         &UserRepositoryMock{hasError: true},
				revokedTokenRepo: &RevokedTokenRepositoryMock{},
				jwtSecret:        []byte("jwtSecret"),
				idTokenValidator: &GoogleIDTokenValidatorMock{},
			},
			args: args{
				idToken: "idToken",
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

			if got == nil && tt.wantErr {
				return
			}

			if _, err := a.VerifyToken(context.Background(), got.AccessToken); (err == nil) != tt.want {
				t.Errorf("authService.VerifyToken() error = %v, wantErr %v", err, tt.want)
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

	return &idtoken.Payload{
		Subject: "subject",
		Claims:  map[string]interface{}{"email": "email"},
	}, nil
}
