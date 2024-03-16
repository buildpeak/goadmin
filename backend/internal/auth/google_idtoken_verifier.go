package auth

import (
	"context"
	"fmt"

	"google.golang.org/api/oauth2/v2"

	"goadmin-backend/internal/domain"
)

// GoogleTokeninfoCall is an interface for the oauth2 package's TokeninfoCall.
type GoogleTokeninfoCall interface {
	IdToken(string) GoogleTokeninfoCall
	Do() (*oauth2.Tokeninfo, error)
}

// GoogleOAuth2Service is an interface for the oauth2 package's Service.
type GoogleOAuth2Service interface {
	Tokeninfo() GoogleTokeninfoCall
}

func (a *authService) VerifyGoogleIDToken(
	ctx context.Context,
	idToken string,
) (*domain.User, error) {
	tokenInfo, err := verifyIDToken(ctx, a.oauth2Service, idToken)
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
	ctx context.Context,
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

// googleTokeninfoCall is a wrapper for the oauth2 package's TokeninfoCall.
type googleTokeninfoCall struct {
	*oauth2.TokeninfoCall
}
