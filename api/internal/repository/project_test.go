package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/repository"
	"github.com/faissalmaulana/21/api/internal/utils"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestProject(t *testing.T) {

	t.Run("AddProject", func(t *testing.T) {
		t.Cleanup(func() {
			_, err := testDB.Exec(`
					TRUNCATE TABLE
						projects
					RESTART IDENTITY CASCADE
				`)
			require.NoError(t, err)
		})

		project := repository.New(testDB, zap.NewNop())

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		prj := model.Project{
			Name: "Test Project",
		}

		err := project.AddProject(ctx, prj)
		require.NoError(t, err)
	})

	t.Run("GetProjects", func(t *testing.T) {
		repo := repository.New(testDB, zap.NewExample())

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		testCases := []struct {
			name          string
			seedProjects  []model.Project
			expectedNames []string
			queryParams   repository.ProjectsParam
		}{
			{
				name:         "returns all projects",
				seedProjects: seedProjects(),
				expectedNames: []string{
					"Work Tasks",
					"Shopping List",
					"Daily Routine",
					"Weekend Plans",
					"Fitness Goals",
					"Study Tasks",
					"Home Chores",
				},
				queryParams: repository.ProjectsParam{},
			},
			{
				name:          "returns found projects",
				seedProjects:  seedProjects(),
				expectedNames: []string{"Fitness Goals"},
				queryParams: repository.ProjectsParam{
					Search: "Goals",
				},
			},
			{
				name:         "returns archived projects",
				seedProjects: seedProjects(),
				expectedNames: []string{
					"Project Ideas",
					"Errands",
					"Gym Goals",
				},
				queryParams: repository.ProjectsParam{
					IsArchive: true,
				},
			},
			{
				name:         "returns found archived projects",
				seedProjects: seedProjects(),
				expectedNames: []string{
					"Errands",
				},
				queryParams: repository.ProjectsParam{
					Search:    "errands",
					IsArchive: true,
				},
			},
			{
				name:          "returns empty projects (not found)",
				seedProjects:  seedProjects(),
				expectedNames: []string{},
				queryParams: repository.ProjectsParam{
					Search: "lizzy",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := testDB.Exec(`TRUNCATE TABLE projects RESTART IDENTITY CASCADE`)
				require.NoError(t, err)

				for _, p := range tc.seedProjects {
					err := repo.AddProject(ctx, p)
					require.NoError(t, err)
				}

				projects, err := repo.Projects(ctx, tc.queryParams)
				require.NoError(t, err)

				var projectNames []string
				for _, p := range projects {
					projectNames = append(projectNames, p.Name)
				}

				require.ElementsMatch(t, tc.expectedNames, projectNames)
			})
		}

	})

}

func seedProjects() []model.Project {
	return []model.Project{
		{Name: "Work Tasks", IsArchive: utils.BoolPtr(false)},
		{Name: "Shopping List", IsArchive: utils.BoolPtr(false)},
		{Name: "Daily Routine", IsArchive: utils.BoolPtr(false)},
		{Name: "Weekend Plans", IsArchive: utils.BoolPtr(false)},
		{Name: "Fitness Goals", IsArchive: utils.BoolPtr(false)},
		{Name: "Study Tasks", IsArchive: utils.BoolPtr(false)},
		{Name: "Home Chores", IsArchive: utils.BoolPtr(false)},
		{Name: "Project Ideas", IsArchive: utils.BoolPtr(true)},
		{Name: "Errands", IsArchive: utils.BoolPtr(true)},
		{Name: "Gym Goals", IsArchive: utils.BoolPtr(true)},
	}
}
