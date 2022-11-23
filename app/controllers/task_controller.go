package controllers

import (
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sayopaul/sendchamp-go-test/config"
	"github.com/sayopaul/sendchamp-go-test/models"
	"github.com/sayopaul/sendchamp-go-test/repositories"
	"github.com/sayopaul/sendchamp-go-test/services"
)

type TaskController struct {
	userRepository repositories.UserRepository
	taskRepository repositories.TaskRepository
	queueService   services.QueueService
	socketService  services.SocketService
	configEnv      config.Config
}

type createTaskInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updateTaskInfo struct {
	TaskUUID    string `json:"task_uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func NewTaskController(userRepository repositories.UserRepository, taskRepository repositories.TaskRepository, queueService services.QueueService, socketService services.SocketService, configEnv config.Config) TaskController {
	return TaskController{
		userRepository: userRepository,
		taskRepository: taskRepository,
		queueService:   queueService,
		socketService:  socketService,
		configEnv:      configEnv,
	}
}

func (tc TaskController) CreateTask(c *gin.Context) {
	// get ID of authenticated user
	claims := jwt.ExtractClaims(c)
	userID := uint(claims["id"].(float64))
	var createTask createTaskInfo
	if err := c.ShouldBindJSON(&createTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//publish to queue (rabbitMQ instance)
	errp := tc.publishTask(createTask)
	if errp != nil {
		log.Panic(errp.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errp.Error(),
			"data":    nil,
		})
		return
	}

	//insert task into the database
	task, err := tc.taskRepository.CreateTask(models.Task{
		UserID:      userID,
		Name:        createTask.Name,
		Description: createTask.Description,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	tc.socketService.Send(map[string]interface{}{
		"task_name":   task.Name,
		"description": task.Description,
	})
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task was created succesfully.",
		"data": gin.H{
			"task": task,
		},
	})

}

func (tc TaskController) publishTask(task createTaskInfo) error {
	err := tc.queueService.PublishMessage(map[string]interface{}{
		"task_name":   task.Name,
		"description": task.Description,
	}, tc.configEnv)
	return err
}

func (tc TaskController) UpdateTask(c *gin.Context) {
	// get ID of authenticated user
	claims := jwt.ExtractClaims(c)
	userID := uint(claims["id"].(float64))
	var updateTask updateTaskInfo
	if err := c.ShouldBindJSON(&updateTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//update the task
	taskCondition := models.Task{
		UserID: userID,
		UUID:   updateTask.TaskUUID,
	}
	updatedTask := models.Task{
		Name:        updateTask.Name,
		Description: updateTask.Description,
		Status:      updateTask.Status,
	}
	task, err := tc.taskRepository.UpdateTask(taskCondition, updatedTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "The task was updated succesfully.",
		"data": gin.H{
			"task": task,
		},
	})
}

func (tc TaskController) DeleteTask(c *gin.Context) {
	// validate signup details
	taskUUID := c.Param("uuid")
	err := tc.taskRepository.DeleteTask(models.Task{UUID: taskUUID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "there was a problem deleting the task, please try again" + err.Error(),
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task deleted successfully.",
		"data":    nil,
	})
}
