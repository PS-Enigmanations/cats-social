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

type UserAuthRepository interface {
	Save(ctx context.Context, model *user.User) (*user.UserSession, error)
	GetIfExists(ctx context.Context, userId int) (*user.UserSession, error)
}

type userAuthRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewUserAuthRepository(pool *pgxpool.Pool) UserAuthRepository {
	return &userAuthRepositoryDB{pool: pool}
}

// Create user session
func (db *userAuthRepositoryDB) Save(ctx context.Context, model *user.User) (*user.UserSession, error) {
	const sessionLengthSeconds = 134784000 // 1 year

	var session = &user.UserSession{
		ExpiresAt: time.Now().Add(time.Duration(sessionLengthSeconds) * time.Second),
	}

	if err := database.BeginTransaction(ctx, db.pool, func(tx pgx.Tx) error {
		const sql = `
			INSERT INTO sessions (token, expires_at, user_id, created_at)
			VALUES($1, $2, $3, now())
			RETURNING token, user_id;
		`

		// Generate access token
		token, err := jwt.GenerateAccessToken(uint64(model.Id), model)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		userSessionRow := tx.QueryRow(
			ctx,
			sql,
			token,
			session.ExpiresAt,
			model.Id,
		)
		uSession := new(user.UserSession)
		uSessionErr := userSessionRow.Scan(
			&uSession.Token,
			&uSession.UserId,
		)
		if uSessionErr != nil {
			return fmt.Errorf("%w", uSessionErr)
		}

		session = uSession
		return nil
	}); err != nil {
		return nil, fmt.Errorf("Save transaction %w", err)
	}

	return session, nil
}

func (db *userAuthRepositoryDB) GetIfExists(ctx context.Context, userId int) (*user.UserSession, error) {
	const sql = `
		SELECT EXISTS (
			SELECT s."token" from sessions s WHERE s.user_id = $1 AND deleted_at IS NULL LIMIT 1
		);`

	row := db.pool.QueryRow(ctx, sql, userId)
	s := new(Exists)
	err := row.Scan(
		&s.exists,
	)
	if err != nil {
		return nil, err
	}

	if s.exists {
		const sql = `
			SELECT s."token", s.user_id from sessions s WHERE s.user_id = $1 AND deleted_at IS NULL LIMIT 1
		`
		row := db.pool.QueryRow(ctx, sql, userId)
		u := new(user.UserSession)
		err := row.Scan(
			&u.Token,
			&u.UserId,
		)
		if err != nil {
			return nil, err
		}

		return u, nil
	}

	return nil, nil
}
