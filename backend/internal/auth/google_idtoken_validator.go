package auth

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/idtoken"

	"goadmin-backend/internal/domain"
)

var ErrInvalidIDToken = errors.New("invalid token")

// GoogleIDTokenValidator is an interface for validating Google ID tokens.
type GoogleIDTokenValidator interface {
	Validate(ctx context.Context, idToken, audience string) (*idtoken.Payload, error)
}

func (a *authService) ValidateGoogleIDToken(
	ctx context.Context,
	idToken string,
	audience string,
) (*domain.User, error) {
	if audience == "" {
		audience = a.audience
	}

	tokenInfo, err := a.idTokenValidator.Validate(ctx, idToken, audience)
	if err != nil {
		return nil, errors.Join(ErrInvalidIDToken, err)
	}

	user, err := a.userRepo.FindByUsername(ctx, tokenInfo.Claims["email"].(string))
	if err != nil {
		return nil, fmt.Errorf("error find user by username: %w", err)
	}

	return user, nil
}
