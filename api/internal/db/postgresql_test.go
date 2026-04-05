package db_test

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"testing"

	"github.com/faissalmaulana/21/api/internal/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func NewForTest(tb testing.TB, opts ...fx.Option) *fx.App {
	testOpts := []fx.Option{
		fx.Logger(fxtest.NewTestPrinter(tb)),
		fxtest.WithTestLogger(tb),
		fx.Provide(func() *zap.Logger {
			return zaptest.NewLogger(tb)
		}),
	}
	opts = append(testOpts, opts...)

	return fx.New(opts...)
}

func TestNewPostgresqlDB(t *testing.T) {
	require.NoError(t, godotenv.Load("../../.env"))

	dbport, err := strconv.Atoi(os.Getenv("PSQL_DB_PORT"))
	require.NoError(t, err)

	t.Run("ConnectionEstablish", func(t *testing.T) {
		app := NewForTest(t,
			fx.Provide(
				func(lc fx.Lifecycle, log *zap.Logger) (*sql.DB, error) {
					pg := db.NewPostgresqlDB(
						dbport,
						os.Getenv("PSQL_DB_HOST"),
						os.Getenv("PSQL_DB_USER"),
						os.Getenv("PSQL_DB_CREDENTIAL"),
						os.Getenv("PSQL_DB_NAME"),
					)
					return pg.Connect(lc, log)
				},
			),
			fx.Invoke(func(_ *sql.DB) {}),
		)

		require.NoError(t, app.Start(context.Background()))
		t.Cleanup(func() {
			require.NoError(t, app.Stop(context.Background()))
		})
	})
}
