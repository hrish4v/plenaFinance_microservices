package mysqlrepo

import (
	"context"
	"tasksvc/models"
)

type TaskMySQLRepo interface {
	CreateTask(ctx context.Context, task models.Task) error
	EditTask(ctx context.Context, taskId int64, task models.Task) error
	AcceptTask(ctx context.Context, task models.TaskDetails) error
	GetTasks(ctx context.Context, userId int64) ([]models.TaskResponse, error)
	GetAllTasks(ctx context.Context, sortField, sortOrder, search string) ([]models.TaskResponse, error)
	MarkComplete(ctx context.Context, task models.TaskDetails) error
	CheckAlreadyCompleted(ctx context.Context, task models.TaskDetails) (bool, error)
}
