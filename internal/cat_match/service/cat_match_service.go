package service

import (
	"context"
	catmatch "enigmanations/cats-social/internal/cat_match"
	"enigmanations/cats-social/internal/cat_match/errs"
	"enigmanations/cats-social/internal/cat_match/repository"
	"enigmanations/cats-social/internal/cat_match/request"
	userRepository "enigmanations/cats-social/internal/user/repository"
	"reflect"
)

type CatMatchService interface {
	Create(req *request.CatMatchRequest, actorId int64) error
}

type CatMatchServiceDependency struct {
	CatMatch repository.CatMatchRepository
	User     userRepository.UserRepository
}
type catMatchService struct {
	repo    *CatMatchServiceDependency
	Context context.Context
}

func NewCatMatchService(repo *CatMatchServiceDependency, ctx context.Context) CatMatchService {
	return &catMatchService{repo: repo, Context: ctx}
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

	// Ensure cat owner between receiver -> issuer should be not equal
	if matchCatFound.UserId == userCatFound.UserId {
		return errs.CatMatchErrInvalidAuthor
	}
	// Ensure cat owner between issuer -> receiver should be not equal
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

	// Update already match
	err = repo.CatMatch.UpdateCatAlreadyMatches(
		svc.Context,
		[]int{
			int(req.MatchCatId),
			int(req.UserCatId),
		},
		true,
	)
	if err != nil {
		return err
	}

	model := catmatch.CatMatch{
		IssuedBy:   actorId,
		MatchCatId: req.MatchCatId,
		UserCatId:  req.UserCatId,
		Message:    req.Message,
	}
	err = repo.CatMatch.Save(svc.Context, &model)
	if err != nil {
		return err
	}

	return nil
}
