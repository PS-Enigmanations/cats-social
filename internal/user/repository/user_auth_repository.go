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
}

type userAuthRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewUserAuthRepository(pool *pgxpool.Pool) UserAuthRepository {
	return &userAuthRepositoryDB{pool: pool}
}

// Create user session
func (db *userAuthRepositoryDB) Save(ctx context.Context, model *user.User) (*user.UserSession, error) {
	sessionLengthSeconds := jwt.AccessTokenDurationSeconds

	var session = &user.UserSession{
		ExpiresAt: time.Now().Add(sessionLengthSeconds),
	}

	if err := database.BeginTransaction(ctx, db.pool, func(tx pgx.Tx) error {
		const sql = `
			INSERT INTO sessions (token, expires_at, user_id, created_at)
			VALUES($1, $2, $3, now())
			RETURNING token;
		`

		// Generate access token
		token, err := jwt.GenerateAccessToken(uint64(model.Id), model)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		userSessionRow := db.pool.QueryRow(
			ctx,
			sql,
			token,
			session.ExpiresAt,
			model.Id,
		)
		uSession := new(user.UserSession)
		uSessionErr := userSessionRow.Scan(
			&uSession.Token,
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
