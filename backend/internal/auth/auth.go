package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"goadmin-backend/internal/domain"
	"goadmin-backend/internal/mailer"
	"goadmin-backend/internal/totp"
)

const (
	DefaultTokenDuration         = 60 * time.Minute
	DefaultRefreshTokenDuration  = 24 * time.Hour
	DefaultBCryptCost            = 15
	DefaultResetTokenExpiry      = 1 * time.Hour
)

var (
	ErrInvalidToken          = errors.New("invalid token")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInvalidResetToken     = errors.New("invalid or expired reset token")
	ErrResetTokenUsed        = errors.New("reset token has already been used")
	ErrAccountLocked         = errors.New("account is locked due to too many failed attempts")
	MaxFailedAttempts        = 5
	LockDuration             = 15 * time.Minute
)

type Service interface {
	Login(ctx context.Context, credentials domain.Credentials) (*domain.JWTToken, error)
	VerifyToken(ctx context.Context, tokenString string) (*domain.User, error)
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	ValidateGoogleIDToken(
		ctx context.Context,
		idToken, audience string,
	) (*domain.JWTToken, error)
	Logout(ctx context.Context, tokenString string) error
	Profile(ctx context.Context, tokenString string) (*domain.User, error)
	RequestPasswordReset(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
	RefreshToken(ctx context.Context, refreshToken string) (*domain.JWTToken, error)
	SetupTOTP(ctx context.Context, tokenString string) (secret, qrURL string, err error)
	VerifyTOTP(ctx context.Context, tokenString, code string) error
	DisableTOTP(ctx context.Context, tokenString, password string) error
	Login2FA(ctx context.Context, tempToken, code string) (*domain.JWTToken, error)
}

var _ Service = &authService{}

type authService struct {
	userRepo              domain.UserRepository
	revokedTokenRepo      domain.RevokedTokenRepository
	passwordResetTokenRepo domain.PasswordResetTokenRepository
	jwtSecret             interface{}
	idTokenValidator      GoogleIDTokenValidator
	audience              string
	mailer                mailer.Mailer
	frontendURL           string
}

func NewAuthService(
	userRepo domain.UserRepository,
	revokedTokenRepo domain.RevokedTokenRepository,
	passwordResetTokenRepo domain.PasswordResetTokenRepository,
	jwtSecret []byte,
	idTokenValidator GoogleIDTokenValidator,
	audience string,
	mailer mailer.Mailer,
	frontendURL string,
) Service {
	return &authService{
		userRepo:              userRepo,
		revokedTokenRepo:      revokedTokenRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		jwtSecret:             jwtSecret,
		idTokenValidator:      idTokenValidator,
		audience:              audience,
		mailer:                mailer,
		frontendURL:           frontendURL,
	}
}

func (a *authService) Login(
	ctx context.Context,
	credentials domain.Credentials,
) (*domain.JWTToken, error) {
	user, err := a.userRepo.FindByUsername(ctx, credentials.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return nil, ErrAccountLocked
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(credentials.Password),
	)
	if err != nil {
		attempts := user.FailedLoginAttempts + 1
		var lockUntil *time.Time
		if attempts >= MaxFailedAttempts {
			t := time.Now().Add(LockDuration)
			lockUntil = &t
		}
		_ = a.userRepo.UpdateFailedAttempts(ctx, user.ID, attempts, lockUntil)
		return nil, ErrInvalidCredentials
	}

	if user.FailedLoginAttempts > 0 {
		_ = a.userRepo.UpdateFailedAttempts(ctx, user.ID, 0, nil)
	}

	if user.TotpEnabled && user.TotpSecret != nil && *user.TotpSecret != "" {
		tempToken, err := a.generateTempToken(user.Username)
		if err != nil {
			return nil, fmt.Errorf("generate temp token error: %w", err)
		}
		return &domain.JWTToken{
			AccessToken: tempToken,
			Requires2FA: true,
		}, nil
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

func generateResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate reset token error: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// RequestPasswordReset generates a reset token for the user with the given email.
// It always returns nil error (even if the email doesn't exist) to prevent
// email enumeration.
func (a *authService) RequestPasswordReset(
	ctx context.Context,
	email string,
) (string, error) {
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal whether the email exists
		return "", nil
	}

	rawToken, err := generateResetToken()
	if err != nil {
		return "", fmt.Errorf("generate reset token error: %w", err)
	}

	tokenHash := hashToken(rawToken)

	_, err = a.passwordResetTokenRepo.Create(ctx, &domain.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(DefaultResetTokenExpiry),
	})
	if err != nil {
		return "", fmt.Errorf("create reset token error: %w", err)
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", a.frontendURL, rawToken)
	if err := a.mailer.SendPasswordReset(user.Email, user.Username, resetURL); err != nil {
		return "", fmt.Errorf("send password reset email error: %w", err)
	}

	return rawToken, nil
}

// ResetPassword resets the user's password using a valid reset token.
func (a *authService) ResetPassword(
	ctx context.Context,
	token, newPassword string,
) error {
	tokenHash := hashToken(token)

	resetToken, err := a.passwordResetTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return ErrInvalidResetToken
	}

	if resetToken.Used {
		return ErrResetTokenUsed
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return ErrInvalidResetToken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		DefaultBCryptCost,
	)
	if err != nil {
		return fmt.Errorf("hash password error: %w", err)
	}

	if err := a.userRepo.UpdatePassword(ctx, resetToken.UserID, string(hashedPassword)); err != nil {
		return fmt.Errorf("update password error: %w", err)
	}

	if err := a.passwordResetTokenRepo.MarkUsed(ctx, resetToken.ID); err != nil {
		return fmt.Errorf("mark token used error: %w", err)
	}

	return nil
}

const DefaultTempTokenDuration = 2 * time.Minute

func (a *authService) generateTempToken(username string) (string, error) {
	claims := &domain.JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(DefaultTempTokenDuration)),
			Issuer:    "goadmin-backend-2fa",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

func (a *authService) SetupTOTP(
	ctx context.Context,
	tokenString string,
) (string, string, error) {
	user, err := a.VerifyToken(ctx, tokenString)
	if err != nil {
		return "", "", fmt.Errorf("verify token: %w", err)
	}

	secret, err := totp.GenerateSecret()
	if err != nil {
		return "", "", err
	}

	secretStr := secret; user.TotpSecret = &secretStr
	if _, err := a.userRepo.Update(ctx, user); err != nil {
		return "", "", fmt.Errorf("update user: %w", err)
	}

	qrURL := totp.QRURL(secret, user.Username, "GoAdmin")
	return secret, qrURL, nil
}

func (a *authService) VerifyTOTP(
	ctx context.Context,
	tokenString, code string,
) error {
	user, err := a.VerifyToken(ctx, tokenString)
	if err != nil {
		return fmt.Errorf("verify token: %w", err)
	}

	if user.TotpSecret == nil || *user.TotpSecret == "" {
		return errors.New("totp not set up")
	}

	if !totp.Validate(*user.TotpSecret, code) {
		return errors.New("invalid totp code")
	}

	user.TotpEnabled = true
	if _, err := a.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (a *authService) DisableTOTP(
	ctx context.Context,
	tokenString, password string,
) error {
	user, err := a.VerifyToken(ctx, tokenString)
	if err != nil {
		return fmt.Errorf("verify token: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}

	emptySecret := ""
	user.TotpSecret = &emptySecret
	user.TotpEnabled = false
	if _, err := a.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (a *authService) Login2FA(
	ctx context.Context,
	tempToken, code string,
) (*domain.JWTToken, error) {
	claims := &domain.JWTClaims{}
	token, err := jwt.ParseWithClaims(tempToken, claims, func(_ *jwt.Token) (any, error) {
		return a.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Issuer != "goadmin-backend-2fa" {
		return nil, ErrInvalidToken
	}

	user, err := a.userRepo.FindByUsername(ctx, claims.Username)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !totp.Validate(*user.TotpSecret, code) {
		return nil, errors.New("invalid totp code")
	}

	// Revoke the temp token so it can't be reused
	if err := a.revokedTokenRepo.AddRevokedToken(ctx, tempToken); err != nil {
		return nil, fmt.Errorf("revoke temp token: %w", err)
	}

	return a.generateToken(user.Username)
}

// RefreshToken validates a refresh token and returns a new token pair.
func (a *authService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*domain.JWTToken, error) {
	isRevoked, err := a.revokedTokenRepo.IsRevoked(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("check revoked token error: %w", err)
	}

	if isRevoked {
		return nil, ErrInvalidToken
	}

	claims := &domain.JWTClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(_ *jwt.Token) (any, error) {
		return a.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	user, err := a.userRepo.FindByUsername(ctx, claims.Username)
	if err != nil {
		return nil, ErrInvalidToken
	}

	_ = user // user exists, proceed

	return a.generateToken(claims.Username)
}

// Profile returns the current user profile
func (a *authService) Profile(
	ctx context.Context,
	tokenString string,
) (*domain.User, error) {
	user, err := a.VerifyToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("verify token error %w", err)
	}

	return user, nil
}
