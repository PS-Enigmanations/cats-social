package router_v1

import (
	"context"
	"enigmanations/cats-social/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type V1Router interface {
	Load(r *gin.Engine, m middleware.Middleware)
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

func (v *v1Router) Load(router *gin.Engine, m middleware.Middleware) {
	// @see https://gin-gonic.com/docs/examples/grouping-routes/
	v1 := router.Group("/v1")
	{
		// Users api endpoint
		user := v1.Group("/user")
		{
			user.POST("/register", v.User.Controller.UserRegister)
			user.POST("/login", v.User.Controller.UserLogin)
		}

		// Cats api endpoint
		cat := v1.Group("/cat").Use(m.Auth.UseAuthMiddleware())
		{
			cat.GET("/", v.Cat.Controller.CatGetAllController)
			cat.POST("/", v.Cat.Controller.CatCreateController)
			cat.DELETE("/:id", v.Cat.Controller.CatDeleteController)
			cat.PUT("/:id", v.Cat.Controller.CatUpdateController)

			// Cat Match api endpoint
			cat.POST("/match", v.CatMatch.Controller.CatMatchCreate)
			cat.GET("/match", v.CatMatch.Controller.CatMatchGetAll)
			cat.DELETE("/match/:id", v.CatMatch.Controller.CatMatchDestroy)

			// Cat Match approve/reject api endpoint
			cat.POST("/match/approve", v.CatMatch.Controller.CatMatchApprove)
			cat.POST("/match/reject", v.CatMatch.Controller.CatMatchReject)
		}
	}
}
