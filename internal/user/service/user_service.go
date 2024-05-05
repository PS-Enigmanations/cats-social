package service

import (
	"context"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/internal/user/errs"
	"enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/internal/user/request"
	"enigmanations/cats-social/pkg/bcrypt"
	"enigmanations/cats-social/pkg/jwt"
	"enigmanations/cats-social/util"
	"fmt"

	"enigmanations/cats-social/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Create(req *request.UserRegisterRequest) (*createReturn, error)
}

type UserDependency struct {
	User    repository.UserRepository
	Session repository.UserAuthRepository
}

type userService struct {
	pool    *pgxpool.Pool
	repo    *UserDependency
	Context context.Context
}

func NewUserService(ctx context.Context, pool *pgxpool.Pool, repo *UserDependency) UserService {
	return &userService{Context: ctx, pool: pool, repo: repo}
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
		userFound, _ := repo.User.GetByEmailIfExists(svc.Context, req.Email)
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
	User *user.User
	AccessToken    string
}

func (svc *userService) Create(req *request.UserRegisterRequest) (*createReturn, error) {
	repo := svc.repo

	// Validate first
	payload, err := svc.validate(req)
	if err != nil {
		return nil, err
	}

	var (
		userCredential *user.User
		accessToken    string
	)

	if err := database.BeginTransaction(svc.Context, svc.pool, func(tx pgx.Tx) error {
		model := user.User{
			Email:    payload.Email,
			Name:     payload.Name,
			Password: payload.Password,
		}

		// Create user
		userCreated, err := repo.User.Save(svc.Context, model, tx)
		if err != nil {
			return err
		}
		userCredential = userCreated

		_, err = repo.Session.SaveOrGet(svc.Context, userCredential, tx)
		if err != nil {
			return err
		}

		// Generate access token
		token, err := jwt.GenerateAccessToken(uint64(userCredential.Id), userCredential)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		accessToken = token

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction %w", err)
	}

	return &createReturn{
		User: userCredential,
		AccessToken: accessToken,
	}, nil
}
