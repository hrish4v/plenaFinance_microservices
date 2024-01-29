package impl

import (
	"context"
	"errors"
	"fmt"
	"tasksvc/constants"
	"tasksvc/models"
	"tasksvc/repository/mysqlrepo"
	"tasksvc/service"
)

type TaskService struct {
	repo mysqlrepo.TaskMySQLRepo
}

func NewTaskService(taskMySQLRepo mysqlrepo.TaskMySQLRepo) service.TaskService {
	return TaskService{
		taskMySQLRepo,
	}
}

func (service TaskService) CreateTask(ctx context.Context, task models.Task) error {
	return service.repo.CreateTask(ctx, task)
}

func (service TaskService) EditTask(ctx context.Context, taskId int64, task models.Task) error {
	return service.repo.EditTask(ctx, taskId, task)
}

func (service TaskService) AcceptTask(ctx context.Context, task models.TaskDetails) error {
	return service.repo.AcceptTask(ctx, task)
}

func (service TaskService) GetTasks(ctx context.Context, userId int64) ([]models.TaskResponse, error) {
	return service.repo.GetTasks(ctx, userId)
}

func (service TaskService) GetAllTasks(ctx context.Context, sortField, sortOrder, search string) ([]models.TaskResponse, error) {
	if sortField == constants.EmptyString {
		sortField = constants.DefaultSortField
	}
	if sortOrder == constants.EmptyString {
		sortOrder = constants.DefaultSortOrder
	}
	return service.repo.GetAllTasks(ctx, sortField, sortOrder, search)
}

func (service TaskService) MarkComplete(ctx context.Context, task models.TaskDetails) error {
	isAlreadyCompleted, err := service.repo.CheckAlreadyCompleted(ctx, task)
	if err != nil {
		fmt.Println("ERROR :: TaskService :: MarkComplete :: Error while checking already completed tasks", err.Error())
	}
	if isAlreadyCompleted {
		return errors.New(constants.ErrAlreadyMarkedCompleted)
	}
	return service.repo.MarkComplete(ctx, task)
}
