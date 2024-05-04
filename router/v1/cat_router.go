package router_v1

import (
	"context"
	"enigmanations/cats-social/internal/cat/controller"
	"enigmanations/cats-social/internal/cat/repository"
	"enigmanations/cats-social/internal/cat/service"
	catMatchRepo "enigmanations/cats-social/internal/cat_match/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CatRouter struct {
	Controller controller.CatController
}

func NewCatRouter(ctx context.Context, pool *pgxpool.Pool) *CatRouter {
	catRepository := repository.NewCatRepository(pool)
	catMatchRepository := catMatchRepo.NewCatMatchRepository(pool)

	catService := service.NewCatService(
		ctx,
		pool,
		&service.CatDependency{
			Cat:      catRepository,
			CatMatch: catMatchRepository,
		},
	)

	return &CatRouter{
		Controller: controller.NewCatController(catService),
	}
}
