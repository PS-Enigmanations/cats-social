package controller

import (
	"encoding/json"
	"enigmanations/cats-social/internal/cat/service"
	"log"
	"net/http"

	"enigmanations/cats-social/internal/cat"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

type catController struct {
	Service service.CatService
}

func NewCatController(svc service.CatService) catController {
	return catController{Service: svc}
}

func (c *catController) CatGetController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	cats, err := c.Service.GetAll()
	jsonResp, err := json.Marshal(cats)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func (c *catController) CatCreateController(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var reqBody cat.CatCreateRequest

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

	// send data to service layer to further process (create record)
	cat, err := c.Service.Create(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	jsonResp, err := json.Marshal(cat)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}
