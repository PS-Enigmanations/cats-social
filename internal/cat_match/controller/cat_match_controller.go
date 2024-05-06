package controller

import (
	"enigmanations/cats-social/internal/cat_match/errs"
	"enigmanations/cats-social/internal/cat_match/request"
	"enigmanations/cats-social/internal/cat_match/response"
	"enigmanations/cats-social/internal/cat_match/service"
	"enigmanations/cats-social/internal/common/auth"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type CatMatchController interface {
	CatMatchCreate(ctx *gin.Context)
	CatMatchDestroy(ctx *gin.Context)
	CatMatchApprove(ctx *gin.Context)
	CatMatchReject(ctx *gin.Context)
	CatMatchGetAll(ctx *gin.Context)
}

type catMatchController struct {
	Service service.CatMatchService
}

type byIdRequest struct {
	ID int `uri:"id" binding:"required,min=1" example:"1"`
}

func NewCatMatchController(svc service.CatMatchService) CatMatchController {
	return &catMatchController{Service: svc}
}

func (c *catMatchController) CatMatchCreate(ctx *gin.Context) {
	var reqBody request.CatMatchRequest
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

	currUser := auth.GetCurrentUser(ctx)

	if err = c.Service.Create(&reqBody, int64(currUser.Uid)); err != nil {
		switch {
		case errors.Is(err, errs.CatMatchErrNotFound),
			errors.Is(err, errs.CatMatchErrOwner):
			ctx.AbortWithError(http.StatusNotFound, err)
			break
		case errors.Is(err, errs.CatMatchErrGender),
			errors.Is(err, errs.CatMatchErrInvalidAuthor),
			errors.Is(err, errs.CatMatchErrAlreadyMatched):
			ctx.AbortWithError(http.StatusBadRequest, err)
			break
		default:
			ctx.AbortWithError(http.StatusInternalServerError, err)
			break
		}
		return
	}

	ctx.Status(http.StatusCreated)
	return
}

func (c *catMatchController) CatMatchDestroy(ctx *gin.Context) {
	var reqUri *byIdRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := c.Service.Destroy(int64(reqUri.ID)); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
	return
}

func (c *catMatchController) CatMatchApprove(ctx *gin.Context) {
	var reqBody request.CatMatchApproveRejectRequest
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

	if err = c.Service.Approve(int(reqBody.MatchId)); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusCreated)
	return
}

func (c *catMatchController) CatMatchReject(ctx *gin.Context) {
	var reqBody request.CatMatchApproveRejectRequest
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

	if err = c.Service.Reject(int(reqBody.MatchId)); err != nil {
		return
	}

	ctx.Status(http.StatusOK)
	return
}

func (c *catMatchController) CatMatchGetAll(ctx *gin.Context) {
	currUser := auth.GetCurrentUser(ctx)
	catMatches, err := c.Service.GetByIssuedOrReceiver(int(currUser.Uid))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Mapping data from service to response
	catMatchShows := response.ToCatMatchShows(catMatches.CatMatches)
	catMatchMappedResults := response.CatToCatGetAllResponse(catMatchShows)

	// Write the JSON response
	ctx.JSON(http.StatusOK, catMatchMappedResults)
	return
}
