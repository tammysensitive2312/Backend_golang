package router

import (
	"Backend_golang_project/infrastructure/config"
	"Backend_golang_project/infrastructure/middleware"
	"Backend_golang_project/infrastructure/middleware/jwt"
	"Backend_golang_project/internal/handlers"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RegisterRoutersIn struct {
	fx.In
	Engine         *gin.Engine
	ProjectHandler *handlers.ProjectHandler
	UserHandler    *handlers.UserHandler
	Config         *config.Config
}

func NewRegisterRouters(p RegisterRoutersIn) {
	r := p.Engine
	v1 := r.Group("/golang-web/api/")
	v1.Use(middleware.LoggingMiddleware(), middleware.GinRecovery(true))
	{
		v1.POST("/refresh", p.UserHandler.RefreshToken)
		projectGroup := v1.Group("projects")
		projectGroup.Use(jwt.AuthMiddleware(p.Config))
		{
			projectGroup.POST("/create", p.ProjectHandler.Create)
			projectGroup.GET("/:id", p.ProjectHandler.GetById)
			projectGroup.DELETE("/:name", p.ProjectHandler.Delete)
			projectGroup.PUT("/:id", p.ProjectHandler.Update)
			projectGroup.GET("/all", p.ProjectHandler.GetProjects)
		}

		userGroup := v1.Group("users")
		{
			userGroup.POST("/create", p.UserHandler.CreateNewUser)
			userGroup.POST("/login", p.UserHandler.LoginUser)
			userGroup.GET("/:id/projects", p.UserHandler.GetUserById)
		}

	}
}
