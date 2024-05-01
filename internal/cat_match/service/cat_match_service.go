package service

import (
	"context"
	"enigmanations/cats-social/internal/cat_match/repository"
	"enigmanations/cats-social/internal/cat_match/request"
	"enigmanations/cats-social/internal/cat_match/response"
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

/**
func (svc *catMatchService) validate(l cat, r cat.Cat) error {
	// Check cat and owner is exist
	catFound, _ := svc.catRepo.GetIfExists(svc.Context, l.Id)

	// Check owner
	if l.UserId == r.UserId {
		return errs.CatMatchErrOwner
	}

	// Check gender
	equalMatches := reflect.DeepEqual(l.Sex, r.Sex)
	if equalMatches {
		return errs.CatMatchErrGender
	}

	return nil
}
*/

func (svc *catMatchService) Create(req *request.CatMatchRequest) (*response.CatMatchCreateResponse, error) {
	return nil, nil
}
