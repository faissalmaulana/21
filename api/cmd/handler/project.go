package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/faissalmaulana/21/api/cmd/dto"
	constant "github.com/faissalmaulana/21/api/internal/const"
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
	Paginate model.Pagination  `json:"pagination"`
}

func (p *GetProjectsHandler) HandleFunc(c *echo.Context) error {
	var (
		search      string
		onlyArchive bool
		page        int
	)

	paramValues := c.QueryParams()

	if paramValues.Has("search") && paramValues.Get("search") != "" {
		search = paramValues.Get("search")
	}

	if paramValues.Has("archive") && paramValues.Get("archive") == "true" {
		onlyArchive = true
	}

	if v := paramValues.Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			page = p
		} else {
			page = 1
		}
	}

	// TODO: add validation
	projectParam := repository.ProjectsParam{
		Search:    search,
		IsArchive: onlyArchive,
		Size:      constant.PaginateSize,
		Page:      page,
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

type DeleteProjectHandler struct {
	ProjectRepository repository.ProjectRepository
}

func NewDeleteProjectHandler(pr repository.ProjectRepository) *DeleteProjectHandler {

	return &DeleteProjectHandler{
		ProjectRepository: pr,
	}
}

func (dp *DeleteProjectHandler) HandleFunc(c *echo.Context) error {
	id := strings.TrimPrefix(c.Param("id"), "/")

	deletedProjectId, err := dp.ProjectRepository.DeleteProjectByID(c.Request().Context(), id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return c.JSON(http.StatusNotFound, JSONResponse[any]{
				Status: http.StatusNotFound,
				Data:   nil,
				Error:  &ErrorResponse{Message: err.Error()},
			})
		default:
			return c.JSON(http.StatusInternalServerError, JSONResponse[any]{
				Status: http.StatusInternalServerError,
				Data:   nil,
				Error:  &ErrorResponse{Message: err.Error()},
			})
		}
	}

	return c.JSON(http.StatusOK, JSONResponse[string]{
		Status: http.StatusOK,
		Data:   fmt.Sprintf("Project with id %s deleted successfully", deletedProjectId),
		Error:  nil,
	})
}

type UpdateProjectHandler struct {
	Validator           *validator.Validate
	ValidatorSugaredMsg *service.SugaredErrorMessageValidator
	ProjectRepository   repository.ProjectRepository
}

func NewUpdateProjectHandler(pr repository.ProjectRepository, val *validator.Validate, sugaredErr *service.SugaredErrorMessageValidator) *UpdateProjectHandler {

	return &UpdateProjectHandler{
		ProjectRepository:   pr,
		Validator:           val,
		ValidatorSugaredMsg: sugaredErr,
	}
}

func (up *UpdateProjectHandler) HandleFunc(c *echo.Context) error {
	var (
		updateParams dto.UpdateProject
	)

	if err := c.Bind(&updateParams); err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse[any]{
			Status: http.StatusBadRequest,
			Data:   nil,
			Error:  &ErrorResponse{Message: err.Error()},
		})
	}

	if err := up.Validator.Struct(updateParams); err != nil {
		var errmsg strings.Builder

		errs := up.ValidatorSugaredMsg.TranslateValidationErrors(err)
		errName, ok := errs["name"]
		if ok {
			errmsg.WriteString(errName)
		}

		errArchived, ok := errs["to_be_archived"]
		if ok {
			errmsg.WriteString(",")
			errmsg.WriteString(errArchived)
		}

		return c.JSON(http.StatusBadRequest, JSONResponse[any]{
			Status: http.StatusBadRequest,
			Data:   nil,
			Error:  &ErrorResponse{Message: errmsg.String()},
		})
	}

	project, err := up.ProjectRepository.GetProjectByID(c.Request().Context(), strings.TrimPrefix(updateParams.ID, "/"))
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return c.JSON(http.StatusNotFound, JSONResponse[any]{
				Status: http.StatusNotFound,
				Data:   nil,
				Error:  &ErrorResponse{Message: err.Error()},
			})
		default:
			return c.JSON(http.StatusInternalServerError, JSONResponse[any]{
				Status: http.StatusInternalServerError,
				Data:   nil,
				Error:  &ErrorResponse{Message: err.Error()},
			})
		}
	}

	if updateParams.Name != nil {
		project.Name = *updateParams.Name
	}

	if updateParams.ToBeArchived != nil {
		project.IsArchive = updateParams.ToBeArchived
	}

	if updateParams.Name == nil && updateParams.ToBeArchived == nil {
		return c.JSON(http.StatusBadRequest, JSONResponse[any]{
			Status: http.StatusBadRequest,
			Data:   nil,
			Error:  &ErrorResponse{Message: "No fields provided to update"},
		})
	}

	if err := up.ProjectRepository.UpdateProject(c.Request().Context(), project); err != nil {
		return c.JSON(http.StatusInternalServerError, JSONResponse[any]{
			Status: http.StatusInternalServerError,
			Data:   nil,
			Error:  &ErrorResponse{Message: err.Error()},
		})
	}

	return c.JSON(http.StatusOK, JSONResponse[string]{
		Status: http.StatusOK,
		Data:   "Update Project Successfully",
		Error:  nil,
	})

}
