package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/domain"
)

// Ordenados por tolerancia ASC: del más crítico (zombi=0) al más relajado.
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
