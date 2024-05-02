package service

import (
	"context"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/internal/user/errs"
	"enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/internal/user/request"
	"enigmanations/cats-social/internal/user/response"
	"enigmanations/cats-social/pkg/bcrypt"
	"fmt"
	"strings"

	"enigmanations/cats-social/pkg/database"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Create(req *request.UserRegisterRequest) (*response.UserCreateResponse, error)
}

type UserDependency struct {
	User repository.UserRepository
}

type userService struct {
	pool    *pgxpool.Pool
	repo    *UserDependency
	Context context.Context
}

func NewUserService(ctx context.Context, pool *pgxpool.Pool, repo *UserDependency) UserService {
	return &userService{Context: ctx, pool: pool, repo: repo}
}

func (svc *userService) Create(req *request.UserRegisterRequest) (*response.UserCreateResponse, error) {
	repo := svc.repo

	var result *user.User

	if err := database.BeginTransaction(svc.Context, svc.pool, func(tx pgx.Tx) error {
		var payload = &user.User{
			Name: req.Name,
		}

		if req.Email != "" {
			lowerCasedEmail := strings.ToLower(req.Email)
			payload.Email = lowerCasedEmail

			// Check email format
			var verifier = emailverifier.NewVerifier()
			ret, err := verifier.Verify(payload.Email)
			if err != nil {
				return errs.UserErrEmailInvalidFormat
			}
			if !ret.Syntax.Valid {
				return errs.UserErrEmailInvalidFormat
			}

			// Check exisiting user
			userFound, _ := repo.User.GetByEmailIfExists(svc.Context, req.Email)
			if userFound != nil {
				return errs.UserErrEmailExist
			}
		}
		if req.Password != "" {
			hashedPassword, err := bcrypt.HashPassword(req.Password)
			if err != nil {
				return err
			}
			payload.Password = hashedPassword
		}

		model := user.User{
			Email:    payload.Email,
			Name:     payload.Name,
			Password: payload.Password,
		}

		user, err := repo.User.Save(svc.Context, model, tx)
		if err != nil {
			return err
		}

		result = user
		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction %w", err)
	}

	return response.UserToUserCreateResponse(*result), nil
}
