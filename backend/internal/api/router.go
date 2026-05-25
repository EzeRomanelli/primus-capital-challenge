// Package api expone el HTTP router y los handlers.
package api

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	pool              *pgxpool.Pool
	corsAllowedOrigin string
}

func NewRouter(pool *pgxpool.Pool, corsAllowedOrigin string) http.Handler {
	rt := &Router{
		pool:              pool,
		corsAllowedOrigin: corsAllowedOrigin,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware(corsAllowedOrigin))

	r.Get("/health", healthHandler)

	r.Get("/swagger", http.RedirectHandler("/swagger/", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/swagger/", swaggerUIHandler)
	r.Get("/openapi.yaml", openapiYAMLHandler)

	r.Route("/api", func(r chi.Router) {
		r.Get("/segmentos", rt.listSegmentosHandler)
		r.Get("/clientes", rt.listClientesHandler)
		r.Get("/clientes/{id}", rt.getClienteHandler)
		r.Post("/clientes/{id}/gestiones", rt.crearGestionHandler)
	})

	return r
}

func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func isValidUUID(s string) bool {
	return uuidPattern.MatchString(s)
}
