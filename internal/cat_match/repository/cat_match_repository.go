package repository

import (
	"context"
	catmatch "enigmanations/cats-social/internal/cat_match"
	"errors"
	"fmt"
	"log"
	"log/slog"

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
	Save(ctx context.Context, model *catmatch.CatMatch, tx pgx.Tx) error
	UpdateCatAlreadyMatches(ctx context.Context, ids []int, matched bool, tx pgx.Tx) error
	UpdateCatMatchStatus(ctx context.Context, id int, status string, tx pgx.Tx) error
	Get(ctx context.Context, id int) (*catmatch.CatMatch, error)
	GetByCatId(ctx context.Context, id int) (*catmatch.CatMatch, error)
	GetAssociationByCatId(ctx context.Context, id int) (*AssociationByCatIdValue, error)
	GetByIssuedOrReceiver(ctx context.Context, id int) ([]*catmatch.CatMatchValue, error)
	Destroy(ctx context.Context, id int64, tx pgx.Tx) error
}

type catMatchRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewCatMatchRepository(pool *pgxpool.Pool) CatMatchRepository {
	return &catMatchRepositoryDB{pool: pool}
}

func (db *catMatchRepositoryDB) Save(ctx context.Context, model *catmatch.CatMatch, tx pgx.Tx) error {
	const sql = `
		INSERT INTO cat_matches (issued_by, match_cat_id, user_cat_id, message, created_at)
		VALUES($1, $2, $3, $4, now());
	`
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
}

func (db *catMatchRepositoryDB) Get(ctx context.Context, id int) (*catmatch.CatMatch, error) {
	const sql = `SELECT id, match_cat_id, issued_by, user_cat_id, message, status, created_at
	FROM cat_matches
	WHERE id = $1 AND deleted_at IS NULL LIMIT 1;`

	row := db.pool.QueryRow(ctx, sql, id)
	v := new(catmatch.CatMatch)

	err := row.Scan(
		&v.Id,
		&v.MatchCatId,
		&v.IssuedBy,
		&v.UserCatId,
		&v.Message,
		&v.Status,
		&v.CreatedAt,
	)
	if err != nil {
		log.Print("Error getting cat match", err)
		return nil, err
	}

	return v, nil
}

func (db *catMatchRepositoryDB) GetAssociationByCatId(ctx context.Context, id int) (*AssociationByCatIdValue, error) {
	const sql = `
		SELECT c.id, c.user_id, c.name, c.race, c.sex, c.age_in_month, c.description, c.has_matched
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

func (db *catMatchRepositoryDB) GetByCatId(ctx context.Context, catId int) (*catmatch.CatMatch, error) {
	const q = `SELECT id, match_cat_id, issued_by, user_cat_id, message, status, created_at
	FROM cat_matches
	WHERE (match_cat_id = COALESCE($1, "match_cat_id") OR user_cat_id = COALESCE($1, "user_cat_id"))
	AND deleted_at IS NULL LIMIT 1;`

	row := db.pool.QueryRow(ctx, q, catId)
	v := new(catmatch.CatMatch)

	err := row.Scan(
		&v.Id,
		&v.MatchCatId,
		&v.IssuedBy,
		&v.UserCatId,
		&v.Message,
		&v.Status,
		&v.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Print("Error getting cat match", err)
		return nil, err
	}

	return v, nil
}

func (db *catMatchRepositoryDB) UpdateCatAlreadyMatches(ctx context.Context, ids []int, matched bool, tx pgx.Tx) error {
	const sql = `
		UPDATE cats SET
			has_matched=@alreadyMatched,
			updated_at = NOW()
		WHERE id = @catId;
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
}

func (db *catMatchRepositoryDB) Destroy(ctx context.Context, id int64, tx pgx.Tx) error {
	const sql = `
		UPDATE cat_matches SET deleted_at=now() WHERE id = $1
	`
	_, err := tx.Exec(
		ctx,
		sql,
		id,
	)

	if err != nil {
		log.Fatal("Cannot delete cat match on database", slog.Any("error", err))
		return errors.New("Cannot delete cat match on database")
	}

	return nil
}

func (db *catMatchRepositoryDB) UpdateCatMatchStatus(ctx context.Context, id int, status string, tx pgx.Tx) error {
	const sql = `
		UPDATE cat_matches SET
			status=$1,
			updated_at = NOW()
		WHERE id = $2
	`
	_, err := tx.Exec(
		ctx,
		sql,
		status,
		id,
	)
	if err != nil {
		log.Fatal("Cannot update cat match status on database", slog.Any("error", err))
		return errors.New("Cannot update cat match status on database")
	}

	return nil
}

func (db *catMatchRepositoryDB) GetByIssuedOrReceiver(ctx context.Context, id int) ([]*catmatch.CatMatchValue, error) {
	const sql = `
	SELECT
		cm.id AS cat_match_id,
		cm.message,
		cm.status,
		cm.created_at AS match_created_at,
		u.name AS user_name,
		u.email AS user_email,
		u.created_at AS user_created_at,
		uc.id AS user_cat_id,
		uc.user_id AS user_cat_user_id,
		uc.name AS user_cat_name,
		uc.race AS user_cat_race,
		uc.sex AS user_cat_sex,
		uc.age_in_month AS user_cat_age_in_month,
		uc.description AS user_cat_description,
		uc.has_matched AS user_cat_has_matched,
		uc.created_at AS user_cat_created_at,
		mc.id AS match_cat_id,
		mc.user_id AS match_cat_user_id,
		mc.name AS match_cat_name,
		mc.race AS match_cat_race,
		mc.sex AS match_cat_sex,
		mc.age_in_month AS match_cat_age_in_month,
		mc.description AS match_cat_description,
		mc.has_matched AS match_cat_has_matched,
		mc.created_at AS match_cat_created_at
	FROM
		cat_matches AS cm
	JOIN
		users AS u ON cm.issued_by = u.id
	JOIN
		cats AS uc ON cm.user_cat_id = uc.id
	JOIN
		cats AS mc ON cm.match_cat_id = mc.id
	LEFT JOIN
		cat_images AS uci ON uc.id = uci.cat_id
	LEFT JOIN
		cat_images AS mci ON mc.id = mci.cat_id
	WHERE
		(cm.issued_by = $1 OR mc.user_id = $1)
		AND cm.deleted_at IS NULL
	GROUP BY
		cm.id,
		u.id,
		uc.id,
		mc.id;

	`
	// execute query
	rows, err := db.pool.Query(ctx, sql, id)

	if err != nil {
		return nil, err
	}

	// close rows if error ocur
	defer rows.Close()
	// iterate Rows
	var cms []*catmatch.CatMatchValue
	if rows != nil {
		for rows.Next() {
			// 		// create 'cm' for struct 'CatMatch'
			cm := new(catmatch.CatMatchValue)

			// scan rows and place it in 'cm' (cat match) container
			err := rows.Scan(
				&cm.CatMatchId,
				&cm.Message,
				&cm.Status,
				&cm.MatchCreatedAt,
				&cm.UserName,
				&cm.UserEmail,
				&cm.UserCreatedAt,
				&cm.UserCatId,
				&cm.UserCatUserId,
				&cm.UserCatName,
				&cm.UserCatRace,
				&cm.UserCatSex,
				&cm.UserCatAgeInMonth,
				&cm.UserCatDescription,
				&cm.UserCatHasMatched,
				&cm.UserCatCreatedAt,
				&cm.MatchCatId,
				&cm.MatchCatUserId,
				&cm.MatchCatName,
				&cm.MatchCatRace,
				&cm.MatchCatSex,
				&cm.MatchCatAgeInMonth,
				&cm.MatchCatDescription,
				&cm.MatchCatHasMatched,
				&cm.MatchCatCreatedAt,
				// &cm.UserCatImageUrls,
				// &cm.MatchCatImageUrls,
			)

			// return nil and error if scan operation fail
			if err != nil {
				return nil, err
			}

			// add c to cats slice
			cms = append(cms, cm)
		}
	}

	// return cats slice and nil for the error
	return cms, nil
}
