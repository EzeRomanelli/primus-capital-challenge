package domain

import "time"

const (
	TipoGestionLlamada  = "llamada"
	TipoGestionEmail    = "email"
	TipoGestionWhatsapp = "whatsapp"
	TipoGestionVisita   = "visita"
)

const (
	ResultadoSinRespuesta = "sin_respuesta"
	ResultadoPromesaPago  = "promesa_pago"
	ResultadoDisputa      = "disputa"
	ResultadoPagado       = "pagado"
	ResultadoOtro         = "otro"
)

var TiposGestionValidos = map[string]bool{
	TipoGestionLlamada:  true,
	TipoGestionEmail:    true,
	TipoGestionWhatsapp: true,
	TipoGestionVisita:   true,
}

var ResultadosValidos = map[string]bool{
	ResultadoSinRespuesta: true,
	ResultadoPromesaPago:  true,
	ResultadoDisputa:      true,
	ResultadoPagado:       true,
	ResultadoOtro:         true,
}

type Gestion struct {
	ID        string    `json:"id"         db:"id"`
	ClienteID string    `json:"cliente_id" db:"cliente_id"`
	Fecha     time.Time `json:"fecha"      db:"fecha"`
	Tipo      string    `json:"tipo"       db:"tipo"`
	Resultado string    `json:"resultado"  db:"resultado"`
	Notas     string    `json:"notas"      db:"notas"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
