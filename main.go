package main

import (
	"PS-Enigmanations/cats-social/app"
	"PS-Enigmanations/cats-social/controller"
	"PS-Enigmanations/cats-social/repository"
	"PS-Enigmanations/cats-social/service"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

func main() {
	db := app.NewDB()
	validate := validator.New()
	catRepository := repository.NewCatRepository()
	catService := service.NewCatService(catRepository, db, validate)
	catController := controller.NewCatController(catService)

	router := httprouter.New()

	router.POST("/cats", catController.Create)
}
