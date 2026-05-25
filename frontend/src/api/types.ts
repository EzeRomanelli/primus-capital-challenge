export type SegmentoNombre = "corporativo" | "pyme_sana" | "en_riesgo" | "zombi"

export type TipoGestion = "llamada" | "email" | "whatsapp" | "visita"

export type ResultadoGestion =
  | "sin_respuesta"
  | "promesa_pago"
  | "disputa"
  | "pagado"
  | "otro"

export type EstadoFactura = "pendiente" | "pagada" | "vencida"

export interface Segmento {
  nombre: SegmentoNombre
  tolerancia_dias: number
  descripcion: string
}

export interface Cliente {
  id: string
  nombre: string
  industria?: string
  fecha_alta: string
  mrr_usd: number
  payment_terms_dias: number
  segmento: SegmentoNombre
  created_at: string
  updated_at: string
}

export interface Factura {
  id: string
  cliente_id: string
  numero: string
  fecha_emision: string
  fecha_vencimiento: string
  fecha_pago?: string
  monto_usd: number
  estado: EstadoFactura
  created_at: string
  updated_at: string
}

export interface Gestion {
  id: string
  cliente_id: string
  fecha: string
  tipo: TipoGestion
  resultado: ResultadoGestion
  notas: string
  created_at: string
  updated_at: string
}

export interface ClientePriorizadoDTO {
  id: string
  nombre: string
  industria?: string
  segmento: SegmentoNombre
  mrr_usd: number
  payment_terms_dias: number
  monto_pendiente_total: number
  dias_atraso_max: number
  ultima_gestion_fecha?: string
  score: number
  urgencia: number
  impacto: number
}

// PascalCase porque scoring.Resultado en Go no tiene json tags
// y el encoder serializa los nombres exactos del campo exportado.
export interface ScoreResultado {
  Score: number
  Urgencia: number
  Impacto: number
  DiasAtrasoMax: number
  MontoPendienteTotal: number
}

export interface ClienteDetalleDTO {
  cliente: Cliente
  score: ScoreResultado
  facturas: Factura[]
  gestiones: Gestion[]
}

export interface ErrorResponse {
  error: string
  code: string
}

export interface CrearGestionReq {
  tipo: TipoGestion
  resultado: ResultadoGestion
  notas: string
}
