package main

import (
	"context"
	"enigmanations/cats-social/pkg/database"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env %v\n", err)
		os.Exit(1)
	}

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
}
