package api

import (
	"net/http"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/db"
)

func (rt *Router) listSegmentosHandler(w http.ResponseWriter, r *http.Request) {
	segs, err := db.ListSegmentos(r.Context(), rt.pool)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "db_error", "no se pudieron obtener segmentos")
		return
	}
	writeJSON(w, http.StatusOK, segs)
}
