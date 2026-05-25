package segments

import (
	"testing"
	"time"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/domain"
)

func parseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

// Cubrimos los 4 segmentos + un edge case importante: zombi
// gana sobre corporativo cuando hay atraso severo (orden de las reglas).
func TestSugerir(t *testing.T) {
	hoy := parseDate("2026-05-23")

	tests := []struct {
		name  string
		input Input
		want  string
	}{
		{
			name: "default: cliente sin facturas vencidas y sin features especiales -> pyme_sana",
			input: Input{
				Hoy:                hoy,
				MrrUSD:             500,
				PaymentTermsDias:   30,
				FacturasPendientes: nil,
			},
			want: domain.SegmentoPyMESana,
		},
		{
			name: "MRR alto + payment_terms >= 60 + atraso leve -> corporativo",
			input: Input{
				Hoy:              hoy,
				MrrUSD:           4000,
				PaymentTermsDias: 75,
				FacturasPendientes: []FacturaInput{
					{FechaVencimiento: hoy.AddDate(0, 0, -10)},
				},
			},
			want: domain.SegmentoCorporativo,
		},
		{
			name: "atraso 95d -> zombi (incluso sin features de corporativo)",
			input: Input{
				Hoy:              hoy,
				MrrUSD:           400,
				PaymentTermsDias: 30,
				FacturasPendientes: []FacturaInput{
					{FechaVencimiento: hoy.AddDate(0, 0, -95)},
				},
			},
			want: domain.SegmentoZombi,
		},
		{
			name: "atraso 45d (entre 15 y 89) y sin features de corporativo -> en_riesgo",
			input: Input{
				Hoy:              hoy,
				MrrUSD:           600,
				PaymentTermsDias: 30,
				FacturasPendientes: []FacturaInput{
					{FechaVencimiento: hoy.AddDate(0, 0, -45)},
				},
			},
			want: domain.SegmentoEnRiesgo,
		},
		{
			name: "edge: MRR alto + payment_terms 60 PERO atraso 100d -> zombi (zombi gana sobre corporativo)",
			input: Input{
				Hoy:              hoy,
				MrrUSD:           5000,
				PaymentTermsDias: 75,
				FacturasPendientes: []FacturaInput{
					{FechaVencimiento: hoy.AddDate(0, 0, -100)},
				},
			},
			want: domain.SegmentoZombi,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sugerir(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
