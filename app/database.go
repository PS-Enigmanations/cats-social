package app

import (
	"database/sql"
	"enigmanations/cats-social/helper"
	"github.com/mattn/go-sqlite3"
)

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./cats.db")
	helper.PanicIfError(err)
	return db
}
