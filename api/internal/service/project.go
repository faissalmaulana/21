package service

import (
	"context"
	"database/sql"
	"fmt"

	constant "github.com/faissalmaulana/21/api/internal/const"
	"github.com/faissalmaulana/21/api/internal/model"
)

type ProjectProvider interface {
	AddProject(ctx context.Context, prj model.Project) error
}

type Project struct {
	db *sql.DB
}

func New(db *sql.DB) *Project {
	return &Project{db: db}
}

func (p *Project) AddProject(ctx context.Context, prj model.Project) error {

	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	_, err := p.db.ExecContext(ctx, "INSERT INTO projects(name) VALUES($1)", prj.Name)
	if err != nil {
		return fmt.Errorf("AddProject: %v", err)
	}

	return nil
}
