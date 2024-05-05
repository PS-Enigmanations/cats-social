package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatImageRepository interface {
	SaveImageUrls(ctx context.Context, tx pgx.Tx, catId int, urls []string) error
	DeleteImageUrls(ctx context.Context, tx pgx.Tx, catId int) error
}

type Database struct {
	pool *pgxpool.Pool
}

func NewCatImageRepository(pool *pgxpool.Pool) CatImageRepository {
	return &Database{
		pool: pool,
	}
}

func (db *Database) SaveImageUrls(ctx context.Context, tx pgx.Tx, catId int, urls []string) error {
	const sql = `INSERT into cat_images
		("cat_id", "url")
		VALUES ($1, $2);`

	for _, url := range urls {
		_, err := tx.Exec(ctx, sql, catId, url)
		if err != nil {
			return fmt.Errorf("SaveImageUrls %w", err)
		}
	}

	return nil
}

func (db *Database) DeleteImageUrls(ctx context.Context, tx pgx.Tx, catId int) error {
	const sql = `DELETE FROM cat_images WHERE cat_id = $1`

	_, err := tx.Exec(ctx, sql, catId)
	if err != nil {
		return fmt.Errorf("Delete Image Urls %w", err)
	}

	return nil
}
