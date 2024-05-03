package controller

import (
	"encoding/json"
	"enigmanations/cats-social/internal/cat_match/errs"
	"enigmanations/cats-social/internal/cat_match/request"
	"enigmanations/cats-social/internal/cat_match/service"
	"enigmanations/cats-social/internal/user"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

type CatMatchController interface {
	CatMatchCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	CatMatchDestroy(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	CatMatchApprove(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	CatMatchReject(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

type catMatchController struct {
	Service service.CatMatchService
}

func NewCatMatchController(svc service.CatMatchService) CatMatchController {
	return &catMatchController{Service: svc}
}

func (c *catMatchController) CatMatchCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reqBody request.CatMatchRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currUser := user.GetCurrentUser(r.Context())

	if err = c.Service.Create(&reqBody, int64(currUser.Uid)); err != nil {
		switch {
		case errors.Is(err, errs.CatMatchErrNotFound),
			errors.Is(err, errs.CatMatchErrOwner):
			http.Error(w, err.Error(), http.StatusNotFound)
			break
		case errors.Is(err, errs.CatMatchErrGender),
			errors.Is(err, errs.CatMatchErrInvalidAuthor),
			errors.Is(err, errs.CatMatchErrAlreadyMatched):
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	return
}

func (c *catMatchController) CatMatchDestroy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	catId, err := strconv.Atoi(id)

	if err = c.Service.Destroy(int64(catId)); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	return
}


func (c *catMatchController) CatMatchApprove(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reqBody request.CatMatchApproveRejectRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.Service.Approve(int(reqBody.MatchId)); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	return
}

func (c *catMatchController) CatMatchReject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reqBody request.CatMatchApproveRejectRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.Service.Reject(int(reqBody.MatchId)); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	return
}
