import { useState } from "react"
import { Link, useParams } from "react-router-dom"
import { useCliente } from "@/hooks/useCliente"
import { SegmentoBadge } from "@/components/SegmentoBadge"
import { GestionDialog } from "@/components/GestionDialog"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from "@/components/ui/table"
import { cn } from "@/lib/utils"
import type {
  ClienteDetalleDTO, EstadoFactura, Factura, Gestion,
  ResultadoGestion, TipoGestion,
} from "@/api/types"

const fmtUSD = new Intl.NumberFormat("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })
const fmtFecha = new Intl.DateTimeFormat("es-AR", { day: "2-digit", month: "2-digit", year: "numeric" })
const fmtFechaHora = new Intl.DateTimeFormat("es-AR", {
  day: "2-digit", month: "2-digit", year: "numeric", hour: "2-digit", minute: "2-digit",
})

const ESTADO_STYLE: Record<EstadoFactura, string> = {
  pendiente: "bg-amber-100 text-amber-800",
  pagada:    "bg-green-100 text-green-800",
  vencida:   "bg-red-100 text-red-800",
}

const TIPO_LABEL: Record<TipoGestion, string> = {
  llamada: "Llamada", email: "Email", whatsapp: "WhatsApp", visita: "Visita",
}

const RESULTADO_LABEL: Record<ResultadoGestion, string> = {
  sin_respuesta: "Sin respuesta",
  promesa_pago: "Promesa de pago",
  disputa: "Disputa",
  pagado: "Pagado",
  otro: "Otro",
}

function scoreColor(score: number) {
  if (score >= 70) return "text-red-700"
  if (score >= 40) return "text-amber-700"
  return "text-slate-600"
}

function date(iso?: string) { return iso ? fmtFecha.format(new Date(iso)) : "—" }

function pendientes(facturas: Factura[]) {
  return facturas.filter((f) => f.estado === "pendiente" || f.estado === "vencida")
}

export function ClienteDetail() {
  const { id } = useParams<{ id: string }>()
  const { data, isPending, isError } = useCliente(id)

  return (
    <main className="min-h-screen bg-slate-50 p-8">
      <div className="mx-auto max-w-5xl space-y-5">
        <Link to="/" className="inline-block text-sm text-slate-600 hover:text-slate-900">
          ← Volver al dashboard
        </Link>

        {isPending && (
          <div className="space-y-4">
            <Skeleton className="h-32 w-full" />
            <Skeleton className="h-40 w-full" />
            <Skeleton className="h-28 w-full" />
          </div>
        )}

        {isError && (
          <div className="rounded-md border border-red-200 bg-red-50 p-4 text-sm text-red-800">
            No pude cargar el cliente.
          </div>
        )}

        {data && <Body detalle={data} />}
      </div>
    </main>
  )
}

function Body({ detalle }: { detalle: ClienteDetalleDTO }) {
  const { cliente, score, facturas, gestiones } = detalle
  const facsPendientes = pendientes(facturas)
  const [open, setOpen] = useState(false)

  return (
    <>
      <section className="space-y-4 rounded-lg border bg-white p-6">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold text-slate-900">{cliente.nombre}</h1>
            <p className="mt-1 text-sm text-slate-600">
              {cliente.industria ?? "—"} · MRR USD {fmtUSD.format(cliente.mrr_usd)} · payment_terms {cliente.payment_terms_dias}d
            </p>
          </div>
          <SegmentoBadge segmento={cliente.segmento} />
        </div>
      </section>

      <section className="rounded-lg border bg-white p-6">
        <div className="flex items-center justify-between">
          <h2 className="text-base font-semibold text-slate-900">Score</h2>
          <span className="font-mono text-xs text-slate-400">
            score = urgencia × 0.6 + impacto × 0.4
          </span>
        </div>
        <div className="mt-4 flex items-center gap-8">
          <div className={cn("font-mono text-6xl font-bold leading-none tabular-nums", scoreColor(score.Score))}>
            {score.Score}
          </div>
          <div className="flex-1 space-y-2 text-sm">
            <BreakdownRow label="Urgencia" value={score.Urgencia} detail={`${score.DiasAtrasoMax}d atraso max`} />
            <BreakdownRow label="Impacto"  value={score.Impacto}  detail={`USD ${fmtUSD.format(score.MontoPendienteTotal)} adeudados`} />
          </div>
        </div>
      </section>

      <section className="space-y-3 rounded-lg border bg-white p-6">
        <h2 className="text-base font-semibold text-slate-900">Facturas adeudadas ({facsPendientes.length})</h2>
        <FacturasTable facturas={facsPendientes} />
      </section>

      <section className="space-y-4 rounded-lg border bg-white p-6">
        <div className="flex items-center justify-between">
          <h2 className="text-base font-semibold text-slate-900">Gestiones ({gestiones.length})</h2>
          <Button size="sm" onClick={() => setOpen(true)}>+ Registrar gestión</Button>
        </div>
        <GestionesTimeline gestiones={gestiones} />
      </section>

      <GestionDialog clienteId={cliente.id} open={open} onOpenChange={setOpen} />
    </>
  )
}

function BreakdownRow({ label, value, detail }: { label: string; value: number; detail: string }) {
  return (
    <div className="flex items-center gap-3">
      <span className="w-20 text-sm font-medium text-slate-700">{label}</span>
      <span className="w-10 text-right font-mono text-base font-semibold tabular-nums text-slate-900">{value}</span>
      <span className="text-sm text-slate-500">{detail}</span>
    </div>
  )
}

function FacturasTable({ facturas }: { facturas: Factura[] }) {
  if (facturas.length === 0) {
    return (
      <div className="rounded-md border bg-slate-50 p-6 text-center text-sm text-slate-500">
        Sin facturas adeudadas.
      </div>
    )
  }
  return (
    <div className="overflow-hidden rounded-md border">
      <Table>
        <TableHeader>
          <TableRow className="bg-slate-50 hover:bg-slate-50">
            <TableHead className="text-[10px] uppercase tracking-wider text-slate-500">Número</TableHead>
            <TableHead className="text-[10px] uppercase tracking-wider text-slate-500">Emisión</TableHead>
            <TableHead className="text-[10px] uppercase tracking-wider text-slate-500">Venc.</TableHead>
            <TableHead className="text-right text-[10px] uppercase tracking-wider text-slate-500">Monto</TableHead>
            <TableHead className="text-right text-[10px] uppercase tracking-wider text-slate-500">Estado</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {facturas.map((f) => (
            <TableRow key={f.id}>
              <TableCell className="font-mono text-sm text-slate-900">{f.numero}</TableCell>
              <TableCell className="font-mono text-sm text-slate-600">{date(f.fecha_emision)}</TableCell>
              <TableCell className="font-mono text-sm text-slate-600">{date(f.fecha_vencimiento)}</TableCell>
              <TableCell className="text-right font-mono text-sm tabular-nums text-slate-900">
                {fmtUSD.format(f.monto_usd)}
              </TableCell>
              <TableCell className="text-right">
                <span className={cn("inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-medium", ESTADO_STYLE[f.estado])}>
                  {f.estado}
                </span>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}

function GestionesTimeline({ gestiones }: { gestiones: Gestion[] }) {
  if (gestiones.length === 0) {
    return (
      <div className="rounded-md bg-slate-50 px-6 py-10 text-center">
        <p className="text-sm font-medium text-slate-700">Sin gestiones registradas todavía.</p>
        <p className="mt-1 text-xs text-slate-500">
          Cuando registres una, el segmento sugerido se recalcula automáticamente.
        </p>
      </div>
    )
  }
  return (
    <ol className="space-y-3">
      {gestiones.map((g) => (
        <li key={g.id} className="rounded-md border bg-white p-4">
          <p className="text-sm font-medium text-slate-900">
            {TIPO_LABEL[g.tipo]} · {RESULTADO_LABEL[g.resultado]}
          </p>
          <p className="font-mono text-xs text-slate-500">{fmtFechaHora.format(new Date(g.fecha))}</p>
          {g.notas && <p className="mt-2 whitespace-pre-wrap text-sm text-slate-700">{g.notas}</p>}
        </li>
      ))}
    </ol>
  )
}
