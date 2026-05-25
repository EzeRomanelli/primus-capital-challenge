package api

import "net/http"

// No checkea Postgres a propósito: los handlers de negocio devuelven 500 si la
// DB falla; /health solo confirma que el proceso responde.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
