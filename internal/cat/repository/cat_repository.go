package repository

import (
	"context"
	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/internal/cat/request"
	"enigmanations/cats-social/util"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatRepository interface {
	GetAllByParams(ctx context.Context, params *request.CatGetAllQueryParams, ownerId int) ([]*cat.Cat, error)
	FindById(ctx context.Context, catId int) (*cat.Cat, error)
	Save(ctx context.Context, tx pgx.Tx, model cat.Cat) (*cat.Cat, error)
	Update(ctx context.Context, tx pgx.Tx, model cat.Cat) (*cat.Cat, error)
	Delete(ctx context.Context, tx pgx.Tx, catId int) error
}

type Database struct {
	pool *pgxpool.Pool
}

func NewCatRepository(pool *pgxpool.Pool) CatRepository {
	return &Database{
		pool: pool,
	}
}

func (db *Database) GetAllByParams(ctx context.Context, params *request.CatGetAllQueryParams, ownerId int) ([]*cat.Cat, error) {
	var (
		args  []any
		where []string
	)

	sql := fmt.Sprintf(`
		SELECT
			c.id,
			c.name,
			c.race,
			c.sex,
			c.age_in_month,
			c.description,
			c.has_matched
		FROM cats c
		`)

	// Id
	if params.Id != "" {
		args = append(args, params.Id)
		where = append(where, fmt.Sprintf(`"id" = $%d`, len(args)))
	}
	// Race
	if params.Race != "" {
		args = append(args, params.Race)
		where = append(where, fmt.Sprintf(`"race" = $%d`, len(args)))
	}
	// Sex
	if params.Sex != "" {
		args = append(args, params.Sex)
		where = append(where, fmt.Sprintf(`"sex" = $%d`, len(args)))
	}
	// HasMatched
	if params.HasMatched != "" {
		hasMethod, err := strconv.ParseBool(params.HasMatched)
		if nil != err {
			return nil, err
		}

		args = append(args, hasMethod)
		where = append(where, fmt.Sprintf(`"has_matched" = $%d`, len(args)))
	}
	// AgeInMonth
	if params.AgeInMonth != "" {
		// Parse the input value
		ageOperator, err := util.ParseQueryOperator(params.AgeInMonth)
		if err != nil {
			return nil, err
		}

		args = append(args, ageOperator.Value)
		where = append(where, fmt.Sprintf(`"age_in_month" %s $%d`, ageOperator.Operator, len(args)))
	}
	// Owned
	if params.Owned != "" {
		args = append(args, ownerId)
		if params.Owned == "true" {
			where = append(where, fmt.Sprintf(`"user_id" = $%d`, len(args)))
		} else {
			where = append(where, fmt.Sprintf(`"user_id" != $%d`, len(args)))
		}
	}
	// Search
	if params.Search != "" {
		args = append(args, params.Search)
		where = append(where, fmt.Sprintf(`"name" LIKE $%d`, len(args)))
	}

	// Merge where clauses
	if len(where) > 0 {
		w := " WHERE " + strings.Join(where, " AND ") + " AND deleted_at IS NULL" // #nosec G202
		sql += w
	} else {
		w := " WHERE deleted_at IS NULL"
		sql += w
	}

	// Limit (default: 5)
	if params.Limit != "" {
		sql += fmt.Sprintf(` LIMIT %s`, params.Limit)
	} else {
		sql += fmt.Sprintf(` LIMIT %d`, 5)
	}
	// Offset (default: 0)
	if params.Offset != "" {
		sql += fmt.Sprintf(` OFFSET %s`, params.Offset)
	} else {
		sql += fmt.Sprintf(` OFFSET %d`, 0)
	}

	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	// close rows if error ocur
	defer rows.Close()

	// iterate Rows
	var cats []*cat.Cat
	if rows != nil {
		for rows.Next() {
			// create 'c' for struct 'Cat'
			c := new(cat.Cat)

			// scan rows and place it in 'c' (cat) container
			err := rows.Scan(
				&c.Id,
				&c.Name,
				&c.Race,
				&c.Sex,
				&c.AgeInMonth,
				&c.Description,
				&c.HasMatched,
			)

			// return nil and error if scan operation fail
			if err != nil {
				return nil, err
			}

			// add c to cats slice
			cats = append(cats, c)
		}
	}

	// return cats slice and nil for the error
	return cats, nil
}

func (db *Database) FindById(ctx context.Context, catId int) (*cat.Cat, error) {
	const catQuery = `
		SELECT id, name, race, sex, age_in_month, description from cats WHERE id = $1;
	`
	row := db.pool.QueryRow(ctx, catQuery, catId)

	c := new(cat.Cat)
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.Race,
		&c.Sex,
		&c.AgeInMonth,
		&c.Description,
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (db *Database) Save(ctx context.Context, tx pgx.Tx, model cat.Cat) (*cat.Cat, error) {
	const sql = `INSERT into cats
		("user_id", "name", "race", "sex", "age_in_month", "description")
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, race, sex, age_in_month, description;`

	// execute query to insert new record. it takes 'cat' variable as its input
	// the result will be placed in 'row' variable
	row := tx.QueryRow(
		ctx,
		sql,
		model.UserId,
		model.Name,
		model.Race,
		model.Sex,
		model.AgeInMonth,
		model.Description,
	)

	// create 'c' variable as 'Cat' type to contain scanned data value from 'row' variable
	c := new(cat.Cat)

	// scan 'row' variable and place the value to 'c' variable as well as check for error
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.Race,
		&c.Sex,
		&c.AgeInMonth,
		&c.Description,
	)

	// return nil and error if scan operation is fail/ error found
	if err != nil {
		return nil, fmt.Errorf("Save %w", err)
	}

	return c, nil
}

func (db *Database) Update(ctx context.Context, tx pgx.Tx, model cat.Cat) (*cat.Cat, error) {
	const sql = `UPDATE cats
		SET name = COALESCE($2, "name"),
			race = COALESCE($3, "race"),
			sex = COALESCE($4, "sex"),
			age_in_month = COALESCE($5, "age_in_month"),
			description = COALESCE($6, "description"),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, race, sex, age_in_month, description;`

	row := tx.QueryRow(
		ctx,
		sql,
		model.Id,
		model.Name,
		model.Race,
		model.Sex,
		model.AgeInMonth,
		model.Description,
	)

	c := new(cat.Cat)
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.Race,
		&c.Sex,
		&c.AgeInMonth,
		&c.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("Update %w", err)
	}

	return c, nil
}

func (db *Database) Delete(ctx context.Context, tx pgx.Tx, catId int) error {
	const sql = `UPDATE cats SET deleted_at = NOW() WHERE id = $1`

	_, err := tx.Exec(ctx, sql, catId)
	if err != nil {
		return fmt.Errorf("Delete %w", err)
	}

	return nil
}
