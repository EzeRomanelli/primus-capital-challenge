package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/domain"
)

// Vista plana del listado priorizado: cliente + tolerancia del segmento +
// agregados de facturas pendientes y última gestión.
type ClienteResumen struct {
	ID                  string     `db:"id"`
	Nombre              string     `db:"nombre"`
	Industria           *string    `db:"industria"`
	MrrUSD              float64    `db:"mrr_usd"`
	PaymentTermsDias    int        `db:"payment_terms_dias"`
	Segmento            string     `db:"segmento"`
	ToleranciaDias      int        `db:"tolerancia_dias"`
	MontoPendienteTotal float64    `db:"monto_pendiente_total"`
	FechaVencMin        *time.Time `db:"fecha_venc_min"`
	UltimaGestionFecha  *time.Time `db:"ultima_gestion_fecha"`
}

// El score no se calcula en SQL: el ranking lo hace Go (función pura, testeable).
func ListClientesResumen(ctx context.Context, pool *pgxpool.Pool) ([]ClienteResumen, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			c.id::text                                  AS id,
			c.nombre                                    AS nombre,
			c.industria                                 AS industria,
			c.mrr_usd::float8                           AS mrr_usd,
			c.payment_terms_dias                        AS payment_terms_dias,
			c.segmento                                  AS segmento,
			s.tolerancia_dias                           AS tolerancia_dias,
			COALESCE(f.monto_total, 0)::float8          AS monto_pendiente_total,
			f.fecha_venc_min                            AS fecha_venc_min,
			g.ultima_gestion_fecha                      AS ultima_gestion_fecha
		FROM clientes c
		JOIN segmentos s ON s.nombre = c.segmento
		LEFT JOIN (
			SELECT cliente_id,
			       SUM(monto_usd)        AS monto_total,
			       MIN(fecha_vencimiento) AS fecha_venc_min
			FROM facturas
			WHERE estado IN ('pendiente', 'vencida')
			GROUP BY cliente_id
		) f ON f.cliente_id = c.id
		LEFT JOIN (
			SELECT cliente_id, MAX(fecha) AS ultima_gestion_fecha
			FROM gestiones
			GROUP BY cliente_id
		) g ON g.cliente_id = c.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[ClienteResumen])
}

func GetClienteByID(ctx context.Context, pool *pgxpool.Pool, id string) (*domain.Cliente, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			id::text          AS id,
			nombre            AS nombre,
			industria         AS industria,
			fecha_alta        AS fecha_alta,
			mrr_usd::float8   AS mrr_usd,
			payment_terms_dias AS payment_terms_dias,
			segmento          AS segmento,
			created_at        AS created_at,
			updated_at        AS updated_at
		FROM clientes
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	c, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Cliente])
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func UpdateSegmento(ctx context.Context, pool *pgxpool.Pool, id, segmento string) error {
	_, err := pool.Exec(ctx, `
		UPDATE clientes
		SET segmento = $1, updated_at = now()
		WHERE id = $2
	`, segmento, id)
	return err
}
