package main

import (
	"Backend_golang_project/infrastructure/config"
	infrastructure "Backend_golang_project/infrastructure/database"
	"Backend_golang_project/infrastructure/router"
	"Backend_golang_project/infrastructure/server"
	"Backend_golang_project/internal/handlers"
	"Backend_golang_project/internal/repositories"
	"Backend_golang_project/internal/use_cases"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func All() fx.Option {
	return fx.Options(
		//init
		fx.Provide(server.NewGinEngine),
		fx.Invoke(config.NewLogConfig),
		fx.Invoke(func(*gin.Engine) {}),
		fx.Provide(config.NewConfig),
		fx.Invoke(router.NewRegisterRouters),
		fx.Provide(infrastructure.NewInitDatabase),

		//inject repository
		fx.Provide(repositories.NewProjectRepository),
		fx.Provide(repositories.NewUserRepository),
		fx.Provide(repositories.NewS3Repository),

		//inject service
		fx.Provide(use_cases.NewProjectService),
		fx.Provide(use_cases.NewUserService),

		//inject controller
		fx.Provide(handlers.NewProjectHandler),
		fx.Provide(handlers.NewUserHandler),
	)
}
