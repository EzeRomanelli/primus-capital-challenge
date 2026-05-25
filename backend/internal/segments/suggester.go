// Package segments sugiere el segmento de un cliente según reglas explícitas.
// Función pura, testeable con tabla de casos.
package segments

import (
	"time"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/domain"
)

type FacturaInput struct {
	FechaVencimiento time.Time
}

type Input struct {
	Hoy                time.Time
	MrrUSD             float64
	PaymentTermsDias   int
	FacturasPendientes []FacturaInput
}

const (
	mrrCorporativoUmbral     = 2500.0
	corporativoTermsMin      = 60
	zombiUmbralDiasAtraso    = 90
	enRiesgoUmbralDiasAtraso = 15
)

// Suggest aplica las reglas en orden: zombi > corporativo > en_riesgo > pyme_sana.
// Zombi gana sobre corporativo: un cliente con MRR alto pero atraso severo se
// trata como zombi (el dolor inmediato pesa más que el perfil contractual).
func Suggest(in Input) string {
	atrasoMaxDias := 0
	for _, f := range in.FacturasPendientes {
		diff := int(in.Hoy.Sub(f.FechaVencimiento).Hours() / 24)
		if diff > atrasoMaxDias {
			atrasoMaxDias = diff
		}
	}

	if atrasoMaxDias >= zombiUmbralDiasAtraso {
		return domain.SegmentoZombi
	}
	if in.MrrUSD >= mrrCorporativoUmbral && in.PaymentTermsDias >= corporativoTermsMin {
		return domain.SegmentoCorporativo
	}
	if atrasoMaxDias >= enRiesgoUmbralDiasAtraso {
		return domain.SegmentoEnRiesgo
	}
	return domain.SegmentoPyMESana
}
