// Package segments sugiere el segmento de un cliente segun reglas explicitas.
//
// Es una funcion pura: recibe `Hoy`, no toca DB, no escribe logs. La analista
// puede sobreescribir lo sugerido desde la UI (Flujo B del MVP); este paquete
// se ocupa unicamente del lado del sistema.
//
// Las reglas viven aca y deben mantenerse sincronizadas con la spec
// (seccion 3.2). Cualquier cambio amerita actualizar tambien la descripcion
// de cada segmento en la migracion inicial.
package segments

import (
	"time"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/domain"
)

// FacturaInput es el subset que necesita el suggester. Por ahora solo la fecha
// de vencimiento; si en el futuro las reglas miran el monto, se agrega aca.
type FacturaInput struct {
	FechaVencimiento time.Time
}

type Input struct {
	Hoy                time.Time
	MrrUSD             float64
	PaymentTermsDias   int
	FacturasPendientes []FacturaInput
}

// Umbrales. Centralizados aca para que cambiarlos sea una sola edicion.
const (
	mrrCorporativoUmbral     = 2500.0
	corporativoTermsMin      = 60
	zombiUmbralDiasAtraso    = 90
	enRiesgoUmbralDiasAtraso = 15
)

// Sugerir aplica las reglas en orden de prioridad y devuelve uno de los 4
// segmentos canonicos definidos en domain.
//
// Orden importante: zombi > corporativo > en_riesgo > pyme_sana.
// Un cliente con MRR alto pero que dejo de pagar hace 100 dias sigue siendo
// zombi, no corporativo: el dolor inmediato pesa mas que el perfil contractual.
func Sugerir(in Input) string {
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
