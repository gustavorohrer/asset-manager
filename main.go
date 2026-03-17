package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gustavorohrer/ecl-be-challenge/internal/assets"
	"github.com/gustavorohrer/ecl-be-challenge/internal/httpapi"
	"github.com/gustavorohrer/ecl-be-challenge/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbPool, err := newPostgresPool()
	if err != nil {
		log.Fatalf("failed to initialize postgres pool: %v", err)
	}
	defer dbPool.Close()

	assetsRepository := postgres.NewAssetRepository(dbPool)
	assetsService := assets.NewService(assetsRepository)
	assetsHandler := httpapi.NewAssetsHandler(assetsService)

	router := httpapi.NewRouter(assetsHandler)
	if err := router.Run(serverAddress()); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func newPostgresPool() (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			getEnv("DB_USER", "applicant"),
			getEnv("DB_PASSWORD", "goodluck"),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_NAME", "eclypsiumdb"),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func serverAddress() string {
	return ":" + getEnv("PORT", "8080")
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
