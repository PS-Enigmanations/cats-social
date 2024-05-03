package service

import (
	"context"
	"fmt"

	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/internal/cat/repository"
	"enigmanations/cats-social/internal/cat/request"
	"enigmanations/cats-social/internal/cat/response"

	"enigmanations/cats-social/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatService interface {
	GetAllByParams(p *request.CatGetAllQueryParams, actorId int) (*response.CatGetAllResponse, error)
	Create(payload *request.CatCreateRequest, actorId int) (*response.CatCreateResponse, error)
}

type catService struct {
	db      repository.CatRepository
	pool    *pgxpool.Pool
	Context context.Context
}

// NewService creates an API service.
func NewCatService(ctx context.Context, pool *pgxpool.Pool, db repository.CatRepository) CatService {
	return &catService{db: db, pool: pool, Context: ctx}
}

func (svc *catService) GetAllByParams(p *request.CatGetAllQueryParams, actorId int) (*response.CatGetAllResponse, error) {
	cats, err := svc.db.GetAllByParams(svc.Context, p, actorId)

	if err != nil {
		return nil, err
	}

	catShows := response.ToCatShows(cats)
	return response.CatToCatGetAllResponse(catShows), nil
}

func (svc *catService) Create(payload *request.CatCreateRequest, actorId int) (*response.CatCreateResponse, error) {
	var result *response.CatCreateResponse

	if err := database.BeginTransaction(svc.Context, svc.pool, func(tx pgx.Tx) error {
		model := cat.Cat{
			UserId:      actorId,
			Name:        payload.Name,
			Race:        cat.Race(payload.Race),
			Sex:         cat.Sex(payload.Sex),
			AgeInMonth:  payload.AgeInMonth,
			Description: payload.Description,
		}

		// call Create from repository/ datastore
		cat, err := svc.db.Save(svc.Context, tx, model)

		// if error occur, return nil for the response as well as return the error
		if err != nil {
			return nil
		}

		err = svc.db.SaveImageUrls(svc.Context, tx, cat.Id, payload.ImageUrls)
		if err != nil {
			return nil
		}

		cat.ImageUrls = payload.ImageUrls
		result = response.CatToCatCreateResponse(*cat)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("transaction %w", err)
	}

	return result, nil
}
