package repository

import (
	"context"
	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/pkg/database/postgres/interfaces"
)

type ICatRepository interface {
	GetAll(ctx context.Context) (*cat.Cat, error)
}

type Database struct {
	pool interfaces.PGXQuerier
}

func NewCatRepository(pool interfaces.PGXQuerier) Database {
	return Database{
		pool: pool,
	}
}

func (db *Database) GetAll(ctx context.Context) ([]*cat.Cat, error) {
	const q = `SELECT * FROM cats`

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
