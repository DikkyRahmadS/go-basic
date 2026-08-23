package config

import (
	"os"

	"github.com/gin-gonic/gin"

	"erp-internal/internal/middleware"
)

func NewGin() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	if os.Getenv("APP_ENV") == "development" {
		gin.SetMode(gin.DebugMode)
	}

	app := gin.New()

	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	app.Use(middleware.CORSMiddleware())

	return app
}
