package controller

import (
	"enigmanations/cats-social/internal/cat/errs"
	"enigmanations/cats-social/internal/cat/request"
	"enigmanations/cats-social/internal/cat/response"
	"enigmanations/cats-social/internal/cat/service"
	"enigmanations/cats-social/internal/common/auth"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type CatController interface {
	CatGetAllController(ctx *gin.Context)
	CatCreateController(ctx *gin.Context)
	CatDeleteController(ctx *gin.Context)
	CatUpdateController(ctx *gin.Context)
}

type catController struct {
	Service service.CatService
}

type byIdRequest struct {
	ID int `uri:"id" binding:"required,min=1" example:"1"`
}

func NewCatController(svc service.CatService) CatController {
	return &catController{Service: svc}
}

func (c *catController) CatGetAllController(ctx *gin.Context) {
	var reqQueryParams request.CatGetAllQueryParams
	if err := ctx.ShouldBindQuery(&reqQueryParams); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	currUser := auth.GetCurrentUser(ctx)

	cats, err := c.Service.GetAllByParams(&reqQueryParams, currUser.Uid)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Mapping data from service to response
	catShows := response.ToCatShows(cats)
	catMappedResults := response.CatToCatGetAllResponse(catShows)

	ctx.JSON(http.StatusOK, catMappedResults)
	return
}

func (c *catController) CatCreateController(ctx *gin.Context) {
	var reqBody request.CatCreateRequest
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

	// send data to service layer to further process (create record)
	catCreated := <-c.Service.Create(&reqBody, currUser.Uid)
	if catCreated.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Mapping data from service to response
	catCreatedMappedResult := response.CatToCatCreateResponse(*catCreated.Result)

	ctx.JSON(http.StatusCreated, catCreatedMappedResult)
	return
}

func (c *catController) CatDeleteController(ctx *gin.Context) {
	// get cat id from request params
	var reqUri *byIdRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := <-c.Service.Delete(reqUri.ID)
	if err != nil {
		switch {
		case errors.Is(err, errs.CatErrNotFound):
			ctx.AbortWithError(http.StatusNotFound, err)
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

func (c *catController) CatUpdateController(ctx *gin.Context) {
	// get cat id from request params
	var reqUri *byIdRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var reqBody request.CatUpdateRequest
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

	err = <-c.Service.Update(&reqBody, reqUri.ID)
	if err != nil {
		switch {
		case errors.Is(err, errs.CatErrNotFound):
			ctx.AbortWithError(http.StatusNotFound, err)
			break
		case errors.Is(err, errs.CatErrSexNotEditable):
			ctx.AbortWithError(http.StatusBadRequest, err)
			break
		default:
			ctx.AbortWithError(http.StatusInternalServerError, err)
			break
		}
		return
	}

	ctx.Status(http.StatusOK)
	return
}
