package repository

import (
	"PS-Enigmanations/cats-social/model/domain"
	"context"
	"database/sql"
)

type CatRepository interface {
	Save(ctx context.Context, tx *sql.Tx, cat domain.Cat) domain.Cat
	Update(ctx context.Context, tx *sql.Tx, cat domain.Cat) domain.Cat
	Delete(ctx context.Context, tx *sql.Tx, cat domain.Cat)
	FindById(ctx context.Context, tx *sql.Tx, catId int) domain.Cat
	Get(ctx context.Context, tx *sql.Tx) []domain.Cat
}
