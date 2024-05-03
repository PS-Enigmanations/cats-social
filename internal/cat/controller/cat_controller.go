package controller

import (
	"encoding/json"
	"enigmanations/cats-social/internal/cat/service"
	"log"
	"net/http"
	"strconv"

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
	w.WriteHeader(http.StatusOK)
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

func (c *catController) CatUpdateController(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var reqBody cat.CatUpdateRequest
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

	_, err = c.Service.FindById(reqBody.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cat, err := c.Service.Update(&reqBody, reqBody.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	jsonResp, err := json.Marshal(cat)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func (c *catController) CatDeleteController(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// get cat id from request params
	catId := p.ByName("id")
	id, err := strconv.Atoi(catId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.Service.FindById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = c.Service.Delete(id)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
