package main

import (
	"context"
	//"enigmanations/cats-social/middleware"
	"enigmanations/cats-social/middleware"
	"enigmanations/cats-social/pkg/database"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	// Cat
	catControllerInternal "enigmanations/cats-social/internal/cat/controller"
	catRepositoryInternal "enigmanations/cats-social/internal/cat/repository"
	catServiceInternal "enigmanations/cats-social/internal/cat/service"

	// Cat Match
	catMatchControllerInternal "enigmanations/cats-social/internal/cat_match/controller"
	catMatchRepositoryInternal "enigmanations/cats-social/internal/cat_match/repository"
	catMatchServiceInternal "enigmanations/cats-social/internal/cat_match/service"

	// User
	userControllerInternal "enigmanations/cats-social/internal/user/controller"
	userRepositoryInternal "enigmanations/cats-social/internal/user/repository"
	userServiceInternal "enigmanations/cats-social/internal/user/service"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env %v\n", err)
		os.Exit(1)
	}

	// Connect to the database
	pgUrl := `postgres://%s:%s@%s:%d/%s?%s&pool_max_conns=%d`
	pgUrl = fmt.Sprintf(pgUrl,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		5432,
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PARAMS"),
		32,
	)

	pgPool, err := database.NewPGXPool(context.Background(), pgUrl, &database.PGXStdLogger{
		Logger: slog.Default(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pgPool.Close()

	// Check reachability
	if _, err = pgPool.Exec(context.Background(), `SELECT 1`); err != nil {
		fmt.Errorf("pool.Exec() error: %v", err)
	}

	// Shared ctx
	ctx := context.Background()

	// Prepare router
	router := httprouter.New()

	// Prepare middleware
	auth := middleware.NewAuthMiddleware(pgPool, ctx)

	// Users
	userRepository := userRepositoryInternal.NewUserRepository(pgPool)
	userService := userServiceInternal.NewUserService(
		ctx,
		pgPool,
		&userServiceInternal.UserDependency{
			User: userRepository,
		},
	)
	userAuthRepository := userRepositoryInternal.NewUserAuthRepository(pgPool)
	userAuthService := userServiceInternal.NewUserAuthService(
		ctx,
		pgPool,
		&userServiceInternal.UserAuthDependency{
			User:    userRepository,
			Session: userAuthRepository,
		},
	)
	userController := userControllerInternal.NewUserController(userService, userAuthService)

	// Users api endpoint
	router.POST("/v1/register", userController.UserRegister)
	router.POST("/v1/login", userController.UserLogin)

	// Cats
	catRepository := catRepositoryInternal.NewCatRepository(pgPool)
	catService := catServiceInternal.NewCatService(catRepository, ctx)
	catController := catControllerInternal.NewCatController(catService)

	// Cats api endpoint
	router.GET("/v1/cats", auth.ProtectedHandler(catController.CatGetController))
	router.POST("/v1/cats", auth.ProtectedHandler(catController.CatCreateController))
	router.DELETE("/v1/cats/:id", auth.ProtectedHandler(catController.CatDeleteController))

	// Cat Match
	catMatchRepository := catMatchRepositoryInternal.NewCatMatchRepository(pgPool)
	catMatchService := catMatchServiceInternal.NewCatMatchService(
		ctx,
		pgPool,
		&catMatchServiceInternal.CatMatchDependency{
			User:     userRepository,
			CatMatch: catMatchRepository,
		},
	)
	catMatchController := catMatchControllerInternal.NewCatMatchController(catMatchService)

	// Cat Match api endpoint
	router.POST("/v1/cat/match", auth.ProtectedHandler(catMatchController.CatMatchCreate))

	// Run the server
	appServeAddr := ":" + os.Getenv("APP_PORT")
	fmt.Printf("Serving on http://localhost:%s\n", os.Getenv("APP_PORT"))
	log.Fatalf("%v", http.ListenAndServe(appServeAddr, router))
}
