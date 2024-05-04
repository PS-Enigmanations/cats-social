package router_v1

import (
	"context"
	"enigmanations/cats-social/middleware"

	"net/http"

	"github.com/bmizerany/pat"
	"github.com/jackc/pgx/v5/pgxpool"
)

type V1Router interface {
	Load(r *pat.PatternServeMux, m middleware.Middleware)
}

type v1Router struct {
	User     *UserRouter
	Cat      *CatRouter
	CatMatch *CatMatchRouter
}

func NewV1Router(ctx context.Context, pool *pgxpool.Pool) *v1Router {
	return &v1Router{
		User:     NewUserRouter(ctx, pool),
		Cat:      NewCatRouter(ctx, pool),
		CatMatch: NewCatMatchRouter(ctx, pool),
	}
}

func (v *v1Router) Load(router *pat.PatternServeMux, m middleware.Middleware) {
	// Users api endpoint
	router.Post("/v1/user/register", http.HandlerFunc(v.User.Controller.UserRegister))
	router.Post("/v1/user/login", http.HandlerFunc(v.User.Controller.UserLogin))

	// Cats api endpoint
	router.Get("/v1/cat", m.Auth.ProtectedHandler(http.HandlerFunc(v.Cat.Controller.CatGetAllController)))
	router.Post("/v1/cat", m.Auth.ProtectedHandler(http.HandlerFunc(v.Cat.Controller.CatCreateController)))
	router.Del("/v1/cat/:id", m.Auth.ProtectedHandler(http.HandlerFunc(v.Cat.Controller.CatDeleteController)))
	router.Put("/v1/cat/:id", m.Auth.ProtectedHandler(http.HandlerFunc(v.Cat.Controller.CatUpdateController)))

	// Cat Match api endpoint
	router.Post("/v1/cat/match", m.Auth.ProtectedHandler(http.HandlerFunc(v.CatMatch.Controller.CatMatchCreate)))
	router.Get("/v1/cat/match", m.Auth.ProtectedHandler(http.HandlerFunc(v.CatMatch.Controller.CatMatchGetAll)))
	router.Post("/v1/cat/match/approve", m.Auth.ProtectedHandler(http.HandlerFunc(v.CatMatch.Controller.CatMatchApprove)))
	router.Post("/v1/cat/match/reject", m.Auth.ProtectedHandler(http.HandlerFunc(v.CatMatch.Controller.CatMatchReject)))
	router.Del("/v1/cat/match/:id", m.Auth.ProtectedHandler(http.HandlerFunc(v.CatMatch.Controller.CatMatchDestroy)))
}
