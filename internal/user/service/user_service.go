package service

import (
	"context"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/internal/user/request"
	"enigmanations/cats-social/internal/user/response"
	"enigmanations/cats-social/pkg/bcrypt"
	"errors"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
)

type UserService interface {
	Create(req *request.UserCreateRequest) (*response.UserCreateResponse, error)
}

type userService struct {
	db      repository.UserRepository
	Context context.Context
}

func NewUserService(db repository.UserRepository, ctx context.Context) UserService {
	return &userService{db: db, Context: ctx}
}

func (service *userService) Create(req *request.UserCreateRequest) (*response.UserCreateResponse, error) {
	var payload = &user.User{
		Name: req.Name,
	}

	var verifier = emailverifier.NewVerifier()

	if req.Email != "" {
		lowerCasedEmail := strings.ToLower(req.Email)
		payload.Email = lowerCasedEmail

		ret, err := verifier.Verify(payload.Email)
		if err != nil {
			return nil, err
		}

		if !ret.Syntax.Valid {
			return nil, errors.New("Email syntax is invalid!")
		}
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		payload.Password = hashedPassword
	}

	model := user.User{
		Email:    payload.Email,
		Name:     payload.Name,
		Password: payload.Password,
	}

	user, userSession, err := service.db.Save(service.Context, model)
	if err != nil {
		return nil, err
	}

	return response.UserToUserCreateResponse(*user, *userSession), nil
}
