# API — Northwind Cobranza

Contrato HTTP. Base URL en local: `http://localhost:8080`.

> **Tip:** este mismo contrato está disponible como spec **OpenAPI 3.0** embebido en el binario:
> - **Swagger UI interactiva:** http://localhost:8080/swagger/ (probar endpoints desde el browser)
> - **Spec YAML:** http://localhost:8080/openapi.yaml (machine-readable, importable a Postman/Insomnia)
>
> Fuente del spec: `backend/internal/api/openapi.yaml` (single source of truth, embebido vía `//go:embed`).

## Convenciones

- **Format:** JSON. Errores 4xx/5xx devuelven `{"error": string, "code": string}`.
- **Fechas:** ISO 8601 UTC en wire. Frontend formatea a `es-AR`.
- **Auth:** ninguna (single-user asumido).

| HTTP | Cuándo |
|---|---|
| 200 | Lectura OK |
| 201 | Recurso creado |
| 400 | Body mal formado o UUID inválido |
| 404 | No encontrado |
| 422 | Validación semántica (enum fuera de rango, longitud inválida) |
| 500 | Error de DB o interno |

---

## `GET /health`

```json
{ "status": "ok" }
```

## `GET /api/segmentos`

Lista los 4 segmentos canónicos con su tolerancia. Cambian raramente; frontend cachea 10 min.

```json
[
  {"nombre":"zombi",      "tolerancia_dias":0,  "descripcion":"..."},
  {"nombre":"en_riesgo",  "tolerancia_dias":5,  "descripcion":"..."},
  {"nombre":"pyme_sana",  "tolerancia_dias":15, "descripcion":"..."},
  {"nombre":"corporativo","tolerancia_dias":30, "descripcion":"..."}
]
```

## `GET /api/clientes`

Listado priorizado por score descendente. Es la fuente del Dashboard.

**Response 200** — `ClientePriorizadoDTO[]`
```json
[{
  "id": "uuid",
  "nombre": "...",
  "industria": "Logistica",
  "segmento": "zombi",
  "mrr_usd": 480,
  "payment_terms_dias": 30,
  "monto_pendiente_total": 8700,
  "dias_atraso_max": 95,
  "ultima_gestion_fecha": null,
  "score": 83,
  "urgencia": 100,
  "impacto": 58
}]
```

## `GET /api/clientes/{id}`

Detalle: cliente + score con desglose + facturas + gestiones.

**Path params:** `id` UUID v4-ish (regex `^[0-9a-f]{8}-...$`); 400 si no matchea.

**Response 200** — `ClienteDetalleDTO`
```json
{
  "cliente": { "id": "...", "nombre": "...", "segmento": "en_riesgo", "...": "..." },
  "score": {
    "Score": 57, "Urgencia": 51, "Impacto": 65,
    "DiasAtrasoMax": 51, "MontoPendienteTotal": 9800
  },
  "facturas": [ /* Factura[] */ ],
  "gestiones": [ /* Gestion[] */ ]
}
```

> **PascalCase en `score`:** el struct `scoring.Resultado` en Go no tiene json tags, así que Go serializa el nombre exacto del campo exportado. Frontend lo tipa así.

**Errores:** `404 not_found` si el cliente no existe.

## `POST /api/clientes/{id}/gestiones`

Registra una nueva gestión. Side effect fire-and-forget: el segmento del cliente puede recalcularse después del response.

**Body**
```json
{ "tipo": "llamada", "resultado": "promesa_pago", "notas": "Habló con CFO..." }
```

**Validaciones:**

| Campo | Valores válidos |
|---|---|
| `tipo` | `llamada \| email \| whatsapp \| visita` |
| `resultado` | `sin_respuesta \| promesa_pago \| disputa \| pagado \| otro` |
| `notas` | Hasta 2000 chars (puede ser `""`) |

**Response 201** — `Gestion` completa con `id`, `fecha`, `created_at`.

**Errores típicos:** `422 invalid_tipo`, `422 invalid_resultado`, `422 notas_too_long`.

---

## Cosas que NO están por decisión

| Falta | Por qué |
|---|---|
| Paginación de `/api/clientes` | 420 clientes caben en 1 request |
| `PATCH /api/clientes/{id}/segmento` (override) | Cortado del MVP (ver DECISIONS #4) |
| Filtrado server-side por segmento | Cliente filtra in-memory si quiere |
| Auth / multi-user | Single-user asumido |
