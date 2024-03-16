package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"goadmin-backend/internal/domain"
)

const (
	DefaultTokenDuration = 60 * time.Minute
	DefaultBCryptCost    = 15
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	Login(ctx context.Context, credentials domain.Credentials) (string, error)
	VerifyToken(ctx context.Context, tokenString string) (*domain.User, error)
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	VerifyGoogleIDToken(
		ctx context.Context,
		idToken string,
	) (*domain.User, error)
}

var _ Service = &authService{}

type authService struct {
	userRepo         domain.UserRepository
	revokedTokenRepo domain.RevokedTokenRepository
	jwtSecret        []byte
	oauth2Service    GoogleOAuth2Service
}

func NewAuthService(
	userRepo domain.UserRepository,
	revokedTokenRepo domain.RevokedTokenRepository,
	jwtSecret []byte,
	oauth2Service GoogleOAuth2Service,
) Service {
	return &authService{
		userRepo:         userRepo,
		revokedTokenRepo: revokedTokenRepo,
		jwtSecret:        jwtSecret,
		oauth2Service:    oauth2Service,
	}
}

func (a *authService) Login(
	ctx context.Context,
	credentials domain.Credentials,
) (string, error) {
	user, err := a.userRepo.FindByUsername(ctx, credentials.Username)
	if err != nil {
		return "", fmt.Errorf("find user error %w", err)
	}

	if user == nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(credentials.Password),
	)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	expirationTime := time.Now().Add(DefaultTokenDuration)
	claims := &domain.JWTClaims{
		Username: credentials.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign token error %w", err)
	}

	return tokenString, nil
}

// VerifyToken checks if the token is valid and not revoked
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
		func(_ *jwt.Token) (any, error) {
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

// Logout invalidates the token by adding it to the revoked_token list
func (a *authService) Logout(
	ctx context.Context,
	tokenString string,
) error {
	err := a.revokedTokenRepo.AddRevokedToken(ctx, tokenString)
	if err != nil {
		return fmt.Errorf("revoke token error %w", err)
	}

	return nil
}

// Register registers a new user
func (a *authService) Register(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		DefaultBCryptCost,
	)
	if err != nil {
		return nil, fmt.Errorf("hash password error %w", err)
	}

	user.Password = string(hashedPassword)

	savedUser, err := a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user error %w", err)
	}

	return savedUser, nil
}
