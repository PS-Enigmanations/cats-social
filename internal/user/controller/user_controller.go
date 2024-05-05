package controller

import (
	"encoding/json"
	sessionErrs "enigmanations/cats-social/internal/session/errs"
	sessionRequest "enigmanations/cats-social/internal/session/request"
	sessionResponse "enigmanations/cats-social/internal/session/response"
	sessionService "enigmanations/cats-social/internal/session/service"
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
	UserService    service.UserService
	SessionService sessionService.SessionService
}

func NewUserController(userSvc service.UserService, sessionSvc sessionService.SessionService) UserController {
	return &userController{UserService: userSvc, SessionService: sessionSvc}
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

	userCreated, err := c.UserService.Create(&reqBody)
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
	var reqBody sessionRequest.SessionLoginRequest

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

	userLogged, err := c.SessionService.Login(&reqBody)
	if err != nil {
		switch {
		case errors.Is(err, errs.UserErrEmailInvalidFormat):
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		case errors.Is(err, errs.UserErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
			break
		case errors.Is(err, sessionErrs.WrongPassword):
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
	userLoggedMappedResult := sessionResponse.SessionToLoginResponse(*userLogged.User, userLogged.AccessToken)

	// Marshal the response into JSON
	jsonResp, err := json.Marshal(userLoggedMappedResult)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

	return
}
