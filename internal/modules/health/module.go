package health

import (
	"gorm.io/gorm"

	"go-basic/internal/modules/health/handler"
	"go-basic/internal/modules/health/router"
	"go-basic/internal/modules/health/service"
)

type Module struct {
	Router *router.Router
}

func NewModule(db *gorm.DB) *Module {
	healthService := service.NewHealthService(db)
	healthHandler := handler.NewHealthHandler(healthService)
	healthRouter := router.NewHealthRouter(healthHandler)

	return &Module{
		Router: healthRouter,
	}
}
