package http

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"tasksvc/constants"
	"tasksvc/models"
	"tasksvc/service"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService,
	}
}

func (handler *TaskHandler) CreateTask(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var taskReq models.Task
	if err := c.Bind(&taskReq); err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to parse the request :: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	err = handler.service.CreateTask(ctx, taskReq)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to create the task :: taskReq: %v :: %s", taskReq, err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, &models.StatusMessageResponse{Status: "success", Message: "Task created successfully"})
}

func (handler *TaskHandler) EditTask(c echo.Context) (err error) {
	ctx := c.Request().Context()
	taskId, _ := strconv.ParseInt(c.Param("taskId"), 10, 64)

	var taskReq models.Task
	if err := c.Bind(&taskReq); err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to parse the request :: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	err = handler.service.EditTask(ctx, taskId, taskReq)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to edit the task :: taskReq: %v :: %s", taskReq, err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.StatusMessageResponse{Status: "success", Message: "Task edited successfully"})
}

func (handler *TaskHandler) AcceptTask(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var taskReq models.TaskDetails
	if err := c.Bind(&taskReq); err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to parse the request :: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	err = handler.service.AcceptTask(ctx, taskReq)
	if err != nil {
		errMsg := err.Error()
		var mysqlErr *mysql.MySQLError
		fmt.Println(fmt.Sprintf("ERROR :: Failed to add the user :: %s", err.Error()))
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			errMsg = constants.ErrTaskAcceptedAlready
		}
		fmt.Println(fmt.Sprintf("ERROR :: Failed to accept the task :: taskReq: %v :: %s", taskReq, errMsg))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: errMsg})
	}
	return c.JSON(http.StatusOK, &models.StatusMessageResponse{Status: "success", Message: "Task accepted successfully"})
}

func (handler *TaskHandler) GetTasks(c echo.Context) (err error) {
	ctx := c.Request().Context()
	userId, _ := strconv.ParseInt(c.Param("userId"), 10, 64)
	tasks, err := handler.service.GetTasks(ctx, userId)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to get the tasks :: userId: %v :: Error: %s", userId, err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.StatusMessageResponse{Status: "success", Data: tasks})
}

func (handler *TaskHandler) GetAllTasks(c echo.Context) (err error) {
	ctx := c.Request().Context()
	sortField := c.Request().URL.Query().Get("sortField")
	sortOrder := c.Request().URL.Query().Get("sortOrder")
	search := c.Request().URL.Query().Get("search")
	tasks, err := handler.service.GetAllTasks(ctx, sortField, sortOrder, search)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to get all the tasks :: Error: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.StatusMessageResponse{Status: "success", Data: tasks})
}

func (handler *TaskHandler) MarkComplete(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var taskReq models.TaskDetails
	if err := c.Bind(&taskReq); err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to parse the request :: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	err = handler.service.MarkComplete(ctx, taskReq)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: Failed to mark as completed :: Error: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, &models.StatusMessageResponse{Status: "error", Message: err.Error()})
	}
	publishTaskCompletedEvent(taskReq.TaskID, taskReq.UserID)
	return c.JSON(http.StatusOK, &models.StatusMessageResponse{Status: "success", Message: "Marked the task as completed"})
}

func publishTaskCompletedEvent(taskId, userId int64) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %s", err)
	}
	defer producer.Close()

	message := &sarama.ProducerMessage{
		Topic: "task_completed",
		Value: sarama.StringEncoder(fmt.Sprintf("Task with task id %v has been completed by user with user id %v", taskId, userId)),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Printf("Error sending message to Kafka: %s\n", err)
	} else {
		log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	}
}
