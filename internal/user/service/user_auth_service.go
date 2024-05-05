package service

import (
	"context"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/internal/user/errs"
	"enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/internal/user/request"
	"enigmanations/cats-social/pkg/bcrypt"
	"enigmanations/cats-social/pkg/jwt"
	"fmt"

	"enigmanations/cats-social/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAuthService interface {
	Login(req *request.UserLoginRequest) (*loginReturn, error)
}

type UserAuthDependency struct {
	User    repository.UserRepository
	Session repository.UserAuthRepository
}

type userAuthService struct {
	pool    *pgxpool.Pool
	repo    *UserAuthDependency
	Context context.Context
}

func NewUserAuthService(ctx context.Context, pool *pgxpool.Pool, repo *UserAuthDependency) UserAuthService {
	return &userAuthService{Context: ctx, pool: pool, repo: repo}
}

type loginReturn struct {
	User *user.User
	UserSession *user.UserSession
	AccessToken    string
}

func (svc *userAuthService) Login(req *request.UserLoginRequest) (*loginReturn, error) {
	repo := svc.repo

	var (
		userCredential *user.User
		userSession    *user.UserSession
		accessToken    string
	)

	// Check email
	if req.Email != "" {
		// Get user
		userCredentialFound, err := repo.User.GetByEmail(svc.Context, req.Email)
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
	if err := database.BeginTransaction(svc.Context, svc.pool, func(tx pgx.Tx) error {
		session, err := repo.Session.SaveOrGet(svc.Context, userCredential, tx)
		if err != nil {
			return err
		}
		userSession = session

		// Generate access token
		token, err := jwt.GenerateAccessToken(uint64(userSession.UserId), userCredential)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		accessToken = token
		return nil
	}); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &loginReturn{
		User: userCredential,
		UserSession: userSession,
		AccessToken: accessToken,
	}, nil
}
