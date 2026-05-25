package api

import "net/http"

// healthHandler responde 200 si el server esta vivo.
// No checkea Postgres aposta: si el pool falla los handlers de negocio
// devuelven 500 con detalle; /health solo verifica que el proceso responde.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
