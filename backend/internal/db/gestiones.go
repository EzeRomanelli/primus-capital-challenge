package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/domain"
)

// ListGestionesByCliente devuelve las gestiones de un cliente ordenadas
// por fecha DESC (timeline en UI).
func ListGestionesByCliente(ctx context.Context, pool *pgxpool.Pool, clienteID string) ([]domain.Gestion, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			id::text         AS id,
			cliente_id::text AS cliente_id,
			fecha            AS fecha,
			tipo             AS tipo,
			resultado        AS resultado,
			notas            AS notas,
			created_at       AS created_at,
			updated_at       AS updated_at
		FROM gestiones
		WHERE cliente_id = $1
		ORDER BY fecha DESC
	`, clienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Gestion])
}

// CreateGestion inserta una gestion y devuelve la fila creada
// (con id, created_at, updated_at y fecha default si no se paso).
//
// Tomamos los campos minimos como argumentos sueltos en lugar de un struct,
// porque queremos forzar al caller a pasar solo lo que el endpoint POST
// recibe (no ID, no timestamps).
func CreateGestion(ctx context.Context, pool *pgxpool.Pool, clienteID, tipo, resultado, notas string) (*domain.Gestion, error) {
	rows, err := pool.Query(ctx, `
		INSERT INTO gestiones (cliente_id, tipo, resultado, notas)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id::text         AS id,
			cliente_id::text AS cliente_id,
			fecha            AS fecha,
			tipo             AS tipo,
			resultado        AS resultado,
			notas            AS notas,
			created_at       AS created_at,
			updated_at       AS updated_at
	`, clienteID, tipo, resultado, notas)
	if err != nil {
		return nil, err
	}
	g, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Gestion])
	if err != nil {
		return nil, err
	}
	return &g, nil
}
