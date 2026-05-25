package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/db"
	"github.com/ezeromanelli/northwind-cobranza/backend/internal/domain"
	"github.com/ezeromanelli/northwind-cobranza/backend/internal/recalc"
	"github.com/ezeromanelli/northwind-cobranza/backend/internal/scoring"
)

// ticketMaxUSDForScoring es la referencia del MVP para normalizar impacto.
// Sale del enunciado: tickets de USD 200 a USD 15.000.
const ticketMaxUSDForScoring = 15000.0

// ClientePriorizadoDTO es la fila que devuelve GET /api/clientes.
// Incluye los datos del cliente + agregados + el score y su desglose.
type ClientePriorizadoDTO struct {
	ID                  string     `json:"id"`
	Nombre              string     `json:"nombre"`
	Industria           *string    `json:"industria,omitempty"`
	Segmento            string     `json:"segmento"`
	MrrUSD              float64    `json:"mrr_usd"`
	PaymentTermsDias    int        `json:"payment_terms_dias"`
	MontoPendienteTotal float64    `json:"monto_pendiente_total"`
	DiasAtrasoMax       int        `json:"dias_atraso_max"`
	UltimaGestionFecha  *time.Time `json:"ultima_gestion_fecha,omitempty"`
	Score               int        `json:"score"`
	Urgencia            int        `json:"urgencia"`
	Impacto             int        `json:"impacto"`
}

// ClienteDetalleDTO es lo que devuelve GET /api/clientes/{id}: cliente
// completo + score con desglose + facturas + gestiones.
type ClienteDetalleDTO struct {
	Cliente   domain.Cliente    `json:"cliente"`
	Score     scoring.Resultado `json:"score"`
	Facturas  []domain.Factura  `json:"facturas"`
	Gestiones []domain.Gestion  `json:"gestiones"`
}

func (rt *Router) listClientesHandler(w http.ResponseWriter, r *http.Request) {
	resumenes, err := db.ListClientesResumen(r.Context(), rt.pool)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "db_error", "error al listar clientes")
		return
	}

	hoy := time.Now().UTC()
	out := make([]ClientePriorizadoDTO, 0, len(resumenes))
	for _, c := range resumenes {
		// Truco: una "factura virtual" con monto=total y fechaVenc=min(venc)
		// produce el MISMO score que iterar todas (urgencia usa max atraso,
		// impacto usa sum monto). Asi evitamos otra query por cliente.
		var facs []scoring.FacturaInput
		if c.FechaVencMin != nil && c.MontoPendienteTotal > 0 {
			facs = []scoring.FacturaInput{{
				FechaVencimiento: *c.FechaVencMin,
				MontoUSD:         c.MontoPendienteTotal,
			}}
		}
		res := scoring.Calcular(scoring.Input{
			Hoy:                hoy,
			ToleranciaDias:     c.ToleranciaDias,
			TicketMaxUSD:       ticketMaxUSDForScoring,
			FacturasPendientes: facs,
		})
		out = append(out, ClientePriorizadoDTO{
			ID:                  c.ID,
			Nombre:              c.Nombre,
			Industria:           c.Industria,
			Segmento:            c.Segmento,
			MrrUSD:              c.MrrUSD,
			PaymentTermsDias:    c.PaymentTermsDias,
			MontoPendienteTotal: c.MontoPendienteTotal,
			DiasAtrasoMax:       res.DiasAtrasoMax,
			UltimaGestionFecha:  c.UltimaGestionFecha,
			Score:               res.Score,
			Urgencia:            res.Urgencia,
			Impacto:             res.Impacto,
		})
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	writeJSON(w, http.StatusOK, out)
}

func (rt *Router) getClienteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !isValidUUID(id) {
		writeError(w, http.StatusBadRequest, "invalid_id", "id invalido")
		return
	}

	cliente, err := db.GetClienteByID(r.Context(), rt.pool, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "not_found", "cliente no encontrado")
			return
		}
		writeError(w, http.StatusInternalServerError, "db_error", "error obteniendo cliente")
		return
	}

	facturas, err := db.ListFacturasByCliente(r.Context(), rt.pool, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "db_error", "error listando facturas")
		return
	}
	gestiones, err := db.ListGestionesByCliente(r.Context(), rt.pool, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "db_error", "error listando gestiones")
		return
	}

	// Resolver tolerancia del segmento del cliente para el scoring.
	segs, err := db.ListSegmentos(r.Context(), rt.pool)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "db_error", "error obteniendo segmentos")
		return
	}
	tolerancia := 0
	for _, s := range segs {
		if s.Nombre == cliente.Segmento {
			tolerancia = s.ToleranciaDias
			break
		}
	}

	scoringIn := scoring.Input{
		Hoy:            time.Now().UTC(),
		ToleranciaDias: tolerancia,
		TicketMaxUSD:   ticketMaxUSDForScoring,
	}
	for _, f := range facturas {
		if f.Estado == domain.EstadoFacturaPendiente || f.Estado == domain.EstadoFacturaVencida {
			scoringIn.FacturasPendientes = append(scoringIn.FacturasPendientes, scoring.FacturaInput{
				FechaVencimiento: f.FechaVencimiento,
				MontoUSD:         f.MontoUSD,
			})
		}
	}
	res := scoring.Calcular(scoringIn)

	writeJSON(w, http.StatusOK, ClienteDetalleDTO{
		Cliente:   *cliente,
		Score:     res,
		Facturas:  facturas,
		Gestiones: gestiones,
	})
}

// CrearGestionReq es el body del POST /api/clientes/{id}/gestiones.
type CrearGestionReq struct {
	Tipo      string `json:"tipo"`
	Resultado string `json:"resultado"`
	Notas     string `json:"notas"`
}

func (rt *Router) crearGestionHandler(w http.ResponseWriter, r *http.Request) {
	clienteID := chi.URLParam(r, "id")
	if !isValidUUID(clienteID) {
		writeError(w, http.StatusBadRequest, "invalid_id", "id invalido")
		return
	}

	var req CrearGestionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "json invalido")
		return
	}
	if !domain.TiposGestionValidos[req.Tipo] {
		writeError(w, http.StatusUnprocessableEntity, "invalid_tipo", "tipo de gestion invalido")
		return
	}
	if !domain.ResultadosValidos[req.Resultado] {
		writeError(w, http.StatusUnprocessableEntity, "invalid_resultado", "resultado invalido")
		return
	}
	if len(req.Notas) > 2000 {
		writeError(w, http.StatusUnprocessableEntity, "notas_too_long", "notas excede 2000 caracteres")
		return
	}

	created, err := db.CreateGestion(r.Context(), rt.pool, clienteID, req.Tipo, req.Resultado, req.Notas)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "db_error", "no se pudo crear la gestion")
		return
	}

	// Fire-and-forget: el segmento puede cambiar por esta gestion
	// (ej: si el resultado es "pagado", el cliente puede dejar de ser zombi).
	go recalc.RecalcularSegmento(rt.pool, clienteID)

	writeJSON(w, http.StatusCreated, created)
}
