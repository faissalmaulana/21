package handler

import (
	"net/http"
	"strings"

	"github.com/faissalmaulana/21/api/cmd/dto"
	"github.com/faissalmaulana/21/api/internal/model"
	"github.com/faissalmaulana/21/api/internal/repository"
	"github.com/faissalmaulana/21/api/internal/service"
	"github.com/faissalmaulana/21/api/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type PostTaskHandler struct {
	TaskRepository      repository.TaskRepository
	Validator           *validator.Validate
	ValidatorSugaredMsg *service.SugaredErrorMessageValidator
}

func NewPostTaskHandler(tr repository.TaskRepository, val *validator.Validate, sgr *service.SugaredErrorMessageValidator) *PostTaskHandler {
	return &PostTaskHandler{
		TaskRepository:      tr,
		Validator:           val,
		ValidatorSugaredMsg: sgr,
	}
}

func (p *PostTaskHandler) HandleFunc(c *echo.Context) error {
	var postBadyPayload dto.PostTask

	if err := c.Bind(&postBadyPayload); err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse[any]{
			Status: http.StatusBadRequest,
			Data:   nil,
			Error:  &ErrorResponse{Message: err.Error()},
		})
	}

	if err := p.Validator.Struct(postBadyPayload); err != nil {
		var errmsg strings.Builder

		errs := p.ValidatorSugaredMsg.TranslateValidationErrors(err)
		errName, ok := errs["name"]
		if ok {
			errmsg.WriteString(errName)
		}

		errProjectID, ok := errs["project_id"]
		if ok {
			errmsg.WriteString(",")
			errmsg.WriteString(errProjectID)
		}

		errStartAt, ok := errs["start_at"]
		if ok {
			errmsg.WriteString(",")
			errmsg.WriteString(errStartAt)
		}

		return c.JSON(http.StatusBadRequest, JSONResponse[any]{
			Status: http.StatusBadRequest,
			Data:   nil,
			Error:  &ErrorResponse{Message: errmsg.String()},
		})
	}

	_, err := p.TaskRepository.AddTask(c.Request().Context(), model.Task{
		Name:      postBadyPayload.Name,
		ProjectID: &postBadyPayload.ProjectID,
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
		Data:   "Task Created",
		Error:  nil,
	})
}

type UpdateTaskHandler struct {
	TaskRepository      repository.TaskRepository
	Validator           *validator.Validate
	ValidatorSugaredMsg *service.SugaredErrorMessageValidator
}

func NewUpdateTaskHandler(tr repository.TaskRepository, val *validator.Validate, sgr *service.SugaredErrorMessageValidator) *UpdateTaskHandler {
	return &UpdateTaskHandler{
		TaskRepository:      tr,
		Validator:           val,
		ValidatorSugaredMsg: sgr,
	}
}

func (ut *UpdateTaskHandler) HandleFunc(c *echo.Context) error {
	var updateTaskPayload dto.UpdateTask
	if err := c.Bind(&updateTaskPayload); err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse[any]{
			Status: http.StatusBadRequest,
			Data:   nil,
			Error:  &ErrorResponse{Message: err.Error()},
		})
	}

	if err := ut.Validator.Struct(updateTaskPayload); err != nil {
		var errmsg strings.Builder

		errs := ut.ValidatorSugaredMsg.TranslateValidationErrors(err)
		errName, ok := errs["name"]
		if ok {
			errmsg.WriteString(errName)
		}

		errProjectID, ok := errs["project_id"]
		if ok {
			errmsg.WriteString(",")
			errmsg.WriteString(errProjectID)
		}

		errStartAt, ok := errs["start_at"]
		if ok {
			errmsg.WriteString(",")
			errmsg.WriteString(errStartAt)
		}

		errStatus, ok := errs["status"]
		if ok {
			errmsg.WriteString(",")
			errmsg.WriteString(errStatus)
		}

		return c.JSON(http.StatusBadRequest, JSONResponse[any]{
			Status: http.StatusBadRequest,
			Data:   nil,
			Error:  &ErrorResponse{Message: errmsg.String()},
		})
	}

	updatedTask := model.Task{}

	if updateTaskPayload.Name != "" {
		updatedTask.Name = updateTaskPayload.Name
	}

	if updateTaskPayload.ProjectID != "" {
		updatedTask.ProjectID = &updateTaskPayload.ProjectID
	}

	if !updateTaskPayload.StartAt.IsZero() {
		updatedTask.StartAt = &updateTaskPayload.StartAt
	}

	if updateTaskPayload.Status != "" {
		s := utils.ToStatus(updateTaskPayload.Status)
		updatedTask.Status = &s
	}

	updatedTaskID := strings.TrimPrefix(updateTaskPayload.ID, "/")
	if err := ut.TaskRepository.UpdateTask(c.Request().Context(), updatedTaskID, updatedTask); err != nil {
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
		Data:   "Update Task Successfully",
		Error:  nil,
	})
}

type GetTasksHandler struct {
	TaskRepository repository.TaskRepository
}

func NewGetTasksHandler(task repository.TaskRepository) *GetTasksHandler {

	return &GetTasksHandler{
		TaskRepository: task,
	}
}

type TaskResponse struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	ProjectID string       `json:"project_id"`
	Status    string       `json:"status"`
	StartAt   service.Date `json:"start_at"`
	Project   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
}

func (gt *GetTasksHandler) HandleFunc(c *echo.Context) error {
	tasks, err := gt.TaskRepository.Tasks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, JSONResponse[any]{
			Status: http.StatusInternalServerError,
			Data:   nil,
			Error:  &ErrorResponse{Message: err.Error()},
		})
	}

	tasksResponse := make([]TaskResponse, 0)

	if len(tasks) == 0 {
		return c.JSON(http.StatusOK, JSONResponse[[]TaskResponse]{
			Status: http.StatusOK,
			Data:   tasksResponse,
			Error:  nil,
		})
	}

	for _, task := range tasks {
		tasksResponse = append(tasksResponse, TaskResponse{
			ID:        task.ID,
			Name:      task.Name,
			ProjectID: *task.ProjectID,
			Status:    task.Status.String(),
			StartAt:   service.TimeToDate(*task.StartAt),
			Project: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   task.Project.ID,
				Name: task.Project.Name,
			},
		})
	}

	return c.JSON(http.StatusOK, JSONResponse[[]TaskResponse]{
		Status: http.StatusOK,
		Data:   tasksResponse,
		Error:  nil,
	})
}

type GetTaskByIDHandler struct {
	TaskRepository repository.TaskRepository
}

func NewGetTaskByIDHandler(task repository.TaskRepository) *GetTaskByIDHandler {
	return &GetTaskByIDHandler{
		TaskRepository: task,
	}
}

type TaskGetByIDResponse struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	ProjectID  string       `json:"project_id"`
	Status     string       `json:"status"`
	StartAt    service.Date `json:"start_at"`
	CreatedAt  service.Date `json:"created_at"`
	LastUpdate service.Date `json:"last_update"`
}

func (gt *GetTaskByIDHandler) HandleFunc(c *echo.Context) error {
	taskID := strings.TrimPrefix(c.Param("id"), "/")

	task, err := gt.TaskRepository.TaskByID(c.Request().Context(), taskID)
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

	return c.JSON(http.StatusOK, JSONResponse[TaskGetByIDResponse]{
		Status: http.StatusOK,
		Data: TaskGetByIDResponse{
			ID:         task.ID,
			Name:       task.Name,
			ProjectID:  *task.ProjectID,
			Status:     task.Status.String(),
			StartAt:    service.TimeToDate(*task.StartAt),
			CreatedAt:  service.TimeToDate(*task.CreatedAt),
			LastUpdate: service.TimeToDate(*task.LastUpdate),
		},
		Error: nil,
	})
}

type DeleteTaskByIDHandler struct {
	TaskRepository repository.TaskRepository
}

func NewDeleteTaskByIDHandler(task repository.TaskRepository) *DeleteTaskByIDHandler {
	return &DeleteTaskByIDHandler{
		TaskRepository: task,
	}
}

func (dt *DeleteTaskByIDHandler) HandleFunc(c *echo.Context) error {
	taskID := strings.TrimPrefix(c.Param("id"), "/")

	_, err := dt.TaskRepository.DeleteTaskByID(c.Request().Context(), taskID)
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
		Data:   "Task Deleted",
		Error:  nil,
	})
}
