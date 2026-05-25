package domain

import "time"

// Estados de factura. Pagada implica FechaPago != nil.
// Pendiente y Vencida son ambos "debe plata"; la diferencia es semantica
// (pendiente = aun no vencio; vencida = paso fecha_vencimiento).
const (
	EstadoFacturaPendiente = "pendiente"
	EstadoFacturaPagada    = "pagada"
	EstadoFacturaVencida   = "vencida"
)

type Factura struct {
	ID               string     `json:"id"                   db:"id"`
	ClienteID        string     `json:"cliente_id"           db:"cliente_id"`
	Numero           string     `json:"numero"               db:"numero"`
	FechaEmision     time.Time  `json:"fecha_emision"        db:"fecha_emision"`
	FechaVencimiento time.Time  `json:"fecha_vencimiento"    db:"fecha_vencimiento"`
	FechaPago        *time.Time `json:"fecha_pago,omitempty" db:"fecha_pago"`
	MontoUSD         float64    `json:"monto_usd"            db:"monto_usd"`
	Estado           string     `json:"estado"               db:"estado"`
	CreatedAt        time.Time  `json:"created_at"           db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"           db:"updated_at"`
}
