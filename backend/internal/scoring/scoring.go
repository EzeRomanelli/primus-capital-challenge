// Package scoring calcula el score de priorizacion de un cliente.
//
// Es una funcion pura: no toca DB, no toca tiempo del sistema (recibe `Hoy`),
// no tiene logging. Todos los inputs son explicitos, todos los outputs derivan
// ineluctablemente. Por eso es testeable con tablas y por eso vive separada de
// los handlers y del paquete db.
//
// La formula y sus pesos viven en la spec del producto (seccion 3.3). Si cambian,
// se cambia aca y se actualiza la spec.
package scoring

import (
	"math"
	"time"
)

const (
	pesoUrgencia = 0.6
	pesoImpacto  = 0.4

	// Normalizamos la urgencia a una ventana de 90 dias post-tolerancia.
	// 90 dias de atraso "neto" (despues de descontar la tolerancia del segmento)
	// = urgencia 100. Mas de 90 se capea.
	ventanaUrgenciaDias = 90
)

// FacturaInput es el subset minimo de una factura que necesita el scoring.
// Pasar solo los datos relevantes evita acoplar este paquete a domain.Factura.
type FacturaInput struct {
	FechaVencimiento time.Time
	MontoUSD         float64
}

// Input agrupa todo lo necesario para calcular el score de un cliente.
// ToleranciaDias viene del segmento actual del cliente (no del sugerido).
// TicketMaxUSD es la referencia para normalizar impacto (15000 segun spec del MVP).
type Input struct {
	Hoy                time.Time
	ToleranciaDias     int
	TicketMaxUSD       float64
	FacturasPendientes []FacturaInput
}

// Resultado expone el score y su desglose, para mostrar tooltip en UI.
// DiasAtrasoMax y MontoPendienteTotal se devuelven aunque sean derivables
// del Input, asi el caller no tiene que recalcular para presentarlos.
type Resultado struct {
	Score               int     // 0-100
	Urgencia            int     // 0-100
	Impacto             int     // 0-100
	DiasAtrasoMax       int     // max(atraso) entre las facturas pendientes
	MontoPendienteTotal float64 // suma de montos pendientes
}

// Calcular devuelve el score 0-100 de un cliente con su desglose.
//
// Diseño:
//   - Sin facturas pendientes -> Resultado{} (todo cero). Cliente al dia.
//   - urgencia depende del MAYOR atraso, no del promedio: una factura muy
//     vieja debe pesar aunque las nuevas esten al dia.
//   - impacto depende de la SUMA de montos: el riesgo agregado importa.
//   - Ambos componentes se capean a 100 antes de combinarse.
func Calcular(in Input) Resultado {
	if len(in.FacturasPendientes) == 0 {
		return Resultado{}
	}

	var (
		atrasoMaxDias int
		montoTotal    float64
	)
	for _, f := range in.FacturasPendientes {
		diff := int(in.Hoy.Sub(f.FechaVencimiento).Hours() / 24)
		if diff > atrasoMaxDias {
			atrasoMaxDias = diff
		}
		montoTotal += f.MontoUSD
	}

	// Atraso neto: lo que excede la tolerancia del segmento (0 si esta dentro).
	atrasoNeto := max(0, atrasoMaxDias-in.ToleranciaDias)

	urgencia := int(math.Round(math.Min(float64(atrasoNeto)/ventanaUrgenciaDias, 1.0) * 100))

	impacto := 0
	if in.TicketMaxUSD > 0 {
		impacto = int(math.Round(math.Min(montoTotal/in.TicketMaxUSD, 1.0) * 100))
	}

	score := int(math.Round(float64(urgencia)*pesoUrgencia + float64(impacto)*pesoImpacto))

	return Resultado{
		Score:               score,
		Urgencia:            urgencia,
		Impacto:             impacto,
		DiasAtrasoMax:       atrasoMaxDias,
		MontoPendienteTotal: montoTotal,
	}
}
