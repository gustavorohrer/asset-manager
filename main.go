package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
	server := &http.Server{
		Addr:              serverAddress(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	shutdownSignalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err, ok := <-serverErrCh:
		if ok && err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	case <-shutdownSignalCtx.Done():
		log.Printf("shutdown signal received")
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to gracefully shutdown server: %v", err)
	}
}

func newPostgresPool() (*pgxpool.Pool, error) {
	databaseURL, err := resolveDatabaseURL()
	if err != nil {
		return nil, err
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

func resolveDatabaseURL() (string, error) {
	if value := strings.TrimSpace(os.Getenv("DATABASE_URL")); value != "" {
		return value, nil
	}

	dbHost := strings.TrimSpace(os.Getenv("DB_HOST"))
	dbPort := strings.TrimSpace(os.Getenv("DB_PORT"))
	dbName := strings.TrimSpace(os.Getenv("DB_NAME"))
	dbUser := strings.TrimSpace(os.Getenv("DB_USER"))
	dbPassword := strings.TrimSpace(os.Getenv("DB_PASSWORD"))

	missing := make([]string, 0, 5)
	if dbHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if dbPort == "" {
		missing = append(missing, "DB_PORT")
	}
	if dbName == "" {
		missing = append(missing, "DB_NAME")
	}
	if dbUser == "" {
		missing = append(missing, "DB_USER")
	}
	if dbPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}

	if len(missing) > 0 {
		return "", fmt.Errorf(
			"database configuration is incomplete: missing %s; set DATABASE_URL or all DB_* variables",
			strings.Join(missing, ", "),
		)
	}

	databaseURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPassword),
		Host:   net.JoinHostPort(dbHost, dbPort),
		Path:   "/" + dbName,
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String(), nil
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
