package controller

import (
	"encoding/json"
	"enigmanations/cats-social/internal/cat/service"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type catController struct {
	Service service.CatService
}

func NewCatController(svc service.CatService) catController {
	return catController{Service: svc}
}

func (c *catController) CatGetController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	cats, err := c.Service.GetAll()
	jsonResp, err := json.Marshal(cats)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}
