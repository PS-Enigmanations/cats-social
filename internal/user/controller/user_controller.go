package controller

import (
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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserController interface {
	UserRegister(ctx *gin.Context)
	UserLogin(ctx *gin.Context)
}

type userController struct {
	UserService    service.UserService
	SessionService sessionService.SessionService
}

func NewUserController(userSvc service.UserService, sessionSvc sessionService.SessionService) UserController {
	return &userController{UserService: userSvc, SessionService: sessionSvc}
}

func (c *userController) UserRegister(ctx *gin.Context) {
	var reqBody request.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqBody)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userCreated := <-c.UserService.Create(&reqBody)
	if userCreated.Error != nil {
		switch {
		case errors.Is(userCreated.Error, errs.UserErrEmailInvalidFormat):
			ctx.AbortWithError(http.StatusBadRequest, userCreated.Error)
			break
		case errors.Is(userCreated.Error, errs.UserErrEmailExist):
			ctx.AbortWithError(http.StatusConflict, userCreated.Error)
			break
		default:
			ctx.AbortWithError(http.StatusInternalServerError, userCreated.Error)
			break
		}
		return
	}

	// Mapping data from service to response
	userCreatedMappedResult := response.UserToUserCreateResponse(
		*userCreated.Result.User,
		userCreated.Result.AccessToken,
	)

	ctx.JSON(http.StatusCreated, userCreatedMappedResult)
	return
}

func (c *userController) UserLogin(ctx *gin.Context) {
	var reqBody sessionRequest.SessionLoginRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqBody)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Validate email format
	if !util.IsEmail(reqBody.Email) {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("Invalid email format"))
		return
	}

	userLogged, err := c.SessionService.Login(&reqBody)
	if err != nil {
		switch {
		case errors.Is(err, errs.UserErrEmailInvalidFormat):
			ctx.AbortWithError(http.StatusBadRequest, err)
			break
		case errors.Is(err, errs.UserErrNotFound):
			ctx.AbortWithError(http.StatusNotFound, err)
			break
		case errors.Is(err, sessionErrs.WrongPassword):
			ctx.AbortWithError(http.StatusBadRequest, err)
			break
		default:
			ctx.AbortWithError(http.StatusInternalServerError, err)
			break
		}
		return
	}

	// Mapping data from service to response
	userLoggedMappedResult := sessionResponse.SessionToLoginResponse(*userLogged.User, userLogged.AccessToken)

	ctx.JSON(http.StatusOK, userLoggedMappedResult)
	return
}
