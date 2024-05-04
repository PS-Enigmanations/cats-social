package router_v1

import (
	"context"
	"enigmanations/cats-social/internal/cat_match/controller"
	"enigmanations/cats-social/internal/cat_match/repository"
	"enigmanations/cats-social/internal/cat_match/service"
	userRepository "enigmanations/cats-social/internal/user/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CatMatchRouter struct {
	Controller controller.CatMatchController
}

func NewCatMatchRouter(ctx context.Context, pool *pgxpool.Pool) *CatMatchRouter {
	userRepository := userRepository.NewUserRepository(pool)

	catMatchRepository := repository.NewCatMatchRepository(pool)
	catMatchService := service.NewCatMatchService(
		ctx,
		pool,
		&service.CatMatchDependency{
			User:     userRepository,
			CatMatch: catMatchRepository,
		},
	)

	return &CatMatchRouter{
		Controller: controller.NewCatMatchController(catMatchService),
	}
}
