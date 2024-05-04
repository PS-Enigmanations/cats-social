package middleware

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Middleware struct {
	Auth AuthMiddleware
}

func RegisterMiddleware(ctx context.Context, pool *pgxpool.Pool) Middleware {
	return Middleware{
		Auth: NewAuthMiddleware(pool, ctx),
	}
}
