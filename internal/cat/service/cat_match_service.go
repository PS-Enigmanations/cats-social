package service

import (
	"context"
	"enigmanations/cats-social/internal/cat/repository"
	"enigmanations/cats-social/internal/cat/request"
	"enigmanations/cats-social/internal/cat/response"
)

type CatMatchService interface {
	Create(req *request.CatMatchRequest) (*response.CatMatchCreateResponse, error)
}

type catMatchService struct {
	repo    repository.CatMatchRepository
	Context context.Context
}

func NewCatMatchService(repo repository.CatMatchRepository, ctx context.Context) CatMatchService {
	return &catMatchService{repo: repo, Context: ctx}
}

func (svc *catMatchService) Create(req *request.CatMatchRequest) (*response.CatMatchCreateResponse, error) {
	return nil, nil
}
