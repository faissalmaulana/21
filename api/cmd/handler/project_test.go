package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/faissalmaulana/21/api/cmd/handler"
	constant "github.com/faissalmaulana/21/api/internal/const"
	"github.com/faissalmaulana/21/api/internal/mock"
	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/repository"
	"github.com/faissalmaulana/21/api/internal/service"
	"github.com/go-playground/validator/v10"
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
		On("Projects", testifyMock.Anything, repository.ProjectsParam{
			// in handler, size has default value.
			Size: constant.PaginateSize,
		}).
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

func TestPostProject(t *testing.T) {
	e := echo.New()

	projectRepoMock := new(mock.ProjectRepositoryMock)
	validate := validator.New()

	postHandler := handler.PostProjectHandler{
		ProjectRepository:   projectRepoMock,
		Validator:           validate,
		ValidatorSugaredMsg: service.NewSugaredErrorMessageValidator(validate),
	}

	tests := []struct {
		name           string
		body           string
		setupMock      func()
		expectedStatus int
		assertBody     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			body: `{"name":"Test Project"}`,
			setupMock: func() {
				projectRepoMock.
					On("AddProject", testifyMock.Anything, testifyMock.MatchedBy(func(p model.Project) bool {
						return p.Name == "Test Project"
					})).
					Return("", nil).
					Once()
			},
			expectedStatus: http.StatusCreated,
			assertBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				expected := `{"status":201,"data":"Add New Project Successfully","error":null}`
				assert.JSONEq(t, expected, rec.Body.String())
			},
		},
		{
			name:           "validation error - empty name",
			body:           `{"name":""}`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			assertBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp handler.JSONResponse[any]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.Equal(t, http.StatusBadRequest, resp.Status)
				assert.Nil(t, resp.Data)
				require.NotNil(t, resp.Error)
				assert.NotEmpty(t, resp.Error.Message)
				assert.Contains(t, strings.ToLower(resp.Error.Message), "name")
			},
		},
		{
			name:           "validation error - name > 255",
			body:           fmt.Sprintf(`{"name":"%s"}`, strings.Repeat("a", 256)),
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			assertBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp handler.JSONResponse[any]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.Equal(t, http.StatusBadRequest, resp.Status)
				assert.Nil(t, resp.Data)
				require.NotNil(t, resp.Error)
				assert.NotEmpty(t, resp.Error.Message)
				assert.Contains(t, strings.ToLower(resp.Error.Message), "name")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRepoMock.ExpectedCalls = nil // reset mock per case

			tt.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := postHandler.HandleFunc(c)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			tt.assertBody(t, rec)

			projectRepoMock.AssertExpectations(t)
		})
	}
}

func TestDeleteProject(t *testing.T) {
	e := echo.New()

	projectRepoMock := new(mock.ProjectRepositoryMock)

	deleteHandler := handler.DeleteProjectHandler{
		ProjectRepository: projectRepoMock,
	}

	tests := []struct {
		name           string
		projectID      string
		setupMock      func()
		expectedStatus int
		assertBody     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:      "success delete",
			projectID: "123",
			setupMock: func() {
				projectRepoMock.
					On("DeleteProjectByID", testifyMock.Anything, "123").
					Return("123", nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			assertBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				expected := `{
					"status":200,
					"data":"Project with id 123 deleted successfully",
					"error":null
				}`
				assert.JSONEq(t, expected, rec.Body.String())
			},
		},
		{
			name:      "not found project",
			projectID: "999",
			setupMock: func() {
				projectRepoMock.
					On("DeleteProjectByID", testifyMock.Anything, "999").
					Return("", repository.ErrNotFound).
					Once()
			},
			expectedStatus: http.StatusNotFound,
			assertBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp handler.JSONResponse[any]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.Equal(t, http.StatusNotFound, resp.Status)
				assert.Nil(t, resp.Data)
				require.NotNil(t, resp.Error)
				assert.Contains(t, resp.Error.Message, "not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRepoMock.ExpectedCalls = nil

			tt.setupMock()

			req := httptest.NewRequest(http.MethodDelete, "/projects/"+tt.projectID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPathValues(echo.PathValues{
				{Name: "id", Value: tt.projectID},
			})

			err := deleteHandler.HandleFunc(c)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			tt.assertBody(t, rec)

			projectRepoMock.AssertExpectations(t)
		})
	}
}
