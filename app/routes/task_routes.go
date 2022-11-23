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
	jwtMiddleware middlewares.JWTMiddleware,
) TaskRoutes {
	return TaskRoutes{
		handler:        handler,
		taskController: taskController,
		jwtMiddleware:  jwtMiddleware,
	}
}

func (tr TaskRoutes) Setup() {
	log.Println("Setting up Tasks routes")
	// use JWT auth middleware for all task related routes
	authMiddleware := tr.jwtMiddleware.GinJWTMiddleware
	task := tr.handler.Group("/task")
	task.Use(authMiddleware.MiddlewareFunc())
	{
		task.POST("/create", tr.taskController.CreateTask)
		task.PUT("/update", tr.taskController.UpdateTask)
		task.DELETE("/delete/:uuid", tr.taskController.DeleteTask)
	}
}
