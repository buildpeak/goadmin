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
	Requires2FA  bool   `json:"requires_2fa,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Requires2FA  bool   `json:"requires_2fa,omitempty"`
	TempToken    string `json:"temp_token,omitempty"`
}

// RoleRepository defines the methods that a role repository should implement
type RevokedTokenRepository interface {
	AddRevokedToken(ctx context.Context, token string) error
	IsRevoked(ctx context.Context, token string) (bool, error)
}
