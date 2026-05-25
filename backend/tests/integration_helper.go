// Package tests provee helpers para los tests de integracion que necesitan
// una conexion real a Postgres.
//
// Los tests usan la DB apuntada por TEST_DATABASE_URL (default
// postgres://.../northwind_test). Setup:
//
//	make db-test-up      # crea la DB de tests y aplica migraciones
//	make test            # corre todo (unit + integration)
//
// Si la DB no es accesible, los tests de integracion se SKIPean
// (no FAIL), asi `go test ./...` sigue corriendo en entornos donde
// Postgres no esta levantado (por ejemplo en CI sin servicio de DB).
package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MustTestPool abre un pool contra TEST_DATABASE_URL. Si la variable
// no esta seteada o el ping falla, hace t.Skip (no FAIL).
func MustTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL no seteada, salteando integration tests")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("pgxpool.New: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("DB de tests inalcanzable: %v (correr `make db-test-up`)", err)
	}
	return pool
}

// TruncateAll limpia las tablas de datos pero mantiene segmentos
// (que son fixtures cargadas en la migracion).
func TruncateAll(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), "TRUNCATE gestiones, facturas, clientes CASCADE"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

// SeedDosClientes inserta dos clientes contrastantes:
//   - "Saludable PyME" sin facturas pendientes (score esperado bajo)
//   - "Zombi Test" con factura vencida hace 100 dias (score esperado alto)
//
// Devuelve sus IDs en ese orden.
func SeedDosClientes(t *testing.T, pool *pgxpool.Pool) (saludableID, zombiID string) {
	t.Helper()
	ctx := context.Background()

	if err := pool.QueryRow(ctx, `
		INSERT INTO clientes (nombre, fecha_alta, mrr_usd, payment_terms_dias, segmento)
		VALUES ('Saludable PyME', '2024-01-01', 800, 30, 'pyme_sana')
		RETURNING id::text
	`).Scan(&saludableID); err != nil {
		t.Fatalf("insert saludable: %v", err)
	}

	if err := pool.QueryRow(ctx, `
		INSERT INTO clientes (nombre, fecha_alta, mrr_usd, payment_terms_dias, segmento)
		VALUES ('Zombi Test', '2023-01-01', 400, 30, 'zombi')
		RETURNING id::text
	`).Scan(&zombiID); err != nil {
		t.Fatalf("insert zombi: %v", err)
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO facturas (cliente_id, numero, fecha_emision, fecha_vencimiento, monto_usd, estado)
		VALUES ($1, 'F-001', current_date - interval '130 days', current_date - interval '100 days', 8000, 'vencida')
	`, zombiID); err != nil {
		t.Fatalf("insert factura zombi: %v", err)
	}

	return saludableID, zombiID
}
