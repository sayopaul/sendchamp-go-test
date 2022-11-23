package routes

import (
	"log"

	"github.com/sayopaul/sendchamp-go-test/app/controllers"
	"github.com/sayopaul/sendchamp-go-test/infrastructure"
)

type AuthRoutes struct {
	handler        infrastructure.Router
	authController controllers.AuthController
}

func NewAuthRoutes(
	handler infrastructure.Router,
	authController controllers.AuthController,
) AuthRoutes {
	return AuthRoutes{
		handler:        handler,
		authController: authController,
	}
}

func (ar AuthRoutes) Setup() {
	log.Println("Setting up Auth routes")
	auth := ar.handler.Group("/auth")
	{
		auth.POST("/signin", ar.authController.SignInHandler)
		auth.POST("/signup", ar.authController.SignupHandler)
	}
}
