package user

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"go-basic/internal/modules/user/handler"
	"go-basic/internal/modules/user/repository"
	"go-basic/internal/modules/user/router"
	"go-basic/internal/modules/user/service"
)

type Module struct {
	Router *router.Router
}

func NewModule(
	db *gorm.DB,
	validate *validator.Validate,
) *Module {
	userRepository := repository.NewUserRepository(db)

	userService := service.NewUserService(
		db,
		userRepository,
		validate,
	)

	userHandler := handler.NewUserHandler(
		userService,
	)

	userRouter := router.NewUserRouter(
		userHandler,
	)

	return &Module{
		Router: userRouter,
	}
}
