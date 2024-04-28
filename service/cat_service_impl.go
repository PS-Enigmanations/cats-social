package service

import (
	"PS-Enigmanations/cats-social/helper"
	"PS-Enigmanations/cats-social/model/domain"
	"PS-Enigmanations/cats-social/model/web"
	"PS-Enigmanations/cats-social/repository"
	"context"
	"database/sql"

	"github.com/go-playground/validator"
)

type CatServiceImpl struct {
	CatRepository repository.CatRepository
	DB            *sql.DB
	Validate      *validator.Validate
}

func NewCatService(catRepository repository.CatRepository, db *sql.DB, validate *validator.Validate) CatService {
	return &CatServiceImpl{
		CatRepository: catRepository,
		DB:            db,
		Validate:      validate,
	}
}

func (service *CatServiceImpl) Create(ctx context.Context, request web.CatCreateRequest) web.CatResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cat := domain.Cat{
		Name: request.Name,
	}

	cat = service.CatRepository.Save(ctx, tx, cat)

	return helper.ToCatResponse(cat)
}

func (service *CatServiceImpl) Update() {
	panic("not implemented") // TODO: Implement
}

func (service *CatServiceImpl) Delete() {
	panic("not implemented") // TODO: Implement
}

func (service *CatServiceImpl) FindById() {
	panic("not implemented") // TODO: Implement
}

func (service *CatServiceImpl) Get() {
	panic("not implemented") // TODO: Implement
}
