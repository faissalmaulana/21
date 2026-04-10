package repository

import (
	"context"
	"database/sql"

	constant "github.com/faissalmaulana/21/api/internal/const"
	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/utils"
	"go.uber.org/zap"
)

type TaskRepository interface {
	AddTask(ctx context.Context, task model.Task) (string, error)
	TaskByID(ctx context.Context, id string) (model.Task, error)
	UpdateTask(ctx context.Context, id string, updatedTask model.Task) error
	Tasks(ctx context.Context) ([]model.Task, error)
	DeleteTaskByID(ctx context.Context, id string) (string, error)
}

type Task struct {
	DB  *sql.DB
	Log *zap.Logger
}

func NewTask(db *sql.DB, log *zap.Logger) *Task {
	return &Task{db, log}
}

func (t *Task) AddTask(ctx context.Context, task model.Task) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var newTaskID string

	err := t.DB.QueryRowContext(ctx,
		"INSERT INTO tasks(name,project_id,start_at) VALUES($1,$2,$3) RETURNING id",
		task.Name,
		*task.ProjectID,
		task.StartAt,
	).Scan(&newTaskID)
	if err != nil {
		t.Log.Error("Error AddTask", zap.Error(err))
		return "", MapDBError(err)
	}

	return newTaskID, nil
}

func (t *Task) TaskByID(ctx context.Context, id string) (model.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var (
		task      model.Task
		rawStatus string
	)

	if err := t.DB.QueryRowContext(ctx, `SELECT id,name,project_id,status,start_at,created_at,last_update FROM tasks WHERE id = $1`, id).Scan(
		&task.ID,
		&task.Name,
		&task.ProjectID,
		&rawStatus,
		&task.StartAt,
		&task.CreatedAt,
		&task.LastUpdate,
	); err != nil {
		t.Log.Error("Error get task", zap.Error(err))
		return model.Task{}, MapDBError(err)
	}

	s := utils.ToStatus(rawStatus)
	task.Status = &s

	return task, nil
}

func (t *Task) UpdateTask(ctx context.Context, id string, updatedTask model.Task) error {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	task, err := t.TaskByID(ctx, id)
	if err != nil {
		t.Log.Error("Error get task", zap.Error(err))
		return MapDBError(err)
	}

	if updatedTask.Name != "" {
		task.Name = updatedTask.Name
	}

	if updatedTask.ProjectID != nil {
		task.ProjectID = updatedTask.ProjectID
	}

	if updatedTask.StartAt != nil {
		task.StartAt = updatedTask.StartAt
	}

	if updatedTask.Status != nil {
		task.Status = updatedTask.Status
	}

	_, err = t.DB.ExecContext(
		ctx,
		"UPDATE tasks SET name = $1, project_id = $2, start_at = $3, status = $4, last_update = NOW() WHERE id = $5",
		task.Name,
		task.ProjectID,
		*task.StartAt,
		task.Status.String(),
		task.ID,
	)
	if err != nil {
		t.Log.Error("Error update task", zap.Error(err))
		return MapDBError(err)
	}

	return nil
}

func (t *Task) Tasks(ctx context.Context) ([]model.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	rows, err := t.DB.QueryContext(
		ctx,
		`SELECT id,name,project_id,status,start_at,p.id,p.name AS project_name
		FROM tasks JOIN projects p ON p.id = project_id ORDER BY tasks.created_at DESC
		`,
	)
	if err != nil {
		t.Log.Error("Error querying get tasks", zap.Error(err))
		return nil, MapDBError(err)
	}

	tasks := make([]model.Task, 0)

	for rows.Next() {
		rawStatus := ""
		task := model.Task{}
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.ProjectID,
			&rawStatus,
			&task.StartAt,
			&task.Project.ID,
			&task.Project.Name,
		); err != nil {
			t.Log.Error("Error populate task", zap.Error(err))
			return nil, MapDBError(err)
		}
		s := utils.ToStatus(rawStatus)
		task.Status = &s

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *Task) DeleteTaskByID(ctx context.Context, id string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var deletedTaskID string
	if err := t.DB.QueryRowContext(ctx, `DELETE FROM tasks WHERE id = $1 RETURNING id`, id).Scan(&deletedTaskID); err != nil {
		t.Log.Error("Error delete task", zap.Error(err))
		return "", MapDBError(err)
	}

	return id, nil
}
