package controller

import (
	"encoding/json"
	"enigmanations/cats-social/helper"
	"enigmanations/cats-social/model/web"
	"enigmanations/cats-social/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CatControllerImpl struct {
	CatService service.CatService
}

func NewCatController(catService service.CatService) CatController {
	return &CatControllerImpl{
		CatService: catService,
	}
}

func (controller *CatControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(request.Body)
	catCreateRequest := web.CatCreateRequest{}
	err := decoder.Decode(&catCreateRequest)
	helper.PanicIfError(err)

	catResponse := controller.CatService.Create(request.Context(), catCreateRequest)
	webResponse := web.SuccessResponse{
		Message: "Success",
		Data:    catResponse,
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err = encoder.Encode(webResponse)
	helper.PanicIfError(err)
}

func (controller *CatControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	panic("not implemented") // TODO: Implement
}

func (controller *CatControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	panic("not implemented") // TODO: Implement
}

func (controller *CatControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	panic("not implemented") // TODO: Implement
}

func (controller *CatControllerImpl) GetAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	panic("not implemented") // TODO: Implement
}
