package service

import (
	"context"

	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/internal/cat/errs"
	"enigmanations/cats-social/internal/cat/repository"
	"enigmanations/cats-social/internal/cat/request"
	catImageRepository "enigmanations/cats-social/internal/cat_image/repository"
	catMatchRepository "enigmanations/cats-social/internal/cat_match/repository"
	"enigmanations/cats-social/util"

	"enigmanations/cats-social/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatService interface {
	GetAllByParams(p *request.CatGetAllQueryParams, ownerId int) ([]*cat.Cat, error)
	Create(p *request.CatCreateRequest, actorId int) <-chan util.Result[*cat.Cat]
	Update(p *request.CatUpdateRequest, id int) <-chan error
	Delete(id int) <-chan error
}

type CatDependency struct {
	Cat      repository.CatRepository
	CatImage catImageRepository.CatImageRepository
	CatMatch catMatchRepository.CatMatchRepository
}

type catService struct {
	repo    *CatDependency
	pool    *pgxpool.Pool
	context context.Context
}

// NewService creates an API service.
func NewCatService(ctx context.Context, pool *pgxpool.Pool, repo *CatDependency) CatService {
	return &catService{repo: repo, pool: pool, context: ctx}
}

func (svc *catService) GetAllByParams(p *request.CatGetAllQueryParams, ownerId int) ([]*cat.Cat, error) {
	repo := svc.repo

	cats, err := repo.Cat.GetAllByParams(svc.context, p, ownerId)
	if err != nil {
		return nil, err
	}

	return cats, nil
}

func (svc *catService) Create(payload *request.CatCreateRequest, actorId int) <-chan util.Result[*cat.Cat] {
	repo := svc.repo

	result := make(chan util.Result[*cat.Cat])
	go func() {
		if err := database.BeginTransaction(svc.context, svc.pool, func(tx pgx.Tx, ctx context.Context) error {
			model := cat.Cat{
				UserId:      actorId,
				Name:        payload.Name,
				Race:        cat.Race(payload.Race),
				Sex:         cat.Sex(payload.Sex),
				AgeInMonth:  payload.AgeInMonth,
				Description: payload.Description,
			}

			// call Create from repository/ datastore
			catCreated, err := repo.Cat.Save(ctx, tx, model)

			// if error occur, return nil for the response as well as return the error
			if err != nil {
				result <- util.Result[*cat.Cat]{
					Error: err,
				}
				return err
			}

			err = repo.CatImage.SaveImageUrls(ctx, tx, catCreated.Id, payload.ImageUrls)
			if err != nil {
				result <- util.Result[*cat.Cat]{
					Error: err,
				}
				return err
			}

			catCreated.ImageUrls = payload.ImageUrls

			result <- util.Result[*cat.Cat]{
				Result: catCreated,
			}
			close(result)

			return nil
		}); err != nil {
			result <- util.Result[*cat.Cat]{
				Error: err,
			}
		}
	}()

	return result
}

func (svc *catService) Delete(id int) <-chan error {
	repo := svc.repo

	result := make(chan error)
	go func() {
		if err := database.BeginTransaction(svc.context, svc.pool, func(tx pgx.Tx, ctx context.Context) error {
			// Find cat
			catFound, err := repo.Cat.FindById(ctx, id)
			if err != nil {
				result <- errs.CatErrNotFound
				return errs.CatErrNotFound
			}

			// Delete cat
			err = repo.Cat.Delete(ctx, tx, catFound.Id)
			if err != nil {
				result <- err
				return err
			}

			result <- nil
			close(result)

			return nil
		}); err != nil {
			result <- err
		}
	}()

	return result
}

func (svc *catService) Update(p *request.CatUpdateRequest, id int) <-chan error {
	repo := svc.repo

	result := make(chan error)
	go func() {
		if err := database.BeginTransaction(svc.context, svc.pool, func(tx pgx.Tx, ctx context.Context) error {
			// Find cat
			catFound, err := repo.Cat.FindById(ctx, id)
			if err != nil {
				result <- errs.CatErrNotFound
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
				catMatchFound, err := repo.CatMatch.GetByCatId(ctx, catFound.Id)
				if err != nil {
					result <- err
					return err
				}
				// If exists, sex should be not editable
				if catMatchFound != nil {
					if catFound.Sex != cat.Sex(p.Sex) {
						result <- errs.CatErrSexNotEditable
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
				err = repo.CatImage.SaveImageUrls(ctx, tx, payload.Id, payload.ImageUrls)
				if err != nil {
					result <- err
					return err
				}
			}

			// Update cat
			_, err = repo.Cat.Update(ctx, tx, payload)
			if err != nil {
				result <- err
				return err
			}

			result <- nil
			close(result)

			return nil
		}); err != nil {
			result <- err
		}
	}()

	return result
}
