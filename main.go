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
)

func main() {
	// Load config
	cfg := config.GetConfig()

	// Shared ctx
	ctx := context.Background()

	// Connect to the database
	pgUrl := `postgres://%s:%s@%s:%d/%s?%s`
	pgUrl = fmt.Sprintf(pgUrl,
		cfg.DBUsername,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBParams,
	)

	pool, err := database.NewPGXPool(ctx, pgUrl, &database.PGXStdLogger{
		Logger: slog.Default(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Check reachability
	if _, err = pool.Exec(ctx, `SELECT 1`); err != nil {
		fmt.Errorf("pool.Exec() error: %v", err)
	}

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
