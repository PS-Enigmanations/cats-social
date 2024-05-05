package controller

import (
	"encoding/json"
	"enigmanations/cats-social/internal/user/errs"
	"enigmanations/cats-social/internal/user/request"
	"enigmanations/cats-social/internal/user/response"
	"enigmanations/cats-social/internal/user/service"
	"enigmanations/cats-social/util"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)

type UserController interface {
	UserRegister(w http.ResponseWriter, r *http.Request)
	UserLogin(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	Service     service.UserService
	AuthService service.UserAuthService
}

func NewUserController(svc service.UserService, authSvc service.UserAuthService) UserController {
	return &userController{Service: svc, AuthService: authSvc}
}

func (c *userController) UserRegister(w http.ResponseWriter, r *http.Request) {
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

	userCreated, err := c.Service.Create(&reqBody)
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

	// Mapping data from service to response
	userCreatedMappedResult := response.UserToUserCreateResponse(*userCreated.User, userCreated.AccessToken)

	// Marshal the response into JSON
	jsonResp, err := json.Marshal(userCreatedMappedResult)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func (c *userController) UserLogin(w http.ResponseWriter, r *http.Request) {
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

	// Validate email format
	if !util.IsEmail(reqBody.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	userLogged, err := c.AuthService.Login(&reqBody)
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

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	// Mapping data from service to response
	userLoggedMappedResult := response.UserToUserLoginResponse(*userLogged.User, userLogged.AccessToken)

	// Marshal the response into JSON
	jsonResp, err := json.Marshal(userLoggedMappedResult)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

	return
}
