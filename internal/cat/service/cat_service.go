package service

import (
	"context"

	"enigmanations/cats-social/internal/cat/repository"
	"enigmanations/cats-social/internal/cat/request"
	"enigmanations/cats-social/internal/cat/response"
)

type CatService interface {
	GetAllByParams(p *request.CatGetAllQueryParams) (*response.CatGetAllResponse, error)
	// Create(payload *request.CatCreateRequest) (*response.CatCreateResponse, error)
	//Update(payload *request.CatUpdateRequest, catId int) error
	// Delete(catId int) error
}

type catService struct {
	db      repository.CatRepository
	Context context.Context
}

// NewService creates an API service.
func NewCatService(db repository.CatRepository, ctx context.Context) CatService {
	return &catService{db: db, Context: ctx}
}

func (service *catService) GetAllByParams(p *request.CatGetAllQueryParams) (*response.CatGetAllResponse, error) {
	cats, err := service.db.GetAllByParams(service.Context)

	if err != nil {
		return nil, err
	}

	catShows := response.ToCatShows(cats)
	return response.CatToCatGetAllResponse(catShows), nil
}

/**
func (service *catService) Create(payload *request.CatCreateRequest) (*response.CatCreateResponse, error) {
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

	err = service.db.SaveImageUrls(service.Context, cat.Id, payload.ImageUrls)
	if err != nil {
		return nil, err
	}

	cat.ImageUrls = payload.ImageUrls

	return CatToCatResponse(*cat), nil
}


func (service *catService) Update(payload *request.CatUpdateRequest, catId int) error {
	const USER_ID = 2

	model := cat.Cat{
		UserId:      USER_ID,
		Id:          catId,
		Name:        payload.Name,
		Race:        cat.Race(payload.Race),
		Sex:         cat.Sex(payload.Sex),
		AgeInMonth:  payload.AgeInMonth,
		Description: payload.Description,
	}

	cat, err := service.db.Update(service.Context, model)

	if err != nil {
		return nil, err
	}

	err = service.db.DeleteImageUrls(service.Context, cat.Id)
	if err != nil {
		return nil, err
	}

	err = service.db.SaveImageUrls(service.Context, cat.Id, payload.ImageUrls)
	if err != nil {
		return nil, err
	}

	cat.ImageUrls = payload.ImageUrls

	return CatToCatResponse(*cat), nil
}

func (service *catService) Delete(catId int) error {
	_, err := service.db.FindById(service.Context, catId)
	if err != nil {
		return err
	}

	err = service.db.Delete(service.Context, catId)
	if err != nil {
		return err
	}
	return nil
}
*/
