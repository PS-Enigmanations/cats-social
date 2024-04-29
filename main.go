package main

import (
	"context"
	"enigmanations/cats-social/pkg/database"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"enigmanations/cats-social/internal/cat/controller"
	"enigmanations/cats-social/internal/cat/repository"
	"enigmanations/cats-social/internal/cat/service"

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
	pgUrl := `postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d`
	pgUrl = fmt.Sprintf(pgUrl,
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		5432,
		os.Getenv("DATABASE_NAME"),
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

	// Setup server
	catRepository := repository.NewCatRepository(pgPool)
	catService := service.NewCatService(catRepository, context.Background())
	catController := controller.NewCatController(catService)

	// Prepare router
	router := httprouter.New()

	// Cat api endpoint
	router.GET("/v1/cats", catController.CatGetController)
	router.POST("/v1/cats", catController.CatCreateController)

	// Run the server
	appServeAddr := ":" + os.Getenv("APP_PORT")
	fmt.Printf("Serving on http://localhost:%s\n", os.Getenv("APP_PORT"))
	log.Fatalf("%v", http.ListenAndServe(appServeAddr, router))
}
