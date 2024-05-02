package service

import (
	"context"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/internal/user/errs"
	"enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/internal/user/request"
	"enigmanations/cats-social/internal/user/response"
	"enigmanations/cats-social/pkg/bcrypt"
	"enigmanations/cats-social/pkg/jwt"
	"fmt"

	"enigmanations/cats-social/pkg/database"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAuthService interface {
	Login(req *request.UserLoginRequest) (*response.UserLoginResponse, error)
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

func (svc *userAuthService) Login(req *request.UserLoginRequest) (*response.UserLoginResponse, error) {
	repo := svc.repo

	var (
		userCredential *user.User
		userSession    *user.UserSession
		accessToken    string
	)

	// Check email
	if req.Email != "" {
		var verifier = emailverifier.NewVerifier()
		ret, err := verifier.Verify(req.Email)
		if err != nil {
			return nil, errs.UserErrEmailInvalidFormat
		}
		if !ret.Syntax.Valid {
			return nil, errs.UserErrEmailInvalidFormat
		}

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
		userSessionFound, _ := repo.Session.GetIfExists(svc.Context, userCredential.Id)
		if userSessionFound != nil {
			userSession = userSessionFound
		} else {
			userSessionCreated, err := repo.Session.Save(svc.Context, userCredential, tx)
			if err != nil {
				return err
			}
			userSession = userSessionCreated
		}

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

	return response.UserToUserLoginResponse(*userCredential, accessToken), nil
}
