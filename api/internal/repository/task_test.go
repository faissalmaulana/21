package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/repository"
	"github.com/faissalmaulana/21/api/internal/utils"
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

		now := time.Now()
		projectID := seedProject(t, task.DB)
		newTask := model.Task{
			Name:      "Example Task",
			ProjectID: &projectID,
			StartAt:   &now,
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

				valid := *resultTask.Status == model.Open || *resultTask.Status == model.Done

				assert.Equal(t, tc.wantID, resultTask.ID)
				assert.NotEmpty(t, resultTask.Name)
				assert.True(t, valid, "invalid status: %v", resultTask.Status)
				assert.NotNil(t, resultTask.StartAt)
				assert.NotNil(t, resultTask.CreatedAt)
				assert.Nil(t, resultTask.LastUpdate)
				assert.Nil(t, resultTask.Project)
				// i'm inserting without project_id
				assert.Nil(t, resultTask.ProjectID)
			})
		}
	})

	t.Run("updateTask", func(t *testing.T) {
		t.Cleanup(func() {
			_, err := testDB.Exec(`TRUNCATE TABLE projects RESTART IDENTITY CASCADE `)
			require.NoError(t, err)
		})

		task := repository.NewTask(testDB, zap.NewExample())

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		newTaskID := seedTask(t, task.DB)
		projectID := seedProject(t, task.DB)

		startAt := time.Now().Add(time.Minute * 60)
		done := utils.ToStatus("done")

		testCases := []struct {
			name    string
			taskID  string
			input   model.Task
			wantErr bool
		}{
			{
				name:   "update (name,projectID,start_at,status)",
				taskID: newTaskID,
				input: model.Task{
					Name:      "Test Task",
					ProjectID: &projectID,
					StartAt:   &startAt,
					Status:    &done,
				},
				wantErr: false,
			},
			{
				name:    "task not found",
				taskID:  "00000000-0000-0000-0000-000000000000",
				input:   model.Task{},
				wantErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Cleanup(func() {
					_, err := testDB.Exec(`TRUNCATE TABLE tasks RESTART IDENTITY CASCADE `)
					require.NoError(t, err)
				})

				err := task.UpdateTask(ctx, tc.taskID, tc.input)
				if tc.wantErr {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)
			})
		}
	})
}
