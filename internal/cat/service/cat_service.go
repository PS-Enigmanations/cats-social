package service

import (
	"context"

	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/internal/cat/repository"
)

type CatService interface {
	GetAll() ([]*CatResponse, error)
	FindById(catId int) (*CatResponse, error)
	Create(payload *cat.CatCreateRequest) (*CatResponse, error)
	Update(payload *cat.CatUpdateRequest, catId int) (*CatResponse, error)
	Delete(catId int) error
	DeleteImageUrls(catId int) error
}

type catService struct {
	db      repository.CatRepository
	Context context.Context
}

// NewService creates an API service.
func NewCatService(db repository.CatRepository, ctx context.Context) CatService {
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

	err = service.db.SaveImageUrls(service.Context, cat.Id, payload.ImageUrls)
	if err != nil {
		return nil, err
	}

	cat.ImageUrls = payload.ImageUrls

	return CatToCatResponse(*cat), nil
}

func (service *catService) Update(payload *cat.CatUpdateRequest, catId int) (*CatResponse, error) {
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

	cat, err := service.db.FindById(service.Context, catId)

	if cat == nil {
		return nil, err
	}

	cat, err = service.db.Update(service.Context, model)

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

func (service *catService) FindById(catId int) (*CatResponse, error) {
	cat, err := service.db.FindById(service.Context, catId)
	if err != nil {
		return nil, err
	}

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

func (service *catService) DeleteImageUrls(catId int) error {
	return service.db.DeleteImageUrls(service.Context, catId)
}

type CatResponse struct {
	Id          int
	Name        string
	Race        string
	Sex         string
	AgeInMonth  int
	Description string
	ImageUrls   []string
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
		ImageUrls:   c.ImageUrls,
	}
}
