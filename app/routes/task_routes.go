package routes

import (
	"log"

	"github.com/sayopaul/sendchamp-go-test/app/controllers"
	"github.com/sayopaul/sendchamp-go-test/app/middlewares"
	"github.com/sayopaul/sendchamp-go-test/infrastructure"
)

type TaskRoutes struct {
	handler        infrastructure.Router
	taskController controllers.TaskController
	jwtMiddleware  middlewares.JWTMiddleware
}

func NewTaskRoutes(
	handler infrastructure.Router,
	taskController controllers.TaskController,
) TaskRoutes {
	return TaskRoutes{
		handler:        handler,
		taskController: taskController,
	}
}

func (tr TaskRoutes) Setup() {
	log.Println("Setting up Tasks routes")
	task := tr.handler.Group("/task")
	// use JWT auth middleware for all task related routes
	authMiddleware := tr.jwtMiddleware.GinJWTMiddleware
	task.Use(authMiddleware.MiddlewareFunc())
	{
		task.POST("/create", tr.taskController.CreateTask)
		task.PUT("/update", tr.taskController.UpdateTask)
		task.DELETE("/delete/:uuid", tr.taskController.DeleteTask)
	}
}
