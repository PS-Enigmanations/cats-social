package service

import (
	"context"

	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/internal/cat/repository"

	"github.com/go-playground/validator"
)

type CatService interface {
	GetAll() ([]*CatResponse, error)
}

type catService struct {
	db       repository.Database
	Context  context.Context
	Validate *validator.Validate
}

// NewService creates an API service.
func NewCatService(db repository.Database, ctx context.Context, validate *validator.Validate) *catService {
	return &catService{db: db, Context: ctx, Validate: validate}
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

type CatResponse struct {
	Id   int
	Name string
}

// convert 'Cat' model to 'CatResponse' DTO
func CatToCatResponse(c cat.Cat) *CatResponse {
	return &CatResponse{
		Id:   c.Id,
		Name: c.Name,
	}
}
