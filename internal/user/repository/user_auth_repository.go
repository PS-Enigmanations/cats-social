package repository

import (
	"context"
	"enigmanations/cats-social/internal/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAuthRepository interface {
	Save(ctx context.Context, model user.User) (*user.User, *user.UserSession, error)
}

type userAuthRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewUserAuthRepository(pool *pgxpool.Pool) UserAuthRepository {
	return &userAuthRepositoryDB{pool: pool}
}

func (db *userAuthRepositoryDB) Save(ctx context.Context, model user.User) (*user.User, *user.UserSession, error) {
	return nil, nil, nil
}
