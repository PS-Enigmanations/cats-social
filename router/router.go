package middleware

import (
	"context"
	v1 "enigmanations/cats-social/router/v1"

	"enigmanations/cats-social/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRouter(ctx context.Context, pool *pgxpool.Pool, router *gin.Engine, m middleware.Middleware) {
	v1Route := v1.NewV1Router(ctx, pool)
	v1Route.Load(router, m)
}
