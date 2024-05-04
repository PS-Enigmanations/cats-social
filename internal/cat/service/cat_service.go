package service

import (
	"context"
	"database/sql"
	"fmt"

	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/internal/cat/errs"
	"enigmanations/cats-social/internal/cat/repository"
	"enigmanations/cats-social/internal/cat/request"
	"enigmanations/cats-social/internal/cat/response"
	catMatchRepository "enigmanations/cats-social/internal/cat_match/repository"

	"enigmanations/cats-social/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatService interface {
	GetAllByParams(p *request.CatGetAllQueryParams, ownerId int) (*response.CatGetAllResponse, error)
	Create(p *request.CatCreateRequest, actorId int) (*response.CatCreateResponse, error)
	Update(p *request.CatUpdateRequest, id int) error
	Delete(id int) error
}

type CatDependency struct {
	Cat      repository.CatRepository
	CatMatch catMatchRepository.CatMatchRepository
}

type catService struct {
	repo    *CatDependency
	pool    *pgxpool.Pool
	Context context.Context
}

// NewService creates an API service.
func NewCatService(ctx context.Context, pool *pgxpool.Pool, repo *CatDependency) CatService {
	return &catService{repo: repo, pool: pool, Context: ctx}
}

func (svc *catService) GetAllByParams(p *request.CatGetAllQueryParams, ownerId int) (*response.CatGetAllResponse, error) {
	repo := svc.repo

	cats, err := repo.Cat.GetAllByParams(svc.Context, p, ownerId)

	if err != nil {
		return nil, err
	}

	catShows := response.ToCatShows(cats)
	return response.CatToCatGetAllResponse(catShows), nil
}

func (svc *catService) Create(payload *request.CatCreateRequest, actorId int) (*response.CatCreateResponse, error) {
	repo := svc.repo

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
		cat, err := repo.Cat.Save(svc.Context, tx, model)

		// if error occur, return nil for the response as well as return the error
		if err != nil {
			return nil
		}

		err = repo.Cat.SaveImageUrls(svc.Context, tx, cat.Id, payload.ImageUrls)
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

func (svc *catService) Delete(id int) error {
	repo := svc.repo

	if err := database.BeginTransaction(svc.Context, svc.pool, func(tx pgx.Tx) error {
		// Find cat
		catFound, err := repo.Cat.FindById(svc.Context, id)
		if err != nil {
			return errs.CatErrNotFound
		}

		// Delete cat
		err = repo.Cat.Delete(svc.Context, tx, catFound.Id)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("Delete transaction %w", err)
	}

	return nil
}

func (svc *catService) Update(p *request.CatUpdateRequest, id int) error {
	repo := svc.repo

	if err := database.BeginTransaction(svc.Context, svc.pool, func(tx pgx.Tx) error {
		// Find cat
		catFound, err := repo.Cat.FindById(svc.Context, id)
		if err != nil {
			return errs.CatErrNotFound
		}

		var payload = *catFound
		payload.Id = catFound.Id

		if p.Name != "" {
			payload.Name = p.Name
		}
		if p.Race != "" {
			payload.Race = cat.Race(p.Race)
		}
		if p.Sex != "" {
			// Check requested cat match for this cat is exists
			catMatchFound, err := repo.CatMatch.GetByCatId(svc.Context, catFound.Id)
			if err != nil && err == sql.ErrNoRows {
				return err
			}
			// If exists, sex should be not editable
			if catMatchFound != nil {
				if catFound.Sex != cat.Sex(p.Sex) {
					return errs.CatErrSexNotEditable
				}
			}
			payload.Sex = cat.Sex(p.Sex)
		}
		if p.AgeInMonth != 0 {
			payload.AgeInMonth = p.AgeInMonth
		}
		if p.Description != "" {
			payload.Description = p.Description
		}
		if len(p.ImageUrls) != 0 {
			payload.ImageUrls = p.ImageUrls

			// Currently we always create new record instead of deleted
			err = repo.Cat.SaveImageUrls(svc.Context, tx, payload.Id, payload.ImageUrls)
			if err != nil {
				return err
			}
		}

		// Update cat
		_, err = repo.Cat.Update(svc.Context, tx, payload)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("Update transaction %w", err)
	}

	return nil
}
