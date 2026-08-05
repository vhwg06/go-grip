package task

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Task capability.
type Handler struct {
	taskUC usecase.Task
	logger logger.Interface
}

// NewHandler constructs a new Task vertical handler instance.
func NewHandler(taskUC usecase.Task, l logger.Interface) *Handler {
	return &Handler{
		taskUC: taskUC,
		logger: l,
	}
}

func getActor(ctx context.Context) entity.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(entity.Actor); ok {
			return a
		}
	}
	return entity.Actor{}
}

// CreateTask handles POST /tasks
func (h *Handler) CreateTask(ctx context.Context, request openapi.CreateTaskRequestObject) (openapi.CreateTaskResponseObject, error) {
	if request.Body == nil {
		return openapi.CreateTask400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	desc := ""
	if request.Body.Description != nil {
		desc = *request.Body.Description
	}

	taskEntity, err := h.taskUC.Create(ctx, actor.UserID, request.Body.Title, desc)
	if err != nil {
		status, errResp := mapTaskError(err)
		switch status {
		case 400:
			return openapi.CreateTask400JSONResponse{}, nil
		case 401:
			return openapi.CreateTask401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.CreateTask500JSONResponse{}, nil
		}
	}

	taskDTO := toTaskResponse(taskEntity)
	return openapi.CreateTask201JSONResponse(taskDTO), nil
}

// ListTasks handles GET /tasks
func (h *Handler) ListTasks(ctx context.Context, request openapi.ListTasksRequestObject) (openapi.ListTasksResponseObject, error) {
	actor := getActor(ctx)
	limit := 10
	if request.Params.Limit != nil && *request.Params.Limit > 0 {
		limit = *request.Params.Limit
	}
	offset := 0
	if request.Params.Offset != nil && *request.Params.Offset >= 0 {
		offset = *request.Params.Offset
	}

	var statusFilter *entity.TaskStatus
	if request.Params.Status != nil && *request.Params.Status != "" {
		st := entity.TaskStatus(*request.Params.Status)
		statusFilter = &st
	}

	tasks, total, err := h.taskUC.List(ctx, actor.UserID, statusFilter, limit, offset)
	if err != nil {
		status, errResp := mapTaskError(err)
		switch status {
		case 401:
			return openapi.ListTasks401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.ListTasks500JSONResponse{}, nil
		}
	}

	listDTO := toTaskListResponse(tasks, total)
	return openapi.ListTasks200JSONResponse(listDTO), nil
}

// GetTaskByID handles GET /tasks/{id}
func (h *Handler) GetTaskByID(ctx context.Context, request openapi.GetTaskByIDRequestObject) (openapi.GetTaskByIDResponseObject, error) {
	actor := getActor(ctx)
	taskEntity, err := h.taskUC.Get(ctx, actor.UserID, request.Id)
	if err != nil {
		status, errResp := mapTaskError(err)
		switch status {
		case 401:
			return openapi.GetTaskByID401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.GetTaskByID403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.GetTaskByID404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetTaskByID500JSONResponse{}, nil
		}
	}

	taskDTO := toTaskResponse(taskEntity)
	return openapi.GetTaskByID200JSONResponse(taskDTO), nil
}

// UpdateTask handles PUT /tasks/{id}
func (h *Handler) UpdateTask(ctx context.Context, request openapi.UpdateTaskRequestObject) (openapi.UpdateTaskResponseObject, error) {
	if request.Body == nil {
		return openapi.UpdateTask400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	title := ""
	if request.Body.Title != nil {
		title = *request.Body.Title
	}
	desc := ""
	if request.Body.Description != nil {
		desc = *request.Body.Description
	}

	taskEntity, err := h.taskUC.Update(ctx, actor.UserID, request.Id, title, desc)
	if err != nil {
		status, errResp := mapTaskError(err)
		switch status {
		case 400:
			return openapi.UpdateTask400JSONResponse{}, nil
		case 401:
			return openapi.UpdateTask401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.UpdateTask403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.UpdateTask404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.UpdateTask500JSONResponse{}, nil
		}
	}

	taskDTO := toTaskResponse(taskEntity)
	return openapi.UpdateTask200JSONResponse(taskDTO), nil
}

// TransitionTask handles POST /tasks/{id}/transition
func (h *Handler) TransitionTask(ctx context.Context, request openapi.TransitionTaskRequestObject) (openapi.TransitionTaskResponseObject, error) {
	if request.Body == nil {
		return openapi.TransitionTask400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	newStatus := entity.TaskStatus(request.Body.Status)

	taskEntity, err := h.taskUC.Transition(ctx, actor.UserID, request.Id, newStatus)
	if err != nil {
		status, errResp := mapTaskError(err)
		switch status {
		case 400:
			return openapi.TransitionTask400JSONResponse{}, nil
		case 401:
			return openapi.TransitionTask401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.TransitionTask403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.TransitionTask404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		case 422:
			return openapi.TransitionTask422JSONResponse{
				UnprocessableEntityResponseJSONResponse: openapi.UnprocessableEntityResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.TransitionTask500JSONResponse{}, nil
		}
	}

	taskDTO := toTaskResponse(taskEntity)
	return openapi.TransitionTask200JSONResponse(taskDTO), nil
}

// DeleteTask handles DELETE /tasks/{id}
func (h *Handler) DeleteTask(ctx context.Context, request openapi.DeleteTaskRequestObject) (openapi.DeleteTaskResponseObject, error) {
	actor := getActor(ctx)
	err := h.taskUC.Delete(ctx, actor.UserID, request.Id)
	if err != nil {
		status, errResp := mapTaskError(err)
		switch status {
		case 401:
			return openapi.DeleteTask401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.DeleteTask403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.DeleteTask404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.DeleteTask500JSONResponse{}, nil
		}
	}

	return openapi.DeleteTask204Response{}, nil
}
