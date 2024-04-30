package controller

import (
	"encoding/json"
	"enigmanations/cats-social/internal/user/errs"
	"enigmanations/cats-social/internal/user/request"
	"enigmanations/cats-social/internal/user/service"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

type UserController interface {
	UserRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	UserLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

type userController struct {
	Service     service.UserService
	AuthService service.UserAuthService
}

func NewUserController(svc service.UserService, authSvc service.UserAuthService) UserController {
	return &userController{Service: svc, AuthService: authSvc}
}

func (c *userController) UserRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reqBody request.UserRegisterRequest

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

	result, err := c.Service.Create(&reqBody)
	if err != nil {
		switch {
		case errors.Is(err, errs.UserErrEmailInvalidFormat):
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		case errors.Is(err, errs.UserErrEmailExist):
			http.Error(w, err.Error(), http.StatusConflict)
			break
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	jsonResp, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func (c *userController) UserLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reqBody request.UserLoginRequest

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

	result, err := c.AuthService.Login(&reqBody)
	if err != nil {
		switch {
		case errors.Is(err, errs.UserErrEmailInvalidFormat):
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		case errors.Is(err, errs.UserErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
			break
		case errors.Is(err, errs.WrongPassword):
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

	jsonResp, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

	return
}
