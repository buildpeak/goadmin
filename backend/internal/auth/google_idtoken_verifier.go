package auth

import (
	"context"
	"fmt"

	"google.golang.org/api/oauth2/v2"

	"goadmin-backend/internal/domain"
)

func (a *authService) VerifyGoogleIDToken(
	ctx context.Context,
	idToken string,
) (*domain.User, error) {
	tokenInfo, err := verifyIDToken(ctx, idToken)
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
	idToken string,
) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("oauth2.NewService: %w", err)
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, fmt.Errorf("tokenInfoCall.Do: %w", err)
	}

	return tokenInfo, nil
}
