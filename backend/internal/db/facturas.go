package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/domain"
)

// ListFacturasByCliente devuelve todas las facturas (de cualquier estado)
// de un cliente, mas recientes primero. La usa el endpoint de detalle.
func ListFacturasByCliente(ctx context.Context, pool *pgxpool.Pool, clienteID string) ([]domain.Factura, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			id::text         AS id,
			cliente_id::text AS cliente_id,
			numero           AS numero,
			fecha_emision    AS fecha_emision,
			fecha_vencimiento AS fecha_vencimiento,
			fecha_pago       AS fecha_pago,
			monto_usd::float8 AS monto_usd,
			estado           AS estado,
			created_at       AS created_at,
			updated_at       AS updated_at
		FROM facturas
		WHERE cliente_id = $1
		ORDER BY fecha_vencimiento DESC
	`, clienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Factura])
}

// FacturaPendienteRow es el subset minimo que necesita el job de recalculo
// de segmento. Vive aca porque es especifico de esta query.
type FacturaPendienteRow struct {
	FechaVencimiento time.Time `db:"fecha_vencimiento"`
	MontoUSD         float64   `db:"monto_usd"`
}

// ListFacturasPendientesByCliente devuelve solo las pendientes/vencidas.
// La usa el recalc para alimentar segments.Sugerir.
func ListFacturasPendientesByCliente(ctx context.Context, pool *pgxpool.Pool, clienteID string) ([]FacturaPendienteRow, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			fecha_vencimiento AS fecha_vencimiento,
			monto_usd::float8 AS monto_usd
		FROM facturas
		WHERE cliente_id = $1 AND estado IN ('pendiente', 'vencida')
	`, clienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[FacturaPendienteRow])
}
