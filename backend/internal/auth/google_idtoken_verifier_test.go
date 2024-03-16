package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"google.golang.org/api/oauth2/v2"

	"goadmin-backend/internal/domain"
)

func Test_authService_VerifyGoogleIDToken(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepo         domain.UserRepository
		revokedTokenRepo domain.RevokedTokenRepository
		jwtSecret        []byte
		oauth2Service    GoogleOAuth2Service
	}

	type args struct {
		idToken string
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
				oauth2Service:    &GoogleOAuth2ServiceMock{},
			},
			args: args{
				idToken: "idToken",
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
				oauth2Service: &GoogleOAuth2ServiceMock{
					hasError: true,
				},
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
				oauth2Service:    tt.fields.oauth2Service,
			}
			got, err := a.VerifyGoogleIDToken(
				context.Background(),
				tt.args.idToken,
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

func Test_verifyIDToken(t *testing.T) {
	t.Parallel()

	type args struct {
		oauth2Service GoogleOAuth2Service
		idToken       string
	}

	tests := []struct {
		name    string
		args    args
		want    *oauth2.Tokeninfo
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				oauth2Service: &GoogleOAuth2ServiceMock{},
				idToken:       "idToken",
			},
			want: &oauth2.Tokeninfo{},
		},
		{
			name: "Error",
			args: args{
				oauth2Service: &GoogleOAuth2ServiceMock{
					hasError: true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := verifyIDToken(
				context.Background(),
				tt.args.oauth2Service,
				tt.args.idToken,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"verifyIDToken() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("verifyIDToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

var _ GoogleTokeninfoCall = &googleTokeninfoCallMock{}

type googleTokeninfoCallMock struct {
	hasError bool
}

//nolint:revive,stylecheck // It's a mock
func (g *googleTokeninfoCallMock) IdToken(string) GoogleTokeninfoCall {
	return &googleTokeninfoCallMock{}
}

func (g *googleTokeninfoCallMock) Do() (*oauth2.Tokeninfo, error) {
	if g.hasError {
		return nil, errors.New("error")
	}

	return &oauth2.Tokeninfo{}, nil
}

var _ GoogleOAuth2Service = &GoogleOAuth2ServiceMock{}

type GoogleOAuth2ServiceMock struct {
	hasError bool
}

func (g *GoogleOAuth2ServiceMock) Tokeninfo() GoogleTokeninfoCall {
	return &googleTokeninfoCallMock{
		hasError: g.hasError,
	}
}
