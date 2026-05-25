package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/domain"
)

// ListSegmentos devuelve las 4 filas de la tabla de configuracion.
// Ordenadas por tolerancia ASC para que el frontend las muestre del
// mas critico (zombi=0) al mas relajado (corporativo=30).
func ListSegmentos(ctx context.Context, pool *pgxpool.Pool) ([]domain.Segmento, error) {
	rows, err := pool.Query(ctx, `
		SELECT nombre, tolerancia_dias, descripcion
		FROM segmentos
		ORDER BY tolerancia_dias ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Segmento])
}
