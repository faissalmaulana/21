package service_test

import (
	"context"
	"testing"

	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/service"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddProject(t *testing.T) {
	db := NewDBForTest(t)
	project := service.New(db, zap.NewExample())

	ctx := context.Background()
	prj := model.Project{
		Name: "Test Project",
	}

	err := project.AddProject(ctx, prj)
	require.NoError(t, err)
}
