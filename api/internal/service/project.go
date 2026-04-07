package service

import (
	"context"
	"database/sql"

	constant "github.com/faissalmaulana/21/api/internal/const"
	"github.com/faissalmaulana/21/api/internal/model"
	"go.uber.org/zap"
)

type ProjectProvider interface {
	AddProject(ctx context.Context, prj model.Project) error
}

type Project struct {
	db  *sql.DB
	log *zap.Logger
}

func New(db *sql.DB, log *zap.Logger) *Project {
	return &Project{
		db:  db,
		log: log,
	}
}

func (p *Project) AddProject(ctx context.Context, prj model.Project) error {

	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	_, err := p.db.ExecContext(ctx, "INSERT INTO projects(name) VALUES($1)", prj.Name)
	if err != nil {
		p.log.Error("Error AddProject", zap.Error(err))
		return MapDBError(err)
	}

	return nil
}
