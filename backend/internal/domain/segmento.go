package domain

const (
	SegmentoCorporativo = "corporativo"
	SegmentoPyMESana    = "pyme_sana"
	SegmentoEnRiesgo    = "en_riesgo"
	SegmentoZombi       = "zombi"
)

var SegmentosValidos = map[string]bool{
	SegmentoCorporativo: true,
	SegmentoPyMESana:    true,
	SegmentoEnRiesgo:    true,
	SegmentoZombi:       true,
}

type Segmento struct {
	Nombre         string `json:"nombre"          db:"nombre"`
	ToleranciaDias int    `json:"tolerancia_dias" db:"tolerancia_dias"`
	Descripcion    string `json:"descripcion"     db:"descripcion"`
}
