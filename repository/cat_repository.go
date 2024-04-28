package repository

import (
	"context"
	"database/sql"
	"enigmanations/cats-social/model/domain"
)

type CatRepository interface {
	Save(ctx context.Context, tx *sql.Tx, cat domain.Cat) domain.Cat
	Update(ctx context.Context, tx *sql.Tx, cat domain.Cat) domain.Cat
	Delete(ctx context.Context, tx *sql.Tx, cat domain.Cat)
	FindById(ctx context.Context, tx *sql.Tx, catId int) domain.Cat
	Get(ctx context.Context, tx *sql.Tx) []domain.Cat
}
