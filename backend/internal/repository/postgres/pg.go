package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(
		ctx context.Context,
		sql string,
		args ...any,
	) (pgconn.CommandTag, error)
}

const defaultRetries = 50

type txConfigKey struct{}

func numRetriesFromContext(ctx context.Context) int {
	if v := ctx.Value(txConfigKey{}); v != nil {
		if retries, ok := v.(int); ok && retries >= 0 {
			return retries
		}
	}

	return defaultRetries
}

// Tx abstracts the operations needed by ExecuteInTx so that different
// frameworks (e.g. go's sql package, pgx, gorm) can be used with ExecuteInTx.
type Tx interface {
	Exec(context.Context, string, ...interface{}) error
	Commit(context.Context) error
	Rollback(context.Context) error
}

// ExecuteInTx runs fn inside tx. This method is primarily intended for internal
// use. See other packages for higher-level, framework-specific ExecuteTx()
// functions.
//
// *WARNING*: It is assumed that no statements have been executed on the
// supplied Tx. ExecuteInTx will only retry statements that are performed within
// the supplied closure (fn). Any statements performed on the tx before
// ExecuteInTx is invoked will *not* be re-run if the transaction needs to be
// retried.
//
// fn is subject to the same restrictions as the fn passed to ExecuteTx.

//nolint:nonamedreturns // named returns for better readability
func ExecuteInTx( //nolint:cyclop // This function is complex by nature.
	ctx context.Context,
	dbTx Tx,
	txFn func() error,
) (err error) {
	defer func() {
		rcvr := recover()

		if rcvr == nil && err == nil {
			// Ignore commit errors. The tx has already been committed by RELEASE.
			_ = dbTx.Commit( //nolint:errcheck // We're already handling an error.
				ctx,
			)

			return
		}

		// We always need to execute a Rollback() so sql.DB releases the
		// connection.
		_ = dbTx.Rollback( //nolint:errcheck // We're already handling an error.
			ctx,
		)

		if rcvr != nil {
			panic(rcvr)
		}
	}()

	// Specify that we intend to retry this txn in case of CockroachDB retryable
	// errors.
	if errExec := dbTx.Exec(ctx, "SAVEPOINT cockroach_restart"); errExec != nil {
		return errExec //nolint:wrapcheck // We're wrapping the error.
	}

	maxRetries := numRetriesFromContext(ctx)
	retryCount := 0

	for {
		releaseFailed := false

		err = txFn()
		if err == nil {
			// RELEASE acts like COMMIT in CockroachDB. We use it since it gives us an
			// opportunity to react to retryable errors, whereas tx.Commit() doesn't.
			if err = dbTx.Exec(ctx, "RELEASE SAVEPOINT cockroach_restart"); err == nil {
				return nil
			}

			releaseFailed = true
		}

		// We got an error; let's see if it's a retryable one and, if so, restart.
		if !errIsRetryable(err) {
			if releaseFailed {
				err = newAmbiguousCommitError(err)
			}

			return err
		}

		if rollbackErr := dbTx.Exec(ctx, "ROLLBACK TO SAVEPOINT cockroach_restart"); rollbackErr != nil {
			return newTxnRestartError(rollbackErr, err)
		}

		retryCount++
		if maxRetries > 0 && retryCount > maxRetries {
			return newMaxRetriesExceededError(err, maxRetries)
		}
	}
}

func errIsRetryable(err error) bool {
	// We look for either:
	//  - the standard PG errcode SerializationFailureError:40001 or
	//  - the Cockroach extension errcode RetriableError:CR000. This extension
	//    has been removed server-side, but support for it has been left here for
	//    now to maintain backwards compatibility.
	code := errCode(err)

	return code == "CR000" || code == "40001"
}

func errCode(err error) string {
	var sqlErr errWithSQLState
	if errors.As(err, &sqlErr) {
		return sqlErr.SQLState()
	}

	return ""
}

// errWithSQLState is implemented by pgx (pgconn.PgError) and lib/pq
type errWithSQLState interface {
	SQLState() string
}

func exec(
	ctx context.Context,
	db Queryer,
	query string,
	args ...any,
) (pgconn.CommandTag, error) { //nolint:unparam // We need to return a value.
	result, err := db.Exec(ctx, query, args...)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec error: %w", err)
	}

	return result, nil
}

func queryRow[T any](
	ctx context.Context,
	db Queryer,
	query string,
	args ...any,
) (*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query row error: %w", err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToAddrOfStructByName[T],
	)
	if err != nil {
		return nil, fmt.Errorf("query row error: %w", err)
	}

	return result, nil
}

func query[T any](
	ctx context.Context,
	db Queryer,
	query string,
	args ...any,
) ([]*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[T])
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return result, nil
}
