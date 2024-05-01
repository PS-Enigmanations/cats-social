package controller

import (
	"encoding/json"
	"enigmanations/cats-social/internal/cat_match/errs"
	"enigmanations/cats-social/internal/cat_match/request"
	"enigmanations/cats-social/internal/cat_match/service"
	"enigmanations/cats-social/internal/user"
	"errors"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

type CatMatchController interface {
	CatMatchCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
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
	err = c.Service.Create(&reqBody, int64(currUser.Uid))
	if err != nil {
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
