-include .env

ADDR := localhost:8000
PROJECTNAME := $(shell basename "$(PWD)")
DATABASE_URL := "postgres://postgres:postgres@localhost:5432/cats-social?sslmode=disable"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-get

## dev: run build and up on dev environment.
dev: build up

## build: run build on dev environment.
build:
	go build .

## up: run docker-compose up with dev environment.
up:
	SECRET_KEY=a-very-secretive-secret-key ./cats-social

## run golang-migrate up
migrateup:
	migrate -database $(DATABASE_URL) -path db/migrations up

## run golang-migrate up
migratedown:
	migrate -database $(DATABASE_URL) -path db/migrations down
