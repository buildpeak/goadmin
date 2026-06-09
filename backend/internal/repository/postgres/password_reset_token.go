package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"goadmin-backend/internal/domain"
)

var _ domain.PasswordResetTokenRepository = &PasswordResetTokenRepo{}

type PasswordResetTokenRepo struct {
	db Queryer
}

func NewPasswordResetTokenRepo(db Queryer) *PasswordResetTokenRepo {
	return &PasswordResetTokenRepo{db: db}
}

func (r *PasswordResetTokenRepo) Create(
	ctx context.Context,
	token *domain.PasswordResetToken,
) (*domain.PasswordResetToken, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING *`,
		passwordResetTokenTable,
	)

	result, err := queryRow[domain.PasswordResetToken](
		ctx, r.db, query,
		token.UserID, token.TokenHash, token.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create password reset token error: %w", err)
	}

	return result, nil
}

func (r *PasswordResetTokenRepo) FindByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*domain.PasswordResetToken, error) {
	query := fmt.Sprintf(
		`SELECT * FROM %s WHERE token_hash = $1`,
		passwordResetTokenTable,
	)

	result, err := queryRow[domain.PasswordResetToken](ctx, r.db, query, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError(
				"PasswordResetToken",
				"token_hash=<redacted>",
			)
		}

		return nil, fmt.Errorf("find password reset token error: %w", err)
	}

	return result, nil
}

func (r *PasswordResetTokenRepo) MarkUsed(ctx context.Context, id string) error {
	query := fmt.Sprintf(
		`UPDATE %s SET used = TRUE WHERE id = $1`,
		passwordResetTokenTable,
	)

	_, err := exec(ctx, r.db, query, id)
	if err != nil {
		return fmt.Errorf("mark password reset token used error: %w", err)
	}

	return nil
}
