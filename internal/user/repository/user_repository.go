package repository

import (
	"context"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/pkg/database"
	"enigmanations/cats-social/pkg/jwt"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Get(ctx context.Context, id int) (*user.User, error)
	GetByEmailIfExists(ctx context.Context, email string) (*user.User, error)
	Save(ctx context.Context, model user.User) (*user.User, *user.UserSession, error)
}

type userRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepositoryDB{pool: pool}
}

func (db *userRepositoryDB) Save(ctx context.Context, model user.User) (*user.User, *user.UserSession, error) {
	sessionLengthSeconds := jwt.AccessTokenDurationSeconds

	var (
		result  *user.User
		session = &user.UserSession{
			ExpiresAt: time.Now().Add(sessionLengthSeconds),
		}
	)

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

		// Create user session
		const qSession = `
			INSERT INTO sessions (token, expires_at, user_id, created_at)
			VALUES($1, $2, $3, now())
			RETURNING token;
		`

		// Generate access token
		token, err := jwt.GenerateAccessToken(uint64(u.Id), &model)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		session.Token = token

		userSessionRow := db.pool.QueryRow(
			ctx,
			qSession,
			token,
			session.ExpiresAt,
			u.Id,
		)
		uSession := new(user.UserSession)
		uSessionErr := userSessionRow.Scan(
			&uSession.Token,
		)
		if uSessionErr != nil {
			return fmt.Errorf("%w", uSessionErr)
		}

		result = u
		session = uSession
		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("Save transaction %w", err)
	}

	return result, session, nil
}

func (db *userRepositoryDB) Get(ctx context.Context, id int) (*user.User, error) {
	const q = `SELECT * FROM users WHERE id = $1 limit 1;`

	row := db.pool.QueryRow(ctx, q, id)
	u := new(user.User)

	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.Email,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

type Exists struct {
	exists bool
}

func (db *userRepositoryDB) GetByEmailIfExists(ctx context.Context, email string) (*user.User, error) {
	const sql = `
		SELECT EXISTS (
			SELECT u.id, u.name, u.email FROM users u WHERE u.email = $1 limit 1
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
			SELECT u.id, u.name, u.email FROM users u WHERE u.email = $1 limit 1
		`
		row := db.pool.QueryRow(ctx, sql, email)
		u := new(user.User)
		err := row.Scan(
			&u.Id,
			&u.Name,
			&u.Email,
		)
		if err != nil {
			return nil, err
		}

		return u, nil
	}

	return nil, nil
}
