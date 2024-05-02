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

## prod: run build and up on production environment.
prod: build-prod up

## build: run build on dev environment.
build:
	GOARCH=amd64 GOOS=darwin go build -o main .

## build: run build on dev environment.
build-prod:
	GOARCH=amd64 GOOS=linux go build -o main .

## up: run docker-compose up with dev environment.
up:
	JWT_SECRET=a-very-secretive-secret-key ./main

## up: run docker-compose up with dev environment.
up-prod:
	JWT_SECRET=a-very-secretive-secret-key ./main

## run golang-migrate up
migrateup:
	migrate -database $(DATABASE_URL) -path db/migrations up

## run golang-migrate up
migratedown:
	migrate -database $(DATABASE_URL) -path db/migrations down
