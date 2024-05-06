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

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Create(req *request.UserRegisterRequest) (*createReturn, error)
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

	model := user.User{
		Email:    payload.Email,
		Name:     payload.Name,
		Password: payload.Password,
	}

	// Create user
	userCreated, err := repo.User.Save(svc.context, model)
	if err != nil {
		return nil, err
	}
	userCredential = userCreated

	_, err = repo.Session.SaveOrGet(svc.context, userCredential)
	if err != nil {
		return nil, err
	}

	// Generate access token
	token, err := jwt.GenerateAccessToken(uint64(userCredential.Id), userCredential)
	if err != nil {
		return nil, err
	}

	accessToken = token

	return &createReturn{
		User:        userCredential,
		AccessToken: accessToken,
	}, nil
}
