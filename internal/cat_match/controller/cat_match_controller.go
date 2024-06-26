package controller

import (
	"encoding/json"
	"enigmanations/cats-social/internal/cat_match/errs"
	"enigmanations/cats-social/internal/cat_match/request"
	"enigmanations/cats-social/internal/cat_match/response"
	"enigmanations/cats-social/internal/cat_match/service"
	"enigmanations/cats-social/internal/session"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
)

type CatMatchController interface {
	CatMatchCreate(w http.ResponseWriter, r *http.Request)
	CatMatchDestroy(w http.ResponseWriter, r *http.Request)
	CatMatchApprove(w http.ResponseWriter, r *http.Request)
	CatMatchReject(w http.ResponseWriter, r *http.Request)
	CatMatchGetAll(w http.ResponseWriter, r *http.Request)
}

type catMatchController struct {
	Service service.CatMatchService
}

func NewCatMatchController(svc service.CatMatchService) CatMatchController {
	return &catMatchController{Service: svc}
}

func (c *catMatchController) CatMatchCreate(w http.ResponseWriter, r *http.Request) {
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

	currUser := session.GetCurrentUser(r.Context())

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

func (c *catMatchController) CatMatchDestroy(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	catId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = c.Service.Destroy(int64(catId)); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	return
}

func (c *catMatchController) CatMatchApprove(w http.ResponseWriter, r *http.Request) {
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

func (c *catMatchController) CatMatchReject(w http.ResponseWriter, r *http.Request) {
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

func (c *catMatchController) CatMatchGetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	currUser := session.GetCurrentUser(r.Context())
	catMatches, err := c.Service.GetByIssuedOrReceiver(int(currUser.Uid))
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	// Mapping data from service to response
	catMatchShows := response.ToCatMatchShows(catMatches.CatMatches)
	catMatchMappedResults := response.CatToCatGetAllResponse(catMatchShows)

	// Marshal the response into JSON
	jsonResp, err := json.Marshal(catMatchMappedResults)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	// Write the JSON response
	w.Write(jsonResp)
}
