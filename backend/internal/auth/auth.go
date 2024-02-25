package auth

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"goadmin-backend/internal/domain"
)

var ErrInvalidToken = fmt.Errorf("invalid token")

type AuthService interface {
	VerifyToken(ctx context.Context, tokenString string) (*domain.User, error)
}

type authService struct {
	userRepo         domain.UserRepository
	revokedTokenRepo domain.RevokedTokenRepository
	jwtSecret        []byte
}

func NewAuthService(
	userRepo domain.UserRepository,
	RevokedTokenRepository domain.RevokedTokenRepository,
	jwtSecret []byte,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		revokedTokenRepo: RevokedTokenRepository,
		jwtSecret:        jwtSecret,
	}
}

func (a *authService) VerifyToken(
	ctx context.Context,
	tokenString string,
) (*domain.User, error) {
	// check if token is in revoked_token list
	isRevoked, err := a.revokedTokenRepo.IsRevoked(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("check revoked token error %w", err)
	}

	if isRevoked {
		return nil, ErrInvalidToken
	}
	claims := &domain.JWTClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return a.jwtSecret, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parse token error %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	user, err := a.userRepo.FindByUsername(ctx, claims.Username)
	if err != nil {
		return nil, fmt.Errorf("find user error %w", err)
	}

	return user, nil
}
