package main

import (
	"context"
	"enigmanations/cats-social/pkg/database"
	"enigmanations/cats-social/pkg/env"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"enigmanations/cats-social/config"
	"enigmanations/cats-social/middleware"
	routes "enigmanations/cats-social/router"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load config
	cfg := config.GetConfig()

	// Shared ctx
	ctx := context.Background()

	// Connect to the database
	pool := initDatabase(cfg, ctx)
	defer pool.Close()

	// Prepare middleware
	middleware := middleware.NewMiddleware(pool)

	// Disable debug mode in production
	if env.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Prepare router
	router := gin.New()

	// Register routes
	routes.RegisterRouter(ctx, pool, router, middleware)

	// Run the server
	appServeAddr := ":" + fmt.Sprint(cfg.AppPort)
	fmt.Printf("Serving on http://localhost:%s\n", fmt.Sprint(cfg.AppPort))
	log.Fatalf("%v", http.ListenAndServe(appServeAddr, router))
}

func initDatabase(cfg *config.Configuration, ctx context.Context) *pgxpool.Pool {
	pgUrl := `postgres://%s:%s@%s:%d/%s?%s&pool_max_conns=%d`
	pgUrl = fmt.Sprintf(pgUrl,
		cfg.DBUsername,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBParams,
		32,
	)

	pgPool, err := database.NewPGXPool(ctx, pgUrl, &database.PGXStdLogger{
		Logger: slog.Default(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Check reachability
	if _, err = pgPool.Exec(ctx, `SELECT 1`); err != nil {
		fmt.Errorf("pool.Exec() error: %v", err)
	}

	return pgPool
}
