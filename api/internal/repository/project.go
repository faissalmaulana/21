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
	AddProject(ctx context.Context, prj model.Project) (string, error)
	Projects(ctx context.Context, pp ProjectsParam) ([]model.Project, model.Pagination, error)
	GetProjectByID(ctx context.Context, id string) (model.Project, error)
	DeleteProjectByID(ctx context.Context, id string) (string, error)
	UpdateProject(ctx context.Context, project model.Project) error
}

type Project struct {
	db  *sql.DB
	log *zap.Logger
}

func NewProject(db *sql.DB, log *zap.Logger) *Project {
	return &Project{
		db:  db,
		log: log,
	}
}

func (p *Project) AddProject(ctx context.Context, prj model.Project) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var isArchive = false

	if prj.IsArchive == nil {
		prj.IsArchive = utils.BoolPtr(false)
	} else if prj.IsArchive != nil && *prj.IsArchive != false {
		isArchive = true
	}

	var newProjectID string

	err := p.db.QueryRowContext(ctx, "INSERT INTO projects(name,is_archive) VALUES($1,$2) RETURNING id", prj.Name, isArchive).Scan(&newProjectID)
	if err != nil {
		p.log.Error("Error AddProject", zap.Error(err))
		return "", MapDBError(err)
	}

	return newProjectID, nil
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

	if pp.Page < 1 {
		pp.Page = 1
	}

	offset := (pp.Page - 1) * pp.Size

	query := `
        SELECT id, name, COUNT(*) OVER() AS total_count
        FROM projects
        WHERE name ILIKE '%' || $1 || '%'
          AND is_archive = $2
        ORDER BY created_at DESC
        LIMIT $3 OFFSET $4
    `

	rows, err := p.db.QueryContext(ctx, query, pp.Search, pp.IsArchive, pp.Size, offset)
	if err != nil {
		p.log.Error("Error query projects", zap.Error(err))
		return nil, model.Pagination{}, MapDBError(err)
	}
	defer rows.Close()

	projects := make([]model.Project, 0)
	var totalItems int64 = 0

	for rows.Next() {
		project := new(model.Project)
		var count int64

		if err := rows.Scan(&project.ID, &project.Name, &count); err != nil {
			p.log.Error("Error scan project", zap.Error(err))
			return nil, model.Pagination{}, MapDBError(err)
		}

		totalItems = count
		projects = append(projects, *project)
	}

	if err := rows.Err(); err != nil {
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

func (p *Project) GetProjectByID(ctx context.Context, id string) (model.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var project model.Project

	if err := p.db.QueryRowContext(ctx, `SELECT id,name,is_archive,created_at FROM projects WHERE id = $1`, id).Scan(
		&project.ID,
		&project.Name,
		&project.IsArchive,
		&project.CreatedAt,
	); err != nil {
		p.log.Error("Error get project", zap.Error(err))
		return model.Project{}, MapDBError(err)
	}

	return project, nil
}

func (p *Project) DeleteProjectByID(ctx context.Context, id string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var deletedProjectID string
	if err := p.db.QueryRowContext(ctx, `DELETE FROM projects WHERE id = $1 RETURNING id`, id).Scan(&deletedProjectID); err != nil {
		p.log.Error("Error delete project", zap.Error(err))
		return "", MapDBError(err)
	}

	return id, nil
}

func (p *Project) UpdateProject(ctx context.Context, project model.Project) error {
	ctx, cancel := context.WithTimeout(ctx, constant.QueryTimeout)
	defer cancel()

	var archive bool
	foundProject, err := p.GetProjectByID(ctx, project.ID)
	if err != nil {
		p.log.Error("Error get project", zap.Error(err))
		return MapDBError(err)
	}

	if project.IsArchive == nil {
		archive = *foundProject.IsArchive
	} else {
		archive = *project.IsArchive
	}

	_, err = p.db.ExecContext(ctx, "UPDATE projects SET name = $1, is_archive = $2 WHERE id = $3", project.Name, archive, foundProject.ID)
	if err != nil {
		p.log.Error("Error update project", zap.Error(err))
		return MapDBError(err)
	}

	return nil
}
