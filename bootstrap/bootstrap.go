package bootstrap

import (
	"context"
	"log"

	"github.com/sayopaul/sendchamp-go-test/app/controllers"
	"github.com/sayopaul/sendchamp-go-test/app/middlewares"
	"github.com/sayopaul/sendchamp-go-test/app/routes"
	"github.com/sayopaul/sendchamp-go-test/config"
	"github.com/sayopaul/sendchamp-go-test/infrastructure"
	"github.com/sayopaul/sendchamp-go-test/repositories"
	"github.com/sayopaul/sendchamp-go-test/services"
	"go.uber.org/fx"
)

// bootstraps the application
var Bootstrap = fx.Options(
	repositories.Bootstrap,
	middlewares.Bootstrap,
	controllers.Bootstrap,
	routes.Bootstrap,
	services.Bootstrap,
	fx.Provide(config.LoadConfig),
	fx.Provide(infrastructure.NewDatabase),
	fx.Provide(infrastructure.NewRouter),
	fx.Invoke(Start),
)

func Start(
	lifecycle fx.Lifecycle,
	handler infrastructure.Router,
	routes routes.Routes,
	configEnv config.Config,
	database infrastructure.Database,
) {
	//actually start the application after injecting dependencies
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Println("Starting Application")
			log.Println("================================")
			log.Println("====== SENDCHAMP GO TEST - SAYOPAUL ======")
			log.Println("================================")

			go func() {
				routes.Setup()
				err := handler.Run(":" + configEnv.ServerPort)
				if err != nil {
					log.Panic("Error running server.")
				}

			}()
			return nil
		},
		OnStop: func(context.Context) error {

			log.Println("Stopping Application")

			db, _ := database.DB.DB()
			err := db.Close()

			if err != nil {
				log.Println("Database connection not closed")
			}

			return nil
		},
	})
}
