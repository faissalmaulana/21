package service_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/faissalmaulana/21/api/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func NewDBForTest(t *testing.T) service.Querier {
	t.Helper()

	require.NoError(t, godotenv.Load("../../.env"))

	dsn := os.Getenv("TEST_DB_DSN")
	require.NotEmpty(t, dsn, "TEST_DB_DSN must be set")

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	tx, err := db.Begin()
	require.NoError(t, err)

	t.Cleanup(func() {
		tx.Rollback()
		db.Close()
	})

	return tx
}
