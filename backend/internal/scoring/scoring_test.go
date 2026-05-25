package scoring

import (
	"testing"
	"time"
)

// parseDate es un helper para tests; YYYY-MM-DD sin TZ (medianoche UTC).
func parseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

// Casos representativos elegidos para cubrir:
//   1) cliente al dia (camino vacio)
//   2) atraso dentro de la tolerancia del segmento -> urgencia 0
//   3) zombi clasico al tope del scoring
//   4) capeo de impacto cuando el monto excede el ticket maximo
//   5) multiples facturas: urgencia toma el max(atraso), impacto suma montos
func TestCalculate(t *testing.T) {
	hoy := parseDate("2026-05-23")
	const ticketMax = 15000.0

	tests := []struct {
		name      string
		input     Input
		wantScore int
		wantUrg   int
		wantImp   int
	}{
		{
			name: "cliente sin facturas pendientes -> score 0",
			input: Input{
				Hoy:                hoy,
				ToleranciaDias:     15,
				TicketMaxUSD:       ticketMax,
				FacturasPendientes: nil,
			},
			wantScore: 0,
			wantUrg:   0,
			wantImp:   0,
		},
		{
			name: "corporativo dentro de ventana (atraso 20, tolerancia 30) -> urgencia 0",
			input: Input{
				Hoy:            hoy,
				ToleranciaDias: 30,
				TicketMaxUSD:   ticketMax,
				FacturasPendientes: []FacturaInput{
					{FechaVencimiento: hoy.AddDate(0, 0, -20), MontoUSD: 5000},
				},
			},
			// urg=0, imp=round(5000/15000*100)=33, score=round(0*0.6 + 33*0.4)=13
			wantScore: 13,
			wantUrg:   0,
			wantImp:   33,
		},
		{
			name: "zombi 95d atraso + USD 8.7K -> score ~83",
			input: Input{
				Hoy:            hoy,
				ToleranciaDias: 0,
				TicketMaxUSD:   ticketMax,
				FacturasPendientes: []FacturaInput{
					{FechaVencimiento: hoy.AddDate(0, 0, -95), MontoUSD: 8700},
				},
			},
			// urg=round(min(95/90,1)*100)=100, imp=round(8700/15000*100)=58, score=round(60+23.2)=83
			wantScore: 83,
			wantUrg:   100,
			wantImp:   58,
		},
		{
			name: "monto excedido se capea a 100",
			input: Input{
				Hoy:            hoy,
				ToleranciaDias: 15,
				TicketMaxUSD:   ticketMax,
				FacturasPendientes: []FacturaInput{
					{FechaVencimiento: hoy.AddDate(0, 0, -20), MontoUSD: 20000},
				},
			},
			// urg=round(5/90*100)=6, imp=round(min(20000/15000,1)*100)=100, score=round(3.6+40)=44
			wantScore: 44,
			wantUrg:   6,
			wantImp:   100,
		},
		{
			name: "multi facturas: urgencia usa max(atraso), impacto suma montos",
			input: Input{
				Hoy:            hoy,
				ToleranciaDias: 15,
				TicketMaxUSD:   ticketMax,
				FacturasPendientes: []FacturaInput{
					{FechaVencimiento: hoy.AddDate(0, 0, -30), MontoUSD: 2000},
					{FechaVencimiento: hoy.AddDate(0, 0, -45), MontoUSD: 3000},
				},
			},
			// max atraso = 45, neto=30. urg=round(30/90*100)=33. imp=round(5000/15000*100)=33. score=33.
			wantScore: 33,
			wantUrg:   33,
			wantImp:   33,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Calculate(tt.input)
			if got.Score != tt.wantScore {
				t.Errorf("Score: got %d, want %d", got.Score, tt.wantScore)
			}
			if got.Urgencia != tt.wantUrg {
				t.Errorf("Urgencia: got %d, want %d", got.Urgencia, tt.wantUrg)
			}
			if got.Impacto != tt.wantImp {
				t.Errorf("Impacto: got %d, want %d", got.Impacto, tt.wantImp)
			}
		})
	}
}
