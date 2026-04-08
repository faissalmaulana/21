package repository

import (
	"context"
	"database/sql"

	constant "github.com/faissalmaulana/21/api/internal/const"
	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/utils"
	"go.uber.org/zap"
)

type ProjectRepository interface {
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

	var isArchive = false

	if prj.IsArchive == nil {
		prj.IsArchive = utils.BoolPtr(false)
	} else if prj.IsArchive != nil && *prj.IsArchive != false {
		isArchive = true
	}

	_, err := p.db.ExecContext(ctx, "INSERT INTO projects(name,is_archive) VALUES($1,$2)", prj.Name, isArchive)
	if err != nil {
		p.log.Error("Error AddProject", zap.Error(err))
		return MapDBError(err)
	}

	return nil
}

type ProjectsParam struct {
	Search    string
	IsArchive bool
	Page      int
	Size      int
}

func (p *Project) Projects(ctx context.Context, pp ProjectsParam) ([]model.Project, model.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var offset int

	if pp.Page <= 1 {
		offset = 0
	} else {
		offset = pp.Page * pp.Size
	}

	query := `
        SELECT id, name
        FROM projects
        WHERE name ILIKE '%' || $1 || '%'
          AND is_archive = $2
        ORDER BY created_at DESC
        LIMIT $3 OFFSET $4`

	rows, err := p.db.QueryContext(ctx, query, pp.Search, pp.IsArchive, pp.Size, offset)
	if err != nil {
		p.log.Error("Error query projects", zap.Error(err))
		return nil, model.Pagination{}, MapDBError(err)
	}
	defer rows.Close()

	projects := make([]model.Project, 0)

	for rows.Next() {
		project := new(model.Project)
		if err := rows.Scan(&project.ID, &project.Name); err != nil {
			p.log.Error("Error scan project", zap.Error(err))
			return nil, model.Pagination{}, MapDBError(err)
		}
		projects = append(projects, *project)
	}

	if err := rows.Err(); err != nil {
		return nil, model.Pagination{}, MapDBError(err)
	}

	var totalItems int64

	err = p.db.QueryRowContext(ctx, `SELECT COUNT(id) FROM projects`).Scan(&totalItems)
	if err != nil {
		p.log.Error("Error counting projects", zap.Error(err))
		return nil, model.Pagination{}, MapDBError(err)
	}

	var totalPages int
	if pp.Size > 0 {
		totalPages = int((totalItems + int64(pp.Size) - 1) / int64(pp.Size))
	}

	paginate := model.Pagination{
		Page:             pp.Page,
		Size:             pp.Size,
		TotalItemsInPage: len(projects),
		TotalItems:       int(totalItems),
		TotalPages:       totalPages,
	}

	return projects, paginate, nil
}
