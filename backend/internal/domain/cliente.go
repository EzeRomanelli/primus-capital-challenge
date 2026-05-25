package domain

import "time"

type Cliente struct {
	ID               string    `json:"id"                 db:"id"`
	Nombre           string    `json:"nombre"             db:"nombre"`
	Industria        *string   `json:"industria,omitempty" db:"industria"`
	FechaAlta        time.Time `json:"fecha_alta"         db:"fecha_alta"`
	MrrUSD           float64   `json:"mrr_usd"            db:"mrr_usd"`
	PaymentTermsDias int       `json:"payment_terms_dias" db:"payment_terms_dias"`
	Segmento         string    `json:"segmento"           db:"segmento"`
	CreatedAt        time.Time `json:"created_at"         db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"         db:"updated_at"`
}
