package config

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"erp-internal/internal/modules/health"
	"erp-internal/internal/modules/user"
)

type BootstrapConfig struct {
	App       *gin.Engine
	DB        *gorm.DB
	Validator *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	healthModule := health.NewModule(config.DB)
	healthModule.Router.RegisterRootRoutes(config.App)

	api := config.App.Group("/api")
	healthModule.Router.RegisterRoutes(api)

	userModule := user.NewModule(
		config.DB,
		config.Validator,
	)
	userModule.Router.RegisterRoutes(api)
}
