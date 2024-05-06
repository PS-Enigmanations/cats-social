package repository

import (
	"context"
	"enigmanations/cats-social/internal/session"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/pkg/jwt"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	GetIfExists(ctx context.Context, userId int) (*session.Session, error)
	Save(ctx context.Context, model *user.User) (*session.Session, error)
	SaveOrGet(ctx context.Context, model *user.User) (*session.Session, error)
}

type sessionRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewUserSessionRepository(pool *pgxpool.Pool) SessionRepository {
	return &sessionRepositoryDB{pool: pool}
}

func (db *sessionRepositoryDB) Save(ctx context.Context, model *user.User) (*session.Session, error) {
	const sessionLengthSeconds = 134784000 // 1 year

	var sessionValue = &session.Session{
		ExpiresAt: time.Now().Add(time.Duration(sessionLengthSeconds) * time.Second),
	}

	const sql = `
		INSERT INTO sessions (token, expires_at, user_id, created_at)
		VALUES($1, $2, $3, now())
		RETURNING token, user_id;
	`

	// Generate access token
	token, err := jwt.GenerateAccessToken(uint64(model.Id), model)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	userSessionRow := db.pool.QueryRow(
		ctx,
		sql,
		token,
		sessionValue.ExpiresAt,
		model.Id,
	)
	v := new(session.Session)
	uSessionErr := userSessionRow.Scan(
		&v.Token,
		&v.UserId,
	)
	if uSessionErr != nil {
		return nil, fmt.Errorf("%w", uSessionErr)
	}
	sessionValue = v

	return sessionValue, nil
}

type getIfExists struct {
	exists bool
}

func (db *sessionRepositoryDB) GetIfExists(ctx context.Context, userId int) (*session.Session, error) {
	const sql = `
		SELECT EXISTS (
			SELECT s."token" from sessions s WHERE s.user_id = $1 AND deleted_at IS NULL LIMIT 1
		);`

	row := db.pool.QueryRow(ctx, sql, userId)
	s := new(getIfExists)
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
		u := new(session.Session)
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

func (db *sessionRepositoryDB) SaveOrGet(ctx context.Context, model *user.User) (*session.Session, error) {
	var userSession *session.Session

	userSessionFound, _ := db.GetIfExists(ctx, model.Id)
	if userSessionFound != nil {
		userSession = userSessionFound
	} else {
		userSessionCreated, err := db.Save(ctx, model)
		if err != nil {
			return nil, err
		}
		userSession = userSessionCreated
	}

	return userSession, nil
}
