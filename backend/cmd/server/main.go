// Command server arranca el backend HTTP.
// Env vars: DATABASE_URL (requerida), SERVER_PORT, CORS_ALLOWED_ORIGIN.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/api"
	"github.com/ezeromanelli/northwind-cobranza/backend/internal/db"
)

func main() {
	dsn := mustGetenv("DATABASE_URL")
	port := getenv("SERVER_PORT", "8080")
	corsOrigin := getenv("CORS_ALLOWED_ORIGIN", "http://localhost:5173")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("db.NewPool: %v", err)
	}
	defer pool.Close()

	handler := api.NewRouter(pool, corsOrigin)

	addr := ":" + port
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("northwind backend listening on %s (CORS: %s)", addr, corsOrigin)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return v
}
