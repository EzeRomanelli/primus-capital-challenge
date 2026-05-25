// Package recalc recalcula el segmento de un cliente cuando hay cambios en
// sus facturas o gestiones. Lo invocan los handlers en goroutine fire-and-forget
// despues de crear/modificar.
//
// El job tiene su propio context (no usa el del request) porque corre
// despues de que el handler ya respondio al cliente HTTP. Si fallara, no
// hay nadie a quien devolverle el error - lo logueamos y seguimos.
package recalc

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/db"
	"github.com/ezeromanelli/northwind-cobranza/backend/internal/segments"
)

// RecalcularSegmento aplica las reglas del suggester y persiste el cambio
// si difiere del segmento actual.
//
// Idempotente: si falla, la proxima gestion lo va a re-disparar.
func RecalcularSegmento(pool *pgxpool.Pool, clienteID string) {
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

	nuevo := segments.Sugerir(in)
	if nuevo == cliente.Segmento {
		return // sin cambio: no escribir
	}
	if err := db.UpdateSegmento(ctx, pool, clienteID, nuevo); err != nil {
		log.Printf("recalc: UpdateSegmento(%s, %s): %v", clienteID, nuevo, err)
		return
	}
	log.Printf("recalc: cliente %s: %s -> %s", clienteID, cliente.Segmento, nuevo)
}
