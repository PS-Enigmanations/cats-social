package repository

import (
	"context"
	catmatch "enigmanations/cats-social/internal/cat_match"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"enigmanations/cats-social/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssociationByCatIdValue struct {
	Id             int
	UserId         int
	Name           string
	Race           string
	Sex            string
	AgeInMonth     int
	Description    string
	AlreadyMatched bool
}

type CatMatchRepository interface {
	Save(ctx context.Context, model *catmatch.CatMatch) error
	UpdateCatAlreadyMatches(ctx context.Context, ids []int, matched bool) error
	GetAssociationByCatId(ctx context.Context, id int) (*AssociationByCatIdValue, error)
}

type catMatchRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewCatMatchRepository(pool *pgxpool.Pool) CatMatchRepository {
	return &catMatchRepositoryDB{pool: pool}
}

func (db *catMatchRepositoryDB) Save(ctx context.Context, model *catmatch.CatMatch) error {
	const sql = `
		INSERT INTO cat_matches (issued_by, match_cat_id, user_cat_id, message, created_at)
		VALUES($1, $2, $3, $4, now());
	`

	if err := database.BeginTransaction(ctx, db.pool, func(tx pgx.Tx) error {
		_, err := tx.Exec(
			ctx,
			sql,
			model.IssuedBy,
			model.MatchCatId,
			model.UserCatId,
			model.Message,
		)
		if err != nil {
			log.Fatal("Cannot create cat match on database", slog.Any("error", err))
			return errors.New("Cannot create cat match on database")
		}

		return nil
	}); err != nil {
		return fmt.Errorf("Save transaction %w", err)
	}

	return nil
}

func (db *catMatchRepositoryDB) GetAssociationByCatId(ctx context.Context, id int) (*AssociationByCatIdValue, error) {
	const sql = `
		SELECT c.id, c.user_id, c.name, c.race, c.sex, c.age_in_month, c.description, c.is_already_matched
		FROM cats c WHERE c.id = $1 AND deleted_at IS NULL LIMIT 1;
	`
	row := db.pool.QueryRow(ctx, sql, id)
	v := new(AssociationByCatIdValue)
	err := row.Scan(
		&v.Id,
		&v.UserId,
		&v.Name,
		&v.Race,
		&v.Sex,
		&v.AgeInMonth,
		&v.Description,
		&v.AlreadyMatched,
	)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (db *catMatchRepositoryDB) UpdateCatAlreadyMatches(ctx context.Context, ids []int, matched bool) error {
	if err := database.BeginTransaction(ctx, db.pool, func(tx pgx.Tx) error {
		const sql = `
			UPDATE cats SET is_already_matched=@alreadyMatched WHERE id = @catId;
		`
		batch := &pgx.Batch{}
		for _, id := range ids {
			args := pgx.NamedArgs{
				"catId":          id,
				"alreadyMatched": matched,
			}
			batch.Queue(sql, args)
		}

		results := tx.SendBatch(ctx, batch)
		defer results.Close()

		for range ids {
			_, err := results.Exec()
			if err != nil {
				return fmt.Errorf("Unable to update row: %w", err)
			}
		}

		return results.Close()
	}); err != nil {
		return fmt.Errorf("Update transaction %w", err)
	}

	return nil
}
