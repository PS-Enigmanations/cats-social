package router_v1

import (
	"context"
	"enigmanations/cats-social/internal/user/controller"
	"enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/internal/user/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRouter struct {
	Controller controller.UserController
}

func NewUserRouter(ctx context.Context, pool *pgxpool.Pool) *UserRouter {
	userRepository := repository.NewUserRepository(pool)
	userAuthRepository := repository.NewUserAuthRepository(pool)

	userService := service.NewUserService(
		ctx,
		pool,
		&service.UserDependency{
			User:    userRepository,
			Session: userAuthRepository,
		},
	)
	userAuthService := service.NewUserAuthService(
		ctx,
		pool,
		&service.UserAuthDependency{
			User:    userRepository,
			Session: userAuthRepository,
		},
	)

	return &UserRouter{
		Controller: controller.NewUserController(userService, userAuthService),
	}
}
