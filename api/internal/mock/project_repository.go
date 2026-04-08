package mock

import (
	"context"

	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/repository"
	"github.com/stretchr/testify/mock"
)

type ProjectRepositoryMock struct {
	mock.Mock
}

func (m *ProjectRepositoryMock) AddProject(ctx context.Context, prj model.Project) (string, error) {
	args := m.Called(ctx, prj)
	return args.String(0), args.Error(1)
}

func (m *ProjectRepositoryMock) Projects(ctx context.Context, pp repository.ProjectsParam) ([]model.Project, model.Pagination, error) {
	args := m.Called(ctx, pp)
	return args.Get(0).([]model.Project), args.Get(1).(model.Pagination), args.Error(2)
}

func (m *ProjectRepositoryMock) GetProjectByID(ctx context.Context, id string) (model.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Project), args.Error(1)
}

func (m *ProjectRepositoryMock) DeleteProjectByID(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *ProjectRepositoryMock) UpdateProject(ctx context.Context, project model.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}
