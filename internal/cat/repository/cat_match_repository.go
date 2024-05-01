package repository

import "github.com/jackc/pgx/v5/pgxpool"

type CatMatchRepository interface {
}

type catMatchRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewCatMatchRepository(pool *pgxpool.Pool) CatMatchRepository {
	return &catMatchRepositoryDB{pool: pool}
}
