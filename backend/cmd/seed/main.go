// Command seed carga datos sinteticos en la DB con seed fijo (determinista).
//
// Estrategia:
//  1) TRUNCATE de gestiones, facturas, clientes
//  2) 420 clientes random con distribucion realista del segmento:
//     70% pyme_sana | 15% corporativo | 10% en_riesgo | 5% zombi
//  3) Cada cliente recibe 12 meses de facturas con patron de pago acorde
//     a su segmento.
//  4) ~30% de clientes reciben una gestion en los ultimos 30 dias.
//
// Reproducibilidad: rand.New(rand.NewSource(42)) en todos los choices random.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ezeromanelli/northwind-cobranza/backend/internal/domain"
)

const (
	seedFijo       = int64(42)
	totalClientes  = 420
	mesesHistorico = 12
)

// Distribucion target de segmentos.
const (
	pctZombi       = 0.05
	pctEnRiesgo    = 0.10
	pctCorporativo = 0.15
	// resto -> pyme_sana
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL requerido (correr via `make seed` que carga .env)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("pool: %v", err)
	}
	defer pool.Close()

	mustExec(ctx, pool, "TRUNCATE gestiones, facturas, clientes CASCADE")

	rnd := rand.New(rand.NewSource(seedFijo))
	hoy := time.Now().UTC().Truncate(24 * time.Hour)

	t0 := time.Now()
	for i := 0; i < totalClientes; i++ {
		segmento := elegirSegmento(rnd)
		c := generarCliente(rnd, segmento)
		id := insertarCliente(ctx, pool, c)
		insertarFacturasHistoricas(ctx, pool, rnd, id, c, hoy)
		if rnd.Float64() < 0.3 {
			insertarGestionRandom(ctx, pool, rnd, id, hoy)
		}
	}

	log.Printf("seed OK: %d clientes en %s (seed=%d)", totalClientes, time.Since(t0).Round(time.Millisecond), seedFijo)
}

type clienteRandom struct {
	nombre, industria, segmento string
	mrr                         float64
	paymentTermsDias            int
}

func elegirSegmento(rnd *rand.Rand) string {
	x := rnd.Float64()
	switch {
	case x < pctZombi:
		return domain.SegmentoZombi
	case x < pctZombi+pctEnRiesgo:
		return domain.SegmentoEnRiesgo
	case x < pctZombi+pctEnRiesgo+pctCorporativo:
		return domain.SegmentoCorporativo
	default:
		return domain.SegmentoPyMESana
	}
}

var (
	palabrasEmpresa = []string{
		"Servicios", "Consultora", "Distribuidora", "Tecnologia", "Comercial",
		"Industrial", "Logistica", "Soluciones", "Sistemas", "Productos",
		"Importadora", "Inversiones", "Grupo", "Centro", "Cadena",
		"Constructora", "Despacho", "Vinos", "Estudio", "Talleres",
	}
	lugaresChile = []string{
		"Andina", "Austral", "del Sur", "Pacifico", "del Norte",
		"Atacama", "Cordillera", "Bio Bio", "del Maule", "Patagonia",
		"Magallanes", "del Centro", "Pillar", "Antuco", "Choapa",
	}
	sufijosChilenos = []string{"SpA", "Limitada", "S.A.", "Ltda.", "EIRL"}
	industrias      = []string{
		"Tecnologia", "Retail", "Logistica", "Construccion", "Servicios",
		"Salud", "Educacion", "AgroTech", "FinTech", "Hospitalidad",
	}
)

func generarCliente(rnd *rand.Rand, segmento string) clienteRandom {
	nombre := fmt.Sprintf("%s %s %s",
		palabrasEmpresa[rnd.Intn(len(palabrasEmpresa))],
		lugaresChile[rnd.Intn(len(lugaresChile))],
		sufijosChilenos[rnd.Intn(len(sufijosChilenos))])
	industria := industrias[rnd.Intn(len(industrias))]
	mrr, terms := caracteristicasPorSegmento(rnd, segmento)
	return clienteRandom{
		nombre:           nombre,
		industria:        industria,
		segmento:         segmento,
		mrr:              mrr,
		paymentTermsDias: terms,
	}
}

func caracteristicasPorSegmento(rnd *rand.Rand, segmento string) (float64, int) {
	switch segmento {
	case domain.SegmentoCorporativo:
		return 2500 + rnd.Float64()*5000, 60 + rnd.Intn(31)
	case domain.SegmentoPyMESana:
		return 200 + rnd.Float64()*800, 30
	case domain.SegmentoEnRiesgo:
		return 300 + rnd.Float64()*700, 30
	case domain.SegmentoZombi:
		return 200 + rnd.Float64()*600, 30
	default:
		return 500, 30
	}
}

func insertarCliente(ctx context.Context, pool *pgxpool.Pool, c clienteRandom) string {
	mesesAtras := 1 + (int(c.mrr) % 36)
	fechaAlta := time.Now().UTC().AddDate(0, -mesesAtras, 0)

	var id string
	err := pool.QueryRow(ctx, `
		INSERT INTO clientes (nombre, industria, fecha_alta, mrr_usd, payment_terms_dias, segmento)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text
	`, c.nombre, c.industria, fechaAlta, c.mrr, c.paymentTermsDias, c.segmento).Scan(&id)
	if err != nil {
		log.Fatalf("insert cliente %s: %v", c.nombre, err)
	}
	return id
}

// insertarFacturasHistoricas genera 12 meses de facturas mensuales con
// patron de pago acorde al segmento. pgx.Batch para 1 round trip por cliente.
func insertarFacturasHistoricas(ctx context.Context, pool *pgxpool.Pool, rnd *rand.Rand, clienteID string, c clienteRandom, hoy time.Time) {
	batch := &pgx.Batch{}
	for i := mesesHistorico; i >= 1; i-- {
		emision := hoy.AddDate(0, -i, 0)
		venc := emision.AddDate(0, 0, c.paymentTermsDias)
		estado, fechaPago := decidirEstadoFactura(rnd, c.segmento, i, venc)
		batch.Queue(
			`INSERT INTO facturas (cliente_id, numero, fecha_emision, fecha_vencimiento, fecha_pago, monto_usd, estado)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			clienteID, fmt.Sprintf("F-%s-%02d", clienteID[:8], mesesHistorico-i+1),
			emision, venc, fechaPago, c.mrr, estado,
		)
	}
	if err := pool.SendBatch(ctx, batch).Close(); err != nil {
		log.Fatalf("batch facturas de %s: %v", clienteID, err)
	}
}

// decidirEstadoFactura: i=1 es la mas reciente, i=12 la mas vieja.
func decidirEstadoFactura(rnd *rand.Rand, segmento string, i int, venc time.Time) (estado string, fechaPago *time.Time) {
	switch segmento {
	case domain.SegmentoZombi:
		if i <= 4+rnd.Intn(3) {
			return domain.EstadoFacturaVencida, nil
		}
		p := venc.AddDate(0, 0, rnd.Intn(15))
		return domain.EstadoFacturaPagada, &p
	case domain.SegmentoEnRiesgo:
		if i <= 1+rnd.Intn(2) {
			return domain.EstadoFacturaVencida, nil
		}
		p := venc.AddDate(0, 0, rnd.Intn(10))
		return domain.EstadoFacturaPagada, &p
	case domain.SegmentoCorporativo:
		if i == 1 {
			return domain.EstadoFacturaPendiente, nil
		}
		p := venc.AddDate(0, 0, 10+rnd.Intn(16))
		return domain.EstadoFacturaPagada, &p
	case domain.SegmentoPyMESana:
		p := venc.AddDate(0, 0, -2+rnd.Intn(8))
		return domain.EstadoFacturaPagada, &p
	default:
		p := venc
		return domain.EstadoFacturaPagada, &p
	}
}

var (
	tiposGestion      = []string{domain.TipoGestionLlamada, domain.TipoGestionEmail, domain.TipoGestionWhatsapp}
	resultadosGestion = []string{
		domain.ResultadoSinRespuesta, domain.ResultadoPromesaPago,
		domain.ResultadoDisputa, domain.ResultadoOtro,
	}
	notasMuestra = []string{
		"No atendio. Reintentar manana.",
		"Confirmo que pagara en los proximos 5 dias.",
		"Disputa el monto, escalado a comercial.",
		"Email enviado, en espera de respuesta.",
		"Sin respuesta tras 3 intentos.",
		"Esta esperando una orden de compra interna.",
		"Cliente pidio extension de plazo.",
	}
)

func insertarGestionRandom(ctx context.Context, pool *pgxpool.Pool, rnd *rand.Rand, clienteID string, hoy time.Time) {
	fecha := hoy.AddDate(0, 0, -rnd.Intn(30))
	tipo := tiposGestion[rnd.Intn(len(tiposGestion))]
	resultado := resultadosGestion[rnd.Intn(len(resultadosGestion))]
	notas := notasMuestra[rnd.Intn(len(notasMuestra))]

	_, err := pool.Exec(ctx, `
		INSERT INTO gestiones (cliente_id, fecha, tipo, resultado, notas)
		VALUES ($1, $2, $3, $4, $5)
	`, clienteID, fecha, tipo, resultado, notas)
	if err != nil {
		log.Fatalf("insert gestion: %v", err)
	}
}

func mustExec(ctx context.Context, pool *pgxpool.Pool, q string) {
	if _, err := pool.Exec(ctx, q); err != nil {
		log.Fatalf("exec %q: %v", q, err)
	}
}
