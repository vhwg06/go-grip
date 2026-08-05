package task

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// toTaskResponse maps a domain entity.Task to an openapi.TaskResponse DTO.
func toTaskResponse(t entity.Task) openapi.TaskResponse {
	desc := t.Description
	statusStr := string(t.Status)

	return openapi.TaskResponse{
		Id:          t.ID,
		Title:       t.Title,
		Description: &desc,
		Status:      statusStr,
		UserId:      t.UserID,
		CreatedAt:   &t.CreatedAt,
		UpdatedAt:   &t.UpdatedAt,
	}
}

// toTaskListResponse maps a slice of domain entity.Task and total count to openapi.TaskListResponse DTO.
func toTaskListResponse(tasks []entity.Task, total int) openapi.TaskListResponse {
	items := make([]openapi.TaskResponse, len(tasks))
	for i, t := range tasks {
		items[i] = toTaskResponse(t)
	}

	return openapi.TaskListResponse{
		Items: items,
		Total: total,
	}
}
