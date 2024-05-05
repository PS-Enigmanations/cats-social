package service

import (
	"context"
	"enigmanations/cats-social/internal/session"
	"enigmanations/cats-social/internal/session/repository"
	"enigmanations/cats-social/internal/session/request"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/internal/user/errs"
	userRepository "enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/pkg/bcrypt"
	"enigmanations/cats-social/pkg/jwt"
	"fmt"

	"enigmanations/cats-social/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionService interface {
	Login(req *request.SessionLoginRequest) (*loginReturn, error)
	GenerateAccessToken(userId int, user *user.User) (string, error)
}

type SessionDependency struct {
	Session repository.SessionRepository
	User    userRepository.UserRepository
}

type sessionService struct {
	pool    *pgxpool.Pool
	repo    *SessionDependency
	context context.Context
}

func NewSessionService(ctx context.Context, pool *pgxpool.Pool, repo *SessionDependency) SessionService {
	return &sessionService{context: ctx, pool: pool, repo: repo}
}

type loginReturn struct {
	User        *user.User
	UserSession *session.Session
	AccessToken string
}

func (svc *sessionService) Login(req *request.SessionLoginRequest) (*loginReturn, error) {
	repo := svc.repo

	var (
		userCredential *user.User
		userSession    *session.Session
		accessToken    string
	)

	// Check email
	if req.Email != "" {
		// Get user
		userCredentialFound, err := repo.User.GetByEmail(svc.context, req.Email)
		if err != nil {
			return nil, errs.UserErrNotFound
		}

		userCredential = userCredentialFound
	}

	// Check password
	if req.Password != "" {
		if !bcrypt.CheckPasswordHash(req.Password, userCredential.Password) {
			return nil, errs.WrongPassword
		}
	}

	// Create or get session
	if err := database.BeginTransaction(svc.context, svc.pool, func(tx pgx.Tx, ctx context.Context) error {
		session, err := repo.Session.SaveOrGet(ctx, userCredential, tx)
		if err != nil {
			return err
		}
		userSession = session

		// Generate access token
		token, err := svc.GenerateAccessToken(session.UserId, userCredential)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		accessToken = token
		return nil
	}); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &loginReturn{
		User:        userCredential,
		UserSession: userSession,
		AccessToken: accessToken,
	}, nil
}

func (svc *sessionService) GenerateAccessToken(userId int, user *user.User) (string, error) {
	token, err := jwt.GenerateAccessToken(uint64(userId), user)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return token, nil

}
