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
	DefaultTokenDuration        = 60 * time.Minute
	DefaultRefreshTokenDuration = 24 * time.Hour
	DefaultBCryptCost           = 15
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	Login(ctx context.Context, credentials domain.Credentials) (*domain.JWTToken, error)
	VerifyToken(ctx context.Context, tokenString string) (*domain.User, error)
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	ValidateGoogleIDToken(
		ctx context.Context,
		idToken, audience string,
	) (*domain.JWTToken, error)
}

var _ Service = &authService{}

type authService struct {
	userRepo         domain.UserRepository
	revokedTokenRepo domain.RevokedTokenRepository
	jwtSecret        interface{}
	idTokenValidator GoogleIDTokenValidator
	audience         string
}

func NewAuthService(
	userRepo domain.UserRepository,
	revokedTokenRepo domain.RevokedTokenRepository,
	jwtSecret []byte,
	idTokenValidator GoogleIDTokenValidator,
	audience string,
) Service {
	return &authService{
		userRepo:         userRepo,
		revokedTokenRepo: revokedTokenRepo,
		jwtSecret:        jwtSecret,
		idTokenValidator: idTokenValidator,
		audience:         audience,
	}
}

func (a *authService) Login(
	ctx context.Context,
	credentials domain.Credentials,
) (*domain.JWTToken, error) {
	user, err := a.userRepo.FindByUsername(ctx, credentials.Username)
	if err != nil {
		return nil, fmt.Errorf("find user error %w", err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(credentials.Password),
	)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return a.generateToken(user.Username)
}

func (a *authService) generateToken(username string) (*domain.JWTToken, error) {
	expirationTime := time.Now().Add(DefaultTokenDuration)
	claims := &domain.JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "goadmin-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("sign token error %w", err)
	}

	// Create the refresh token with longer expiry
	refreshTokenExpirationTime := time.Now().Add(DefaultRefreshTokenDuration)
	refreshClaims := &domain.JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpirationTime),
			Issuer:    "goadmin-backend",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString(a.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token error %w", err)
	}

	return &domain.JWTToken{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}, nil
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
