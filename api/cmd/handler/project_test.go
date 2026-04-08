package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/faissalmaulana/21/api/cmd/handler"
	constant "github.com/faissalmaulana/21/api/internal/const"
	"github.com/faissalmaulana/21/api/internal/mock"
	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/repository"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetProjects(t *testing.T) {
	projectRepoMock := new(mock.ProjectRepositoryMock)

	expectedProjects := []model.Project{
		{ID: "1", Name: "Project Alpha"},
		{ID: "2", Name: "Project Beta"},
	}

	expectedPagination := model.Pagination{
		Page:             1,
		Size:             constant.PaginateSize,
		TotalItemsInPage: 2,
		TotalItems:       2,
		TotalPages:       1,
	}

	projectRepoMock.
		On("Projects", testifyMock.Anything, repository.ProjectsParam{}).
		Return(expectedProjects, expectedPagination, nil)

	getProjectsHandler := handler.GetProjectsHandler{
		ProjectRepository: projectRepoMock,
	}
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	err := getProjectsHandler.HandleFunc(c)
	require.NoError(t, err)

	expectedResponse := handler.JSONResponse[handler.GetProjectsResponse]{
		Status: http.StatusOK,
		Data: handler.GetProjectsResponse{
			Projects: []handler.ProjectResponse{{ID: "1", Name: "Project Alpha"}, {ID: "2", Name: "Project Beta"}},
			Paginate: expectedPagination,
		},
		Error: nil,
	}

	assert.Equal(t, http.StatusOK, rec.Code)

	// unmarshal response body
	var actualResponse handler.JSONResponse[handler.GetProjectsResponse]
	err = json.Unmarshal(rec.Body.Bytes(), &actualResponse)
	require.NoError(t, err)

	assert.Equal(t, expectedResponse, actualResponse)

	projectRepoMock.AssertExpectations(t)
}
