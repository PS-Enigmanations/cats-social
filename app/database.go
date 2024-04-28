package app

import (
	"database/sql"
	"enigmanations/cats-social/helper"
)

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./cats.db")
	helper.PanicIfError(err)
	return db
}
