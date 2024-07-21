package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/trad3r/hskills/apirest/internal/storage"
	"github.com/trad3r/hskills/apirest/internal/testpostgres"
)

func PreparePostgres(t *testing.T) (pgDB *storage.Storage, connStr string) {
	//TODO?
	//skipShort(t)
	ctx := context.Background()

	pgContainer, err := testpostgres.RunContainer(ctx,
		"postgres:16-alpine",
		"test-db",
		"postgres",
		"postgres",
		wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(5*time.Minute),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err = pgContainer.ConnectionDsn(ctx, "sslmode=disable")
	require.NoError(t, err)

	pgDB, err = storage.NewDB(ctx, connStr)
	require.NoError(t, err)

	return pgDB, connStr
}

func RunFixtures(fixturesPath string, dsn string) error {
	var db *sql.DB
	var fixtures *testfixtures.Loader
	var err error

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %s", err)
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(fixturesPath),
	)
	if err != nil {
		return fmt.Errorf("failed to create test fixtures: %s", err)
	}

	if err := fixtures.Load(); err != nil {
		return fmt.Errorf("failed to load test fixtures: %s", err)
	}

	return nil
}
