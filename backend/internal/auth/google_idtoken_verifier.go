package auth

import (
	"context"
	"fmt"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/oauth2/v2"

	"goadmin-backend/internal/domain"
)

// GoogleTokeninfoCall is an interface for the oauth2 package's TokeninfoCall.
type GoogleTokeninfoCall interface {
	IdToken(string) *oauth2.TokeninfoCall
	Do(opts ...googleapi.CallOption) (*oauth2.Tokeninfo, error)
}

// GoogleOAuth2Service is an interface for the oauth2 package's Service.
type GoogleOAuth2Service interface {
	Tokeninfo() GoogleTokeninfoCall
}

func (a *authService) VerifyGoogleIDToken(
	ctx context.Context,
	idToken string,
) (*domain.User, error) {
	tokenInfo, err := verifyIDToken(a.oauth2Service, idToken)
	if err != nil {
		return nil, fmt.Errorf("verifyIDToken: %w", err)
	}

	user, err := a.userRepo.FindByUsername(ctx, tokenInfo.Email)
	if err != nil {
		return nil, fmt.Errorf("userRepo.FindByGoogleID: %w", err)
	}

	return user, nil
}

func verifyIDToken(
	oauth2Service GoogleOAuth2Service,
	idToken string,
) (*oauth2.Tokeninfo, error) {
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, fmt.Errorf("tokenInfoCall.Do: %w", err)
	}

	return tokenInfo, nil
}

// googleOAuth2Service is a concrete implementation of the GoogleOAuth2Service
// interface.
type googleOAuth2Service struct {
	service *oauth2.Service
}

type googleTokeninfoCall struct {
	*oauth2.TokeninfoCall
}

func (g *googleOAuth2Service) Tokeninfo() GoogleTokeninfoCall {
	return &googleTokeninfoCall{
		TokeninfoCall: g.service.Tokeninfo(),
	}
}

// NewGoogleOAuth2Service returns a new GoogleOAuth2Service.
func NewGoogleOAuth2Service(
	service *oauth2.Service,
	err error,
) (GoogleOAuth2Service, error) {
	return &googleOAuth2Service{
		service: service,
	}, err
}
