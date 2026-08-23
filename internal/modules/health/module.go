package health

import (
	"gorm.io/gorm"

	"erp-internal/internal/modules/health/handler"
	"erp-internal/internal/modules/health/router"
	"erp-internal/internal/modules/health/service"
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
