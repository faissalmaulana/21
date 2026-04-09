package handler

import (
	"net/http"

	"github.com/faissalmaulana/21/api/cmd/dto"
	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/repository"
	"github.com/faissalmaulana/21/api/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type GetProjectsHandler struct {
	ProjectRepository repository.ProjectRepository
}

func NewGetProjectsHandler(pr repository.ProjectRepository) *GetProjectsHandler {

	return &GetProjectsHandler{
		ProjectRepository: pr,
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type JSONResponse[T any] struct {
	Status int            `json:"status"`
	Data   T              `json:"data"`
	Error  *ErrorResponse `json:"error"`
}

type ProjectResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetProjectsResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Paginate model.Pagination  `json:"paginate"`
}

func (p *GetProjectsHandler) HandleFunc(c *echo.Context) error {
	var (
		search      string
		onlyArchive bool
	)

	paramValues := c.QueryParams()

	if paramValues.Has("search") && paramValues.Get("search") != "" {
		search = paramValues.Get("search")
	}

	if paramValues.Has("archive") && paramValues.Get("archive") == "true" {
		onlyArchive = true
	}

	// TODO: add validation
	projectParam := repository.ProjectsParam{
		Search:    search,
		IsArchive: onlyArchive,
	}

	projects, paginate, err := p.ProjectRepository.Projects(c.Request().Context(), projectParam)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			JSONResponse[any]{
				Status: http.StatusInternalServerError,
				Data:   nil,
				Error: &ErrorResponse{
					Message: err.Error(),
				},
			},
		)
	}

	projectsResponse := make([]ProjectResponse, 0)

	for _, project := range projects {
		projectsResponse = append(projectsResponse, ProjectResponse{
			ID:   project.ID,
			Name: project.Name,
		})
	}

	return c.JSON(http.StatusOK, JSONResponse[GetProjectsResponse]{
		Status: http.StatusOK,
		Data: GetProjectsResponse{
			Projects: projectsResponse,
			Paginate: paginate,
		},
		Error: nil,
	})
}

type PostProjectHandler struct {
	ProjectRepository   repository.ProjectRepository
	Validator           *validator.Validate
	ValidatorSugaredMsg *service.SugaredErrorMessageValidator
}

func NewPostProjectHandler(
	pr repository.ProjectRepository,
	sugaredErr *service.SugaredErrorMessageValidator,
	validator *validator.Validate,
) *PostProjectHandler {
	return &PostProjectHandler{
		ProjectRepository:   pr,
		ValidatorSugaredMsg: sugaredErr,
		Validator:           validator,
	}
}

func (pp *PostProjectHandler) HandleFunc(c *echo.Context) error {
	var newProjectPayload dto.PostProject

	if err := c.Bind(&newProjectPayload); err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse[any]{
			Status: http.StatusBadRequest,
			Data:   nil,
			Error:  &ErrorResponse{Message: err.Error()},
		})
	}

	if err := pp.Validator.Struct(newProjectPayload); err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse[any]{
			Status: http.StatusBadRequest,
			Data:   nil,
			Error:  &ErrorResponse{Message: pp.ValidatorSugaredMsg.TranslateValidationErrors(err)["name"]},
		})
	}

	_, err := pp.ProjectRepository.AddProject(c.Request().Context(), model.Project{
		Name: newProjectPayload.Name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, JSONResponse[any]{
			Status: http.StatusInternalServerError,
			Data:   nil,
			Error:  &ErrorResponse{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, JSONResponse[string]{
		Status: http.StatusCreated,
		Data:   "Add New Project Successfully",
		Error:  nil,
	})
}
