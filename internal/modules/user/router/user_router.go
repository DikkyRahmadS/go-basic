package router

import (
	"github.com/gin-gonic/gin"

	"erp-internal/internal/modules/user/handler"
)

type Router struct {
	handler *handler.UserHandler
}

func NewUserRouter(
	handler *handler.UserHandler,
) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) RegisterRoutes(
	api *gin.RouterGroup,
) {
	users := api.Group("/users")

	users.POST("", r.handler.Create)
	users.GET("", r.handler.FindAll)
	users.PUT("/:id", r.handler.Update)
	users.DELETE("/:id", r.handler.Delete)
}
