// Package api expone el HTTP router (chi) y los handlers de los 5 endpoints.
//
// Cada handler hace su validacion manual; no usamos go-playground/validator
// porque tenemos solo 2 endpoints de mutacion y la validacion explicita es
// mas leible que un setup de tags.
package api

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Router agrupa dependencias compartidas por los handlers para no
// pasarlas como argumentos de cada funcion.
type Router struct {
	pool              *pgxpool.Pool
	corsAllowedOrigin string
}

// NewRouter arma el router con middleware comun y las rutas de la API.
// El http.Handler resultante se monta directamente en http.Server.
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

	r.Route("/api", func(r chi.Router) {
		r.Get("/segmentos", rt.listSegmentosHandler)
		r.Get("/clientes", rt.listClientesHandler)
		r.Get("/clientes/{id}", rt.getClienteHandler)
		r.Post("/clientes/{id}/gestiones", rt.crearGestionHandler)
	})

	return r
}

// corsMiddleware permite el origen del frontend en desarrollo.
// Simple, sin dependencia externa para una sola politica.
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

// uuidPattern valida formato UUID v4-ish (8-4-4-4-12 hex con guiones).
// No validamos la version porque Postgres lo hace. Solo evitamos enviar
// strings basura a la DB y obtener errores cripticos.
var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func isValidUUID(s string) bool {
	return uuidPattern.MatchString(s)
}
