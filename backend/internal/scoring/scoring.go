// Package scoring calcula el score 0-100 de priorización de un cliente.
// Función pura, testeable con tabla de casos.
package scoring

import (
	"math"
	"time"
)

const (
	pesoUrgencia = 0.6
	pesoImpacto  = 0.4

	// 90 días de atraso neto (post-tolerancia) = urgencia 100; más se capea.
	ventanaUrgenciaDias = 90
)

type FacturaInput struct {
	FechaVencimiento time.Time
	MontoUSD         float64
}

type Input struct {
	Hoy                time.Time
	ToleranciaDias     int
	TicketMaxUSD       float64
	FacturasPendientes []FacturaInput
}

type Resultado struct {
	Score               int
	Urgencia            int
	Impacto             int
	DiasAtrasoMax       int
	MontoPendienteTotal float64
}

func Calculate(in Input) Resultado {
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
