package domain

import (
	"context"
	"time"
)

type PasswordResetToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *PasswordResetToken) (*PasswordResetToken, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string) error
}
