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
	Save(ctx context.Context, req user.User) (*user.User, *user.UserSession, error)
}

type Database struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &Database{pool: pool}
}

func (db *Database) Save(ctx context.Context, req user.User) (*user.User, *user.UserSession, error) {
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
			req.Name,
			req.Email,
			req.Password,
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
		token, err := jwt.GenerateAccessToken(uint64(u.Id), &req)
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
