package main

import (
	"enigmanations/cats-social/app"
	"enigmanations/cats-social/controller"
	"enigmanations/cats-social/repository"
	"enigmanations/cats-social/service"

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
