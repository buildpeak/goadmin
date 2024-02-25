package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func executeInTx(ctx context.Context, db *pgx.Conn, fn func(*pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(&tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func exec(ctx context.Context, db Queryer, query string, args ...any) (pgconn.CommandTag, error) {
	return db.Exec(ctx, query, args...)
}

func queryRow[T any](ctx context.Context, db Queryer, query string, args ...any) (*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[T])
}

func query[T any](ctx context.Context, db Queryer, query string, args ...any) ([]*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[T])
}
