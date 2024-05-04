package controller

import (
	"encoding/json"
	"enigmanations/cats-social/internal/cat/errs"
	"enigmanations/cats-social/internal/cat/request"
	"enigmanations/cats-social/internal/cat/service"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/util"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
)

type CatController interface {
	CatGetAllController(w http.ResponseWriter, r *http.Request)
	CatCreateController(w http.ResponseWriter, r *http.Request)
	CatDeleteController(w http.ResponseWriter, r *http.Request)
	CatUpdateController(w http.ResponseWriter, r *http.Request)
}

type catController struct {
	Service service.CatService
}

func NewCatController(svc service.CatService) CatController {
	return &catController{Service: svc}
}

func (c *catController) CatGetAllController(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	queryParams, err := util.ParseQuery[request.CatGetAllQueryParams](r)
	if err != nil {
		log.Fatalf("Error happened in parse query. Err: %s", err)
	}
	currUser := user.GetCurrentUser(r.Context())
	cats, err := c.Service.GetAllByParams(queryParams, currUser.Uid)
	jsonResp, err := json.Marshal(cats)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func (c *catController) CatCreateController(w http.ResponseWriter, r *http.Request) {
	var reqBody request.CatCreateRequest

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

	// send data to service layer to further process (create record)
	cat, err := c.Service.Create(&reqBody, currUser.Uid)
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

func (c *catController) CatDeleteController(w http.ResponseWriter, r *http.Request) {
	// get cat id from request params
	id := r.URL.Query().Get(":id")

	catId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.Service.Delete(catId)
	if err != nil {
		switch {
		case errors.Is(err, errs.CatErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
			break
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	return
}

func (c *catController) CatUpdateController(w http.ResponseWriter, r *http.Request) {
	// get cat id from request params
	id := r.URL.Query().Get(":id")

	catId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var reqBody request.CatUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.Service.Update(&reqBody, catId)
	if err != nil {
		switch {
		case errors.Is(err, errs.CatErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
			break
		case errors.Is(err, errs.CatErrSexNotEditable):
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	return
}
