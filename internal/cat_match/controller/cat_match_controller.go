package controller

import (
	"encoding/json"
	"enigmanations/cats-social/internal/cat_match/request"
	"enigmanations/cats-social/internal/cat_match/service"
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

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	return
}
