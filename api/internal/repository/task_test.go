package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTask(t *testing.T) {
	t.Run("AddTask", func(t *testing.T) {
		t.Cleanup(func() {
			_, err := testDB.Exec(`
					TRUNCATE TABLE
						tasks,
						projects
					RESTART IDENTITY CASCADE
				`)
			require.NoError(t, err)
		})

		task := repository.NewTask(testDB, zap.NewNop())

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		projectID := seedProject(t, task.DB)
		newTask := model.Task{
			Name:      "Example Task",
			ProjectID: &projectID,
			StartAt:   time.Now(),
		}

		newTaskID, err := task.AddTask(ctx, newTask)
		require.NoError(t, err)
		assert.NotEmpty(t, newTaskID)
	})

	t.Run("GetProjectByID", func(t *testing.T) {
		t.Cleanup(func() {
			_, err := testDB.Exec(`
            TRUNCATE TABLE
            	tasks
            RESTART IDENTITY CASCADE
        `)
			require.NoError(t, err)
		})

		task := repository.NewTask(testDB, zap.NewNop())

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		newTaskID := seedTask(t, task.DB)

		testCases := []struct {
			name    string
			taskID  string
			wantID  string
			wantErr bool
		}{
			{
				name:    "success - get task",
				taskID:  newTaskID,
				wantID:  newTaskID,
				wantErr: false,
			},
			{
				name:    "fail - get task (not found)",
				taskID:  "00000000-0000-0000-0000-000000000000",
				wantID:  "",
				wantErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				resultTask, err := task.TaskByID(ctx, tc.taskID)

				if tc.wantErr {
					assert.Error(t, err)
					assert.Equal(t, tc.wantID, resultTask.ID)
					return
				}

				valid := resultTask.Status == model.Open || resultTask.Status == model.Done

				assert.Equal(t, tc.wantID, resultTask.ID)
				assert.NotEmpty(t, resultTask.Name)
				assert.True(t, valid, "invalid status: %v", resultTask.Status)
				assert.NotEmpty(t, resultTask.StartAt)
				assert.NotNil(t, resultTask.CreatedAt)
				assert.Nil(t, resultTask.LastUpdate)
				// i'm inserting without project_id
				assert.Nil(t, resultTask.ProjectID)
			})
		}
	})
}
