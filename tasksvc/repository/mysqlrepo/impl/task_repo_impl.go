package impl

import (
	"context"
	"fmt"
	"tasksvc/models"
	"tasksvc/repository/mysqlrepo"
	"time"
)

const (
	createTask            string = "INSERT INTO tasks(title, description, priority, due_date) VALUES(?, ?, ?, ?)"
	updateTask            string = "UPDATE tasks SET title=?, description=?, priority=?, due_date=? WHERE id = ?"
	acceptTask            string = "INSERT INTO taskDetails(user_id, task_id) VALUES(?,?)"
	markComplete          string = "UPDATE taskDetails SET is_completed=true WHERE user_id = ? AND task_id = ?"
	CheckAlreadyCompleted string = "SELECT * FROM taskDetails WHERE user_id = ? AND task_id = ?"
	getTasks              string = "SELECT t.*, td.user_id AS user_id, td.task_id AS task_id, td.is_completed AS is_completed  FROM tasks t INNER JOIN taskDetails td ON td.task_id = t.id  WHERE td.user_id = ?"
	getAllTasks                  = "SELECT t.*, td.task_id AS task_id, td.user_id AS user_id, td.created_at AS detail_created_at,td.updated_at AS detail_updated_at, td.is_completed AS is_completed FROM tasks t LEFT JOIN taskdetails td ON t.id = td.task_id"
)

type TaskMySQLRepo struct {
	TaskDB *models.TaskDB
}

func NewTaskMySQLRepo(taskDB *models.TaskDB) mysqlrepo.TaskMySQLRepo {
	return &TaskMySQLRepo{
		taskDB,
	}
}

func (repo TaskMySQLRepo) CreateTask(ctx context.Context, task models.Task) error {
	parsedDueDate := task.DueDate.Format("2006-01-02 15:04:05")

	fmt.Println("INFO :: TaskMySQLRepo :: CreateTask :: Inside method", task.DueDate, parsedDueDate)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := repo.TaskDB.Db.PrepareContext(ctx, createTask)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: CreateTask :: Error %s when preparing SQL statement", err.Error()))
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, task.Title, task.Description, task.Priority, parsedDueDate)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: CreateTask :: Error %s while adding task", err.Error()))
		return err
	}
	return nil
}

func (repo TaskMySQLRepo) AcceptTask(ctx context.Context, task models.TaskDetails) error {
	fmt.Println("INFO :: TaskMySQLRepo :: AcceptTask :: Inside method")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := repo.TaskDB.Db.PrepareContext(ctx, acceptTask)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: CreateTask :: Error %s when preparing SQL statement", err.Error()))
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, task.UserID, task.TaskID)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: AcceptTask :: Error %s while accepting task", err.Error()))
		return err
	}
	return nil
}

func (repo TaskMySQLRepo) GetTasks(ctx context.Context, userId int64) ([]models.TaskResponse, error) {
	fmt.Println("INFO :: TaskMySQLRepo :: GetTasks :: Inside method")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tasks []models.TaskResponse

	err := repo.TaskDB.Db.SelectContext(ctx, &tasks, getTasks, userId)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: CreateTask :: Error %s when preparing SQL statement", err.Error()))
		return []models.TaskResponse{}, err
	}
	if len(tasks) == 0 {
		fmt.Println("INFO :: TaskMySQLRepo :: GetTasks :: no records found")
		return []models.TaskResponse{}, nil
	}
	return tasks, nil
}

func (repo TaskMySQLRepo) GetAllTasks(ctx context.Context, sortField, sortOrder, search string) ([]models.TaskResponse, error) {
	fmt.Println("INFO :: TaskMySQLRepo :: GetTasks :: Inside method")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tasks []models.TaskResponse

	query := getAllTasks
	var args []interface{}

	if search != "" {
		query += " WHERE t.title LIKE ? OR t.description LIKE ?"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	if sortField != "" {
		query += " ORDER BY " + sortField + " " + sortOrder
	}

	err := repo.TaskDB.Db.SelectContext(ctx, &tasks, query, args...)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: CreateTask :: Error %s when preparing SQL statement", err.Error()))
		return []models.TaskResponse{}, err
	}
	if len(tasks) == 0 {
		fmt.Println("INFO :: TaskMySQLRepo :: GetTasks :: no records found")
		return []models.TaskResponse{}, nil
	}
	return tasks, nil
}

func (repo TaskMySQLRepo) MarkComplete(ctx context.Context, task models.TaskDetails) error {
	fmt.Println("INFO :: TaskMySQLRepo :: MarkComplete :: Inside method")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := repo.TaskDB.Db.PrepareContext(ctx, markComplete)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: CreateTask :: Error %s when preparing SQL statement", err.Error()))
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, task.UserID, task.TaskID)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: AcceptTask :: Error %s while accepting task", err.Error()))
		return err
	}
	return nil
}

func (repo TaskMySQLRepo) CheckAlreadyCompleted(ctx context.Context, task models.TaskDetails) (bool, error) {
	fmt.Println("INFO :: TaskMySQLRepo :: CheckAlreadyCompleted :: Inside method")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tasks []models.TaskDetails

	err := repo.TaskDB.Db.SelectContext(ctx, &tasks, CheckAlreadyCompleted, task.UserID, task.TaskID)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: CheckAlreadyCompleted :: Error %s when preparing SQL statement", err.Error()))
		return false, err
	}
	if len(tasks) == 0 {
		fmt.Println("INFO :: TaskMySQLRepo :: GetTasks :: no records found")
		return false, nil
	}
	return tasks[0].IsCompleted, nil
}

func (repo TaskMySQLRepo) EditTask(ctx context.Context, taskId int64, task models.Task) error {
	parsedDueDate := task.DueDate.Format("2006-01-02 15:04:05")

	fmt.Println("INFO :: TaskMySQLRepo :: EditTask :: Inside method", task.DueDate, parsedDueDate)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := repo.TaskDB.Db.PrepareContext(ctx, updateTask)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: CreateTask :: Error %s when preparing SQL statement", err.Error()))
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, task.Title, task.Description, task.Priority, parsedDueDate, taskId)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: TaskMySQLRepo :: CreateTask :: Error %s while editing task", err.Error()))
		return err
	}
	return nil
}
