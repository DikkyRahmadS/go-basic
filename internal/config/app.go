package config

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-basic/docs"
	"go-basic/internal/modules/health"
	"go-basic/internal/modules/user"
)

type BootstrapConfig struct {
	App       *gin.Engine
	DB        *gorm.DB
	Validator *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	healthModule := health.NewModule(config.DB)
	healthModule.Router.RegisterRootRoutes(config.App)

	config.App.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := config.App.Group("/api")
	healthModule.Router.RegisterRoutes(api)

	userModule := user.NewModule(
		config.DB,
		config.Validator,
	)
	userModule.Router.RegisterRoutes(api)
}
