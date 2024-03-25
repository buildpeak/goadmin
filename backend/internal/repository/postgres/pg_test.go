package postgres

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/buildpeak/sqltestutil"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/goleak"
)

// testPgConnStr is the connection string for the test database
var testPgConnStr string

type pgxErr struct {
	error
	code string
}

func (e *pgxErr) SQLState() string {
	return e.code
}

var _ errWithSQLState = &pgxErr{}

var errCR000 = &pgxErr{code: "CR000"}

func TestMain(m *testing.M) {
	// setup
	ctx := context.Background()

	pgc, err := sqltestutil.StartPostgresContainer(ctx, "15")
	if err != nil {
		log.Fatalf("Failed to start postgres container: %v", err)
	}

	testPgConnStr = pgc.ConnectionString()

	log.Printf("Postgres container started: %s\n", testPgConnStr)

	db, errS := sql.Open("pgx", testPgConnStr)
	if errS != nil {
		pgc.Shutdown(ctx)
		log.Fatalf("Failed to open database: %v", errS)
	}

	errM := sqltestutil.RunMigrations(ctx, db, "../../../database/migrations")
	if errM != nil {
		pgc.Shutdown(ctx)
		log.Fatalf("Failed to run migrations: %v", errM)
	}

	sqltestutil.LoadScenario(ctx, db, "./testdata/scenario.yml")

	leak := flag.Bool("leak", false, "leak the container")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	exitCode := m.Run()

	// teardown
	log.Printf("Shutting down Postgres container\n")
	pgc.Shutdown(ctx)

	os.Exit(exitCode)
}

func TestPostgres(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("pgx", testPgConnStr)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Check tables
	rows, err := db.Query(`SELECT table_name FROM information_schema.tables
		WHERE table_schema = 'public'`)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()

	tables := []string{}

	for rows.Next() {
		var table string

		err = rows.Scan(&table)
		if err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}

		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("Failed to iterate rows: %v", err)
	}

	expectedTables := []string{
		"user",
		"role",
		"user_role",
		"permission",
		"role_permission",
		"user_permission",
		"revoked_token",
	}

	if len(tables) != len(expectedTables) {
		t.Fatalf("Expected %d tables, got %d", len(expectedTables), len(tables))
	}
}

func TestExecuteInTx(t *testing.T) {
	t.Parallel()

	type args struct {
		dbTx Tx
		txFn func() error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				dbTx: &txMock{},
				txFn: func() error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				dbTx: &txMock{
					execErr: sql.ErrTxDone,
				},
				txFn: func() error {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "Error 2",
			args: args{
				dbTx: &txMock{
					execErr:      sql.ErrTxDone,
					execErrCount: 2,
				},
				txFn: func() error {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "Error 3",
			args: args{
				dbTx: &txMock{
					execErr:      sql.ErrTxDone,
					execErrCount: 2,
				},
				txFn: func() error {
					return errCR000
				},
			},
			wantErr: true,
		},
		{
			name: "Error 4",
			args: args{
				dbTx: &txMock{
					execErr:      sql.ErrTxDone,
					execErrCount: 3,
				},
				txFn: func() error {
					return errCR000
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := ExecuteInTx(context.Background(), tt.args.dbTx, tt.args.txFn); (err != nil) != tt.wantErr {
				t.Errorf(
					"ExecuteInTx() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func Test_exec(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	type args struct {
		db   Queryer
		sql  string
		args []any
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				db:   conn,
				sql:  "SELECT 1",
				args: []any{},
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				db:   conn,
				sql:  "INSERT INTO non_existent_table VALUES (1)",
				args: []any{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if _, err := exec(
				context.Background(), tt.args.db, tt.args.sql, tt.args.args...,
			); (err != nil) != tt.wantErr {
				t.Errorf("exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_queryRow(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	type one struct {
		One int `db:"one"`
	}

	type args struct {
		db   Queryer
		sql  string
		args []any
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				db:   conn,
				sql:  "SELECT 1 AS one",
				args: []any{},
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				db:   conn,
				sql:  "SELECT non_existent_table",
				args: []any{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if _, err := queryRow[one](context.Background(), tt.args.db, tt.args.sql, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("queryRow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_query(t *testing.T) {
	t.Parallel()

	conn, teardown := before(t)
	t.Cleanup(func() {
		teardown(t)
	})

	type one struct {
		One int `db:"one"`
	}

	type args struct {
		db   Queryer
		sql  string
		args []any
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				db:   conn,
				sql:  "SELECT 1 AS one",
				args: []any{},
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				db:   conn,
				sql:  "SELECT non_existent_table",
				args: []any{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if _, err := query[one](context.Background(), tt.args.db, tt.args.sql, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("query() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type txMock struct {
	execCounter    int
	execCounterMux sync.Mutex
	execErrCount   int
	execErr        error
	commitErr      error
	rollbackErr    error
}

func (t *txMock) Commit(_ context.Context) error {
	return t.commitErr
}

func (t *txMock) Rollback(_ context.Context) error {
	return t.rollbackErr
}

func (t *txMock) Exec(_ context.Context, _ string, _ ...any) error {
	t.execCounterMux.Lock()
	defer t.execCounterMux.Unlock()

	t.execCounter++

	if t.execCounter >= t.execErrCount && t.execErr != nil {
		return t.execErr
	}

	return nil
}

type rowMock struct {
	err error
}

func (r *rowMock) Scan(_ ...any) error {
	return r.err
}

type queryerMock struct {
	err error
}

var _ Queryer = &queryerMock{}

func (q *queryerMock) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return &rowMock{err: q.err}
}

func (q *queryerMock) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, q.err
}

func (q *queryerMock) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, q.err
}
