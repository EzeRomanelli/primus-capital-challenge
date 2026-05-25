// Package recalc actualiza el segmento de un cliente tras cambios en facturas
// o gestiones. Se dispara fire-and-forget desde los handlers, con context propio
// (el del request ya terminó). Si falla, se loguea y la próxima gestión lo
// vuelve a disparar.
package recalc

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/db"
	"github.com/ezeromanelli/northwind-cobranza/backend/internal/segments"
)

func RecalculateSegment(pool *pgxpool.Pool, clienteID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cliente, err := db.GetClienteByID(ctx, pool, clienteID)
	if err != nil {
		log.Printf("recalc: GetClienteByID(%s): %v", clienteID, err)
		return
	}

	pendientes, err := db.ListFacturasPendientesByCliente(ctx, pool, clienteID)
	if err != nil {
		log.Printf("recalc: ListFacturasPendientesByCliente(%s): %v", clienteID, err)
		return
	}

	in := segments.Input{
		Hoy:              time.Now().UTC(),
		MrrUSD:           cliente.MrrUSD,
		PaymentTermsDias: cliente.PaymentTermsDias,
	}
	for _, p := range pendientes {
		in.FacturasPendientes = append(in.FacturasPendientes, segments.FacturaInput{
			FechaVencimiento: p.FechaVencimiento,
		})
	}

	nuevo := segments.Suggest(in)
	if nuevo == cliente.Segmento {
		return
	}
	if err := db.UpdateSegmento(ctx, pool, clienteID, nuevo); err != nil {
		log.Printf("recalc: UpdateSegmento(%s, %s): %v", clienteID, nuevo, err)
		return
	}
	log.Printf("recalc: cliente %s: %s -> %s", clienteID, cliente.Segmento, nuevo)
}
