package v1

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type Connection struct {
	pool *pgxpool.Pool
}

func NewHttpServerAPI(db Connection) {
	router := httprouter.New()

}
