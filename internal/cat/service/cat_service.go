package service

import (
	"context"

	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/internal/cat/repository"
)

type CatService interface {
	GetAll() ([]*CatResponse, error)
	Create(payload *cat.CatCreateRequest) (*CatResponse, error)
}

type catService struct {
	db      repository.Database
	Context context.Context
}

// NewService creates an API service.
func NewCatService(db repository.Database, ctx context.Context) *catService {
	return &catService{db: db, Context: ctx}
}

func (service *catService) GetAll() ([]*CatResponse, error) {
	// call GetAll from repository/ datastore to retrieve all Cat record
	cats, err := service.db.GetAll(service.Context)

	if err != nil {
		return nil, err
	}

	var catRes []*CatResponse
	for _, cat := range cats {
		catRes = append(catRes, CatToCatResponse(*cat))
	}

	return catRes, nil
}

func (service *catService) Create(payload *cat.CatCreateRequest) (*CatResponse, error) {
	const USER_ID = 2

	model := cat.Cat{
		UserId:      USER_ID,
		Name:        payload.Name,
		Race:        cat.Race(payload.Race),
		Sex:         cat.Sex(payload.Sex),
		AgeInMonth:  payload.AgeInMonth,
		Description: payload.Description,
	}

	// call Create from repository/ datastore
	cat, err := service.db.Save(service.Context, model)

	// if error occur, return nil for the response as well as return the error
	if err != nil {
		return nil, err
	}

	return CatToCatResponse(*cat), nil
}

type CatResponse struct {
	Id          int
	Name        string
	Race        string
	Sex         string
	AgeInMonth  int
	Description string
}

// convert 'Cat' model to 'CatResponse' DTO
func CatToCatResponse(c cat.Cat) *CatResponse {
	return &CatResponse{
		Id:          c.Id,
		Name:        c.Name,
		Race:        string(c.Race),
		Sex:         string(c.Sex),
		AgeInMonth:  c.AgeInMonth,
		Description: c.Description,
	}
}
