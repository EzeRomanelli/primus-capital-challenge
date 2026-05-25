package domain

// Nombres canonicos de segmentos. Cualquier valor fuera de este conjunto
// es un bug; los handlers validan contra SegmentosValidos.
const (
	SegmentoCorporativo = "corporativo"
	SegmentoPyMESana    = "pyme_sana"
	SegmentoEnRiesgo    = "en_riesgo"
	SegmentoZombi       = "zombi"
)

// SegmentosValidos es el set para validacion en handlers.
// Usamos map[string]bool en vez de slice para lookup O(1).
var SegmentosValidos = map[string]bool{
	SegmentoCorporativo: true,
	SegmentoPyMESana:    true,
	SegmentoEnRiesgo:    true,
	SegmentoZombi:       true,
}

// Segmento es una fila de la tabla de configuracion.
// 4 filas fijas insertadas en la migracion inicial.
type Segmento struct {
	Nombre         string `json:"nombre"          db:"nombre"`
	ToleranciaDias int    `json:"tolerancia_dias" db:"tolerancia_dias"`
	Descripcion    string `json:"descripcion"     db:"descripcion"`
}
