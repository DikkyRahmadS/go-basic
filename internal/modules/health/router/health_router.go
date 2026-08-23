package router

import (
	"github.com/gin-gonic/gin"

	"go-basic/internal/modules/health/handler"
)

type Router struct {
	handler *handler.HealthHandler
}

func NewHealthRouter(handler *handler.HealthHandler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/health", r.handler.Check)
}

func (r *Router) RegisterRootRoutes(app *gin.Engine) {
	app.GET("/health", r.handler.Check)
}
