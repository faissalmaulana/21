package repository

import (
	"context"
	"database/sql"

	constant "github.com/faissalmaulana/21/api/internal/const"
	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/utils"
	"go.uber.org/zap"
)

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

	task.Status = utils.ToStatus(rawStatus)

	return task, nil
}
