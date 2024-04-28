package app

import (
	"PS-Enigmanations/cats-social/helper"
	"database/sql"
)

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./cats.db")
	helper.PanicIfError(err)
	return db
}
