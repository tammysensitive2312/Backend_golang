package router

import (
	"Backend_golang_project/internal/handlers"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RegisterRoutersIn struct {
	fx.In
	Engine         *gin.Engine
	ProjectHandler *handlers.ProjectHandler
}

func NewRegisterRouters(p RegisterRoutersIn) {
	r := p.Engine
	v1 := r.Group("/golang-web/api/")
	{
		projectGroup := v1.Group("projects")
		{
			projectGroup.POST("/create", p.ProjectHandler.Create)
			projectGroup.GET("/:id", p.ProjectHandler.GetById)
			projectGroup.DELETE("/:name", p.ProjectHandler.Delete)
			projectGroup.PUT("/:id", p.ProjectHandler.Update)
			projectGroup.GET("/all", p.ProjectHandler.GetProjects)
		}
	}
}
