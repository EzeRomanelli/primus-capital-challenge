// Package db expone el pool de conexiones a Postgres y las queries SQL
// usadas por el resto del backend. Sin ORM: SQL plano via pgx.
//
// Convencion de casteos en SELECT:
//   - id UUID  -> SELECT id::text AS id      (pgx representa UUID como [16]byte
//     por defecto; el cast a text permite scanear directo a Go string)
//   - monto NUMERIC -> SELECT monto::float8  (pgx puede manejar NUMERIC, pero
//     float8 evita ambiguedades de precision para el monto en USD del MVP)
package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool crea un pool, lo pingea, y lo devuelve listo para usar.
// Si el Ping falla cerramos el pool y devolvemos el error.
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return pool, nil
}
