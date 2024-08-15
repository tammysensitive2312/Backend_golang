package server

import (
	"Backend_golang_project/infrastructure/config"
	"Backend_golang_project/internal/pkg"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func registerCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password_strength", pkg.PasswordValidator)
		v.RegisterValidation("valid_category", pkg.CategoryValidator)
		v.RegisterValidation("future_date", pkg.ProjectStartDateValidator)
	}
}

func NewGinEngine(lc fx.Lifecycle, config *config.Config) *gin.Engine {
	engine := gin.New()
	registerCustomValidators()
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting HTTP server")

			go func() {
				addr := fmt.Sprintf("%s:%s", config.HttpConfig.Host, config.HttpConfig.Port)
				if err := engine.Run(addr); err != nil {
					log.Fatal(fmt.Sprint("HTTP server failed to start %w", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server")
			return nil
		},
	})
	return engine
}
