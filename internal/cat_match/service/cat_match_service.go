package service

import (
	"context"
	catmatch "enigmanations/cats-social/internal/cat_match"
	"enigmanations/cats-social/internal/cat_match/errs"
	"enigmanations/cats-social/internal/cat_match/repository"
	"enigmanations/cats-social/internal/cat_match/request"
	userRepository "enigmanations/cats-social/internal/user/repository"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CatMatchService interface {
	Create(req *request.CatMatchRequest, actorId int64) error
	Approve(matchId int) error
	Reject(matchId int) error
	Destroy(id int64) error
	GetByIssuedOrReceiver(matchId int) (*getByIssueOrReceiverReturn, error)
}

type CatMatchDependency struct {
	CatMatch repository.CatMatchRepository
	User     userRepository.UserRepository
}
type catMatchService struct {
	pool    *pgxpool.Pool
	repo    *CatMatchDependency
	context context.Context
}

func NewCatMatchService(ctx context.Context, pool *pgxpool.Pool, repo *CatMatchDependency) CatMatchService {
	return &catMatchService{pool: pool, repo: repo, context: ctx}
}

func (svc *catMatchService) validate(req *request.CatMatchRequest) error {
	repo := svc.repo

	// Check cat by match cat id
	matchCatFound, err := repo.CatMatch.GetAssociationByCatId(svc.context, int(req.MatchCatId))
	if err != nil {
		return errs.CatMatchErrNotFound
	}
	if matchCatFound.AlreadyMatched {
		return errs.CatMatchErrAlreadyMatched
	}

	// Check user from match cat id is belong to the user
	_, err = repo.User.Get(svc.context, matchCatFound.UserId)
	if err != nil {
		return errs.CatMatchErrOwnerNotFound
	}

	// Check cat by user cat id
	userCatFound, err := repo.CatMatch.GetAssociationByCatId(svc.context, int(req.UserCatId))
	if err != nil {
		return errs.CatMatchErrNotFound
	}
	if userCatFound.AlreadyMatched {
		return errs.CatMatchErrAlreadyMatched
	}

	// Check user from user cat id is belong to the user
	_, err = repo.User.Get(svc.context, userCatFound.UserId)
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

	// Create cat matches
	model := catmatch.CatMatch{
		IssuedBy:   actorId,
		MatchCatId: req.MatchCatId,
		UserCatId:  req.UserCatId,
		Message:    req.Message,
	}
	err = repo.CatMatch.Save(svc.context, &model)
	if err != nil {
		return err
	}

	return nil
}

func (svc *catMatchService) Destroy(id int64) error {
	repo := svc.repo

	// Destroy cat matches
	err := repo.CatMatch.Destroy(svc.context, id)
	if err != nil {
		return err
	}

	return nil
}

func (svc *catMatchService) Approve(matchId int) error {
	repo := svc.repo

	// Get cat match
	catMatchFound, err := repo.CatMatch.Get(svc.context, matchId)
	if err != nil {
		return err
	}

	// Approve cat matches
	err = repo.CatMatch.UpdateCatMatchStatus(svc.context, catMatchFound.Id, "success")
	if err != nil {
		return err
	}

	// Update cat already match if cat matches approved
	err = repo.CatMatch.UpdateCatAlreadyMatches(
		svc.context,
		[]int{
			int(catMatchFound.UserCatId),
			int(catMatchFound.MatchCatId),
		},
		true,
	)
	if err != nil {
		return err
	}

	return nil
}

func (svc *catMatchService) Reject(matchId int) error {
	repo := svc.repo

	// Reject cat matches
	err := repo.CatMatch.UpdateCatMatchStatus(svc.context, matchId, "reject")
	if err != nil {
		return err
	}

	return nil
}

type getByIssueOrReceiverReturn struct {
	CatMatches []*catmatch.CatMatchValue
}

func (svc *catMatchService) GetByIssuedOrReceiver(matchId int) (*getByIssueOrReceiverReturn, error) {
	repo := svc.repo

	data, err := repo.CatMatch.GetByIssuedOrReceiver(svc.context, int(matchId))
	if err != nil {
		fmt.Print(err)
		return nil, errs.CatMatchErrOwnerNotFound
	}

	return &getByIssueOrReceiverReturn{
		CatMatches: data,
	}, nil
}
