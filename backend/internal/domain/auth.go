package domain

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RoleRepository defines the methods that a role repository should implement
type RevokedTokenRepository interface {
	AddRevokedToken(ctx context.Context, token string) error
	IsRevoked(ctx context.Context, token string) (bool, error)
}
