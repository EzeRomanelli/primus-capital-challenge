package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/api"
	"github.com/ezeromanelli/northwind-cobranza/backend/tests"
)

// TestListClientesPriorizado es el integration test del endpoint mas critico
// del MVP. Cubre el flujo end-to-end: DB real -> query con joins -> calculo
// del score en Go -> respuesta JSON ordenada por score DESC.
//
// Si la DB de tests no esta disponible, el helper hace t.Skip; el test
// no rompe el build.
func TestListClientesPriorizado(t *testing.T) {
	pool := tests.MustTestPool(t)
	defer pool.Close()

	tests.TruncateAll(t, pool)
	_, zombiID := tests.SeedDosClientes(t, pool)

	handler := api.NewRouter(pool, "*")
	req := httptest.NewRequest(http.MethodGet, "/api/clientes", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200. body: %s", rec.Code, rec.Body.String())
	}

	var got []api.ClientePriorizadoDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v. body: %s", err, rec.Body.String())
	}
	if len(got) != 2 {
		t.Fatalf("esperaba 2 clientes, got %d", len(got))
	}

	// El zombi debe ser primero (mayor score).
	if got[0].ID != zombiID {
		t.Errorf("primer cliente: got id=%s (score=%d), want zombi id=%s",
			got[0].ID, got[0].Score, zombiID)
	}
	if got[0].Score <= got[1].Score {
		t.Errorf("zombi score (%d) deberia ser > saludable score (%d)",
			got[0].Score, got[1].Score)
	}
	if got[0].Score == 0 {
		t.Errorf("zombi tiene score 0; algo fallo en el calculo")
	}
}

