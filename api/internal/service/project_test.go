package service_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func NewDBForTest(t *testing.T) *sql.DB {
	require.NoError(t, godotenv.Load("../../.env"))

	t.Helper()

	dsn := os.Getenv("TEST_DB_DSN")
	require.NotEmpty(t, dsn, "TEST_DB_DSN must be set")

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestAddProject(t *testing.T) {

	db := NewDBForTest(t)
	project := service.New(db, zap.NewExample())

	ctx := context.Background()
	prj := model.Project{
		Name: "Test Project",
	}

	err := project.AddProject(ctx, prj)
	require.NoError(t, err)
}
