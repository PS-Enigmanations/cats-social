package middleware

import (
	"context"
	v1 "enigmanations/cats-social/router/v1"

	"enigmanations/cats-social/middleware"

	"github.com/bmizerany/pat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func RegisterRouter(ctx context.Context, pool *pgxpool.Pool, router *pat.PatternServeMux, m middleware.Middleware) {
	v1Route := v1.NewV1Router(ctx, pool)
	v1Route.Load(router, m)
}
