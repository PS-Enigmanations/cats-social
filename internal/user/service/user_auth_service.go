package service

import (
	"context"
	"enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/internal/user/request"
	"enigmanations/cats-social/internal/user/response"
)

type UserAuthService interface {
	Login(req *request.UserLoginRequest) (*response.UserLoginResponse, error)
}

type userAuthService struct {
	db      repository.UserAuthRepository
	Context context.Context
}

func NewUserAuthService(db repository.UserAuthRepository, ctx context.Context) UserAuthService {
	return &userAuthService{db: db, Context: ctx}
}

func (service *userAuthService) Login(req *request.UserLoginRequest) (*response.UserLoginResponse, error) {
	return nil, nil
}
