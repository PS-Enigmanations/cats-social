package repository

import (
	"context"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/pkg/database"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Get(ctx context.Context, id int) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByEmailIfExists(ctx context.Context, email string) (*user.User, error)
	Save(ctx context.Context, model user.User) (*user.User, error)
}

type userRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepositoryDB{pool: pool}
}

func (db *userRepositoryDB) Save(ctx context.Context, model user.User) (*user.User, error) {
	var result *user.User

	if err := database.BeginTransaction(ctx, db.pool, func(tx pgx.Tx) error {
		// Create user
		const qUser = `
			INSERT INTO users ("name", email, "password", created_at)
			VALUES($1, $2, $3, now())
			RETURNING id, email, name;
		`
		userRow := db.pool.QueryRow(
			ctx,
			qUser,
			model.Name,
			model.Email,
			model.Password,
		)
		u := new(user.User)
		uErr := userRow.Scan(
			&u.Id,
			&u.Email,
			&u.Name,
		)
		if uErr != nil {
			return fmt.Errorf("%w", uErr)
		}

		result = u
		return nil
	}); err != nil {
		return nil, fmt.Errorf("Save transaction %w", err)
	}

	return result, nil
}

func (db *userRepositoryDB) Get(ctx context.Context, id int) (*user.User, error) {
	const q = `SELECT u.id, u.name, u.email FROM users u WHERE u.id = $1 AND deleted_at IS NULL LIMIT 1;`

	row := db.pool.QueryRow(ctx, q, id)
	u := new(user.User)

	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.Email,
	)
	if err != nil {
		log.Print("Error getting user", err)
		return nil, err
	}

	return u, nil
}

type Exists struct {
	exists bool
}

func (db *userRepositoryDB) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	const sql = `
		SELECT u.id, u.name, u.email, u.password FROM users u WHERE u.email = $1 AND deleted_at IS NULL LIMIT 1;
	`
	row := db.pool.QueryRow(ctx, sql, email)
	u := new(user.User)
	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&u.Password,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *userRepositoryDB) GetByEmailIfExists(ctx context.Context, email string) (*user.User, error) {
	const sql = `
		SELECT EXISTS (
			SELECT u.id, u.name, u.email, u.password FROM users u WHERE u.email = $1 AND deleted_at IS NULL LIMIT 1
		);`

	row := db.pool.QueryRow(ctx, sql, email)
	e := new(Exists)
	err := row.Scan(
		&e.exists,
	)
	if err != nil {
		return nil, err
	}

	if e.exists {
		const sql = `
			SELECT u.id, u.name, u.email FROM users u WHERE u.email = $1 AND deleted_at IS NULL LIMIT 1;
		`
		row := db.pool.QueryRow(ctx, sql, email)
		u := new(user.User)
		err := row.Scan(
			&u.Id,
			&u.Name,
			&u.Email,
			&u.Password,
		)
		if err != nil {
			return nil, err
		}

		return u, nil
	}

	return nil, nil
}
