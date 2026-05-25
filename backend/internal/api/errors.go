package api

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse es el formato JSON estable que devolvemos en cualquier 4xx/5xx.
// `code` es un identificador maquina-amigable (ej: "invalid_segmento") y
// `error` es el texto humano. Documentado en docs/API.md.
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// writeJSON serializa un valor como JSON y setea Content-Type + status.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError es el helper estandar para responder errores.
// El status NO se infiere del code: el caller pasa ambos explicitamente.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{Error: message, Code: code})
}
