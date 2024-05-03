package service

import (
	"context"
	catmatch "enigmanations/cats-social/internal/cat_match"
	"enigmanations/cats-social/internal/cat_match/errs"
	"enigmanations/cats-social/internal/cat_match/repository"
	"enigmanations/cats-social/internal/cat_match/request"
	userRepository "enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/pkg/database"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatMatchService interface {
	Create(req *request.CatMatchRequest, actorId int64) error
	Approve(matchId int) error
	Reject(matchId int) error
	Destroy(id int64) error
	GetByIssuedOrReceiver(matchId int) ([]*catmatch.CatMatchResponse, error)
}

type CatMatchDependency struct {
	CatMatch repository.CatMatchRepository
	User     userRepository.UserRepository
}
type catMatchService struct {
	pool    *pgxpool.Pool
	repo    *CatMatchDependency
	Context context.Context
}

func NewCatMatchService(ctx context.Context, pool *pgxpool.Pool, repo *CatMatchDependency) CatMatchService {
	return &catMatchService{pool: pool, repo: repo, Context: ctx}
}

func (svc *catMatchService) validate(req *request.CatMatchRequest) error {
	repo := svc.repo

	// Check cat by match cat id
	matchCatFound, err := repo.CatMatch.GetAssociationByCatId(svc.Context, int(req.MatchCatId))
	if err != nil {
		return errs.CatMatchErrNotFound
	}
	if matchCatFound.AlreadyMatched {
		return errs.CatMatchErrAlreadyMatched
	}

	// Check user from match cat id is belong to the user
	_, err = repo.User.Get(svc.Context, matchCatFound.UserId)
	if err != nil {
		return errs.CatMatchErrOwnerNotFound
	}

	// Check cat by user cat id
	userCatFound, err := repo.CatMatch.GetAssociationByCatId(svc.Context, int(req.UserCatId))
	if err != nil {
		return errs.CatMatchErrNotFound
	}
	if userCatFound.AlreadyMatched {
		return errs.CatMatchErrAlreadyMatched
	}

	// Check user from user cat id is belong to the user
	_, err = repo.User.Get(svc.Context, userCatFound.UserId)
	if err != nil {
		return err
	}

	// Ensure cat owner between issuer -> receiver should be not equal
	if matchCatFound.UserId == userCatFound.UserId {
		return errs.CatMatchErrInvalidAuthor
	}
	// Ensure cat owner between receiver -> issuer should be not equal
	if userCatFound.UserId == matchCatFound.UserId {
		return errs.CatMatchErrInvalidAuthor
	}

	// Check gender, should be not equal
	equalMatches := reflect.DeepEqual(matchCatFound.Sex, userCatFound.Sex)
	if equalMatches {
		return errs.CatMatchErrGender
	}

	return nil
}

func (svc *catMatchService) Create(req *request.CatMatchRequest, actorId int64) error {
	repo := svc.repo

	// Validate first
	err := svc.validate(req)
	if err != nil {
		return err
	}

	if err := database.BeginTransaction(svc.Context, svc.pool, func(tx pgx.Tx) error {
		// Create cat matches
		model := catmatch.CatMatch{
			IssuedBy:   actorId,
			MatchCatId: req.MatchCatId,
			UserCatId:  req.UserCatId,
			Message:    req.Message,
		}
		err = repo.CatMatch.Save(svc.Context, &model, tx)
		if err != nil {
			return err
		}

		// Update cat already match if cat matches successfully created
		err = repo.CatMatch.UpdateCatAlreadyMatches(
			svc.Context,
			[]int{
				int(req.MatchCatId),
				int(req.UserCatId),
			},
			true,
			tx,
		)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction %w", err)
	}

	return nil
}

func (svc *catMatchService) Destroy(id int64) error {
	repo := svc.repo

	if err := database.BeginTransaction(svc.Context, svc.pool, func(tx pgx.Tx) error {
		// Destroy cat matches
		err := repo.CatMatch.Destroy(svc.Context, id, tx)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction %w", err)
	}

	return nil
}

func (svc *catMatchService) Approve(matchId int) error {
	repo := svc.repo

	if err := database.BeginTransaction(svc.Context, svc.pool, func(tx pgx.Tx) error {
		// Approve cat matches
		err := repo.CatMatch.UpdateCatMatchStatus(svc.Context, matchId, "success", tx)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction %w", err)
	}
	
	return nil
}

func (svc *catMatchService) Reject(matchId int) error {
	repo := svc.repo

	if err := database.BeginTransaction(svc.Context, svc.pool, func(tx pgx.Tx) error {
		// Reject cat matches
		err := repo.CatMatch.UpdateCatMatchStatus(svc.Context, matchId, "reject", tx)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction %w", err)
	}
	
	return nil
}

func (svc *catMatchService) GetByIssuedOrReceiver(matchId int) ([]*catmatch.CatMatchResponse, error) {
	repo := svc.repo

	data, err := repo.CatMatch.GetByIssuedOrReceiver(svc.Context, int(matchId))

	if err != nil {
		fmt.Print(err)
		return nil, errs.CatMatchErrOwnerNotFound
	}

	var cmRes []*catmatch.CatMatchResponse
	for _, cm := range data {
		cmRes = append(cmRes, catmatch.CatToCatResponse(*cm))
	}

	return cmRes, nil
}

