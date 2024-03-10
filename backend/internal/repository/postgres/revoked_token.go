package postgres

import (
	"context"
	"fmt"
)

const revokedTokenTable = "revoked_token"

type RevokedTokenRepo struct {
	db Queryer
}

func NewRevokedTokenRepo(db Queryer) *RevokedTokenRepo {
	return &RevokedTokenRepo{
		db: db,
	}
}

func (rtr *RevokedTokenRepo) AddRevokedToken(ctx context.Context, token string) error {
	query := fmt.Sprintf(`INSERT INTO %s (token) VALUES ($1)`, revokedTokenTable)

	_, err := rtr.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("add revoked token error: %w", err)
	}

	return nil
}

func (rtr *RevokedTokenRepo) IsRevoked(ctx context.Context, token string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE token = $1)`, revokedTokenTable)

	var exists bool

	err := rtr.db.QueryRow(ctx, query, token).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("is revoked token error: %w", err)
	}

	return exists, nil
}
