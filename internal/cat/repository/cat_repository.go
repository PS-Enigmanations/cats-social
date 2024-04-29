package repository

import (
	"context"
	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/pkg/database"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ICatRepository interface {
	GetAll(ctx context.Context) ([]*cat.Cat, error)
	Save(ctx context.Context, model cat.Cat) (*cat.Cat, error)
}

type Database struct {
	pool *pgxpool.Pool
}

func NewCatRepository(pool *pgxpool.Pool) Database {
	return Database{
		pool: pool,
	}
}

func (db *Database) GetAll(ctx context.Context) ([]*cat.Cat, error) {
	const q = `
	SELECT
		c.id,
		c.name,
		c.race,
		c.sex,
		c.age_in_month,
		c.description
	FROM cats c`

	// execute query
	rows, err := db.pool.Query(ctx, q)

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

func (db *Database) Save(ctx context.Context, model cat.Cat) (*cat.Cat, error) {
	var result *cat.Cat

	if err := database.BeginTransaction(ctx, db.pool, func(tx pgx.Tx) error {
		const q = `INSERT into cats
			("user_id", "name", "race", "sex", "age_in_month", "description")
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, name, race, sex, age_in_month, description;`

		// execute query to insert new record. it takes 'cat' variable as its input
		// the result will be placed in 'row' variable
		row := db.pool.QueryRow(
			ctx, q,
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
			return fmt.Errorf("Save %w", err)
		}

		result = c
		return nil
	}); err != nil {
		return nil, fmt.Errorf("Save transaction %w", err)
	}

	// return 'c' and nil if no error found
	return result, nil
}
