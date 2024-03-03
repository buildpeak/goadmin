package postgres

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/buildpeak/sqltestutil"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// testPgConnStr is the connection string for the test database
var testPgConnStr string

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
