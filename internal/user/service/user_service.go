package service

import (
	"context"
	sessionRepository "enigmanations/cats-social/internal/session/repository"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/internal/user/errs"
	"enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/internal/user/request"
	"enigmanations/cats-social/pkg/bcrypt"
	"enigmanations/cats-social/pkg/jwt"
	"enigmanations/cats-social/util"

	"enigmanations/cats-social/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Create(req *request.UserRegisterRequest) <-chan util.Result[*createReturn]
}

type UserDependency struct {
	User    repository.UserRepository
	Session sessionRepository.SessionRepository
}

type userService struct {
	pool    *pgxpool.Pool
	repo    *UserDependency
	context context.Context
}

func NewUserService(ctx context.Context, pool *pgxpool.Pool, repo *UserDependency) UserService {
	return &userService{context: ctx, pool: pool, repo: repo}
}

func (svc *userService) validate(req *request.UserRegisterRequest) (*user.User, error) {
	repo := svc.repo

	var payload = &user.User{
		Name: req.Name,
	}

	if req.Email != "" {
		// Check email format
		if !util.IsEmail(req.Email) {
			return nil, errs.UserErrEmailInvalidFormat
		}

		// Check existing user
		userFound, _ := repo.User.GetByEmailIfExists(svc.context, req.Email)
		if userFound != nil {
			return nil, errs.UserErrEmailExist
		}
		payload.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		payload.Password = hashedPassword
	}

	return payload, nil
}

type createReturn struct {
	User        *user.User
	AccessToken string
}

func (svc *userService) Create(req *request.UserRegisterRequest) <-chan util.Result[*createReturn] {
	repo := svc.repo

	result := make(chan util.Result[*createReturn])
	go func() {
		// Validate first
		payload, err := svc.validate(req)
		if err != nil {
			result <- util.Result[*createReturn]{
				Error: err,
			}

			return
		}

		if err := database.BeginTransaction(svc.context, svc.pool, func(tx pgx.Tx, ctx context.Context) error {
			model := user.User{
				Email:    payload.Email,
				Name:     payload.Name,
				Password: payload.Password,
			}

			// Create user
			userCreated, err := repo.User.Save(ctx, model, tx)
			if err != nil {
				result <- util.Result[*createReturn]{
					Error: err,
				}

				return err
			}

			_, err = repo.Session.SaveOrGet(ctx, userCreated, tx)
			if err != nil {
				result <- util.Result[*createReturn]{
					Error: err,
				}

				return err
			}

			// Generate access token
			token, err := jwt.GenerateAccessToken(uint64(userCreated.Id), userCreated)
			if err != nil {
				result <- util.Result[*createReturn]{
					Error: err,
				}

				return err
			}

			result <- util.Result[*createReturn]{
				Result: &createReturn{
					User:        userCreated,
					AccessToken: token,
				},
			}
			close(result)

			return nil
		}); err != nil {
			result <- util.Result[*createReturn]{
				Error: err,
			}
		}
	}()

	return result
}
