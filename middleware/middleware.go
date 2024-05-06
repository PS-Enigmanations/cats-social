package middleware

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Middleware struct {
	Auth AuthMiddleware
}

func NewMiddleware(pool *pgxpool.Pool) Middleware {
	return Middleware{
		Auth: NewAuthMiddleware(pool),
	}
}
