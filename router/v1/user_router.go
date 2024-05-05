package router_v1

import (
	"context"
	sessionRepository "enigmanations/cats-social/internal/session/repository"
	sessionService "enigmanations/cats-social/internal/session/service"
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
	userSessionRepository := sessionRepository.NewUserSessionRepository(pool)

	userService := service.NewUserService(
		ctx,
		pool,
		&service.UserDependency{
			User:    userRepository,
			Session: userSessionRepository,
		},
	)
	userSessionService := sessionService.NewSessionService(
		ctx,
		pool,
		&sessionService.SessionDependency{
			User:    userRepository,
			Session: userSessionRepository,
		},
	)

	return &UserRouter{
		Controller: controller.NewUserController(userService, userSessionService),
	}
}
