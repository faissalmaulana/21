package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/repository"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestProject(t *testing.T) {
	t.Cleanup(func() {
		_, err := testDB.Exec(`
				TRUNCATE TABLE
					projects
				RESTART IDENTITY CASCADE
			`)
		require.NoError(t, err)
	})

	t.Run("AddProject", func(t *testing.T) {
		project := repository.New(testDB, zap.NewNop())

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		prj := model.Project{
			Name: "Test Project",
		}

		err := project.AddProject(ctx, prj)
		require.NoError(t, err)
	})

}
