import { useMemo, useState } from "react"
import { useNavigate } from "react-router-dom"
import { useClientes } from "@/hooks/useClientes"
import { SegmentoBadge } from "@/components/SegmentoBadge"
import { FiltroSegmentos, type SegmentoFiltro } from "@/components/FiltroSegmentos"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from "@/components/ui/table"
import { cn } from "@/lib/utils"
import type { ClientePriorizadoDTO, SegmentoNombre } from "@/api/types"

const fmtUSD = new Intl.NumberFormat("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })
const fmtFecha = new Intl.DateTimeFormat("es-AR", { day: "2-digit", month: "2-digit" })

function scoreColor(score: number) {
  if (score >= 70) return "text-red-700"
  if (score >= 40) return "text-amber-700"
  return "text-slate-600"
}

function formatGestion(iso?: string) {
  if (!iso) return "—"
  return fmtFecha.format(new Date(iso))
}

function buildCounts(cs: ClientePriorizadoDTO[]): Record<SegmentoFiltro, number> {
  const acc: Record<SegmentoFiltro, number> = {
    todos: cs.length,
    corporativo: 0,
    pyme_sana: 0,
    en_riesgo: 0,
    zombi: 0,
  }
  for (const c of cs) acc[c.segmento as SegmentoNombre]++
  return acc
}

export function Dashboard() {
  const navigate = useNavigate()
  const { data, isPending, isError } = useClientes()
  const clientes = data ?? []

  const [seleccionado, setSeleccionado] = useState<SegmentoFiltro>("todos")
  const counts = useMemo(() => buildCounts(clientes), [clientes])
  const filtrados = useMemo(
    () => (seleccionado === "todos" ? clientes : clientes.filter((c) => c.segmento === seleccionado)),
    [clientes, seleccionado],
  )

  return (
    <main className="min-h-screen bg-slate-50 p-8">
      <div className="mx-auto max-w-7xl space-y-6">
        <header className="space-y-1">
          <h1 className="text-2xl font-semibold text-slate-900">Northwind — Cobranza</h1>
          <p className="text-sm text-slate-600">
            {isPending
              ? "Cargando cartera…"
              : `Tu día, priorizado por score. ${clientes.length} clientes activos.`}
          </p>
        </header>

        {!isPending && !isError && (
          <FiltroSegmentos seleccionado={seleccionado} onSelect={setSeleccionado} counts={counts} />
        )}

        {isPending && (
          <div className="space-y-2 rounded-md border bg-white p-4">
            {Array.from({ length: 8 }).map((_, i) => <Skeleton key={i} className="h-12 w-full" />)}
          </div>
        )}

        {isError && (
          <div className="rounded-md border border-red-200 bg-red-50 p-4 text-sm text-red-800">
            No pude cargar los clientes. ¿El backend está corriendo en{" "}
            <code className="font-mono">{import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080"}</code>?
          </div>
        )}

        {!isPending && !isError && (
          <>
            <div className="overflow-hidden rounded-md border bg-white">
              <Table>
                <TableHeader>
                  <TableRow className="bg-slate-50 hover:bg-slate-50">
                    <TableHead className="text-[11px] uppercase tracking-wider text-slate-500">Cliente</TableHead>
                    <TableHead className="text-[11px] uppercase tracking-wider text-slate-500">Segmento</TableHead>
                    <TableHead className="text-right text-[11px] uppercase tracking-wider text-slate-500">Score</TableHead>
                    <TableHead className="text-right text-[11px] uppercase tracking-wider text-slate-500">Monto pend.</TableHead>
                    <TableHead className="text-right text-[11px] uppercase tracking-wider text-slate-500">Atraso</TableHead>
                    <TableHead className="text-right text-[11px] uppercase tracking-wider text-slate-500">Ult. gestión</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filtrados.length === 0 && (
                    <TableRow><TableCell colSpan={6} className="py-12 text-center text-sm text-slate-500">
                      No hay clientes en este segmento.
                    </TableCell></TableRow>
                  )}
                  {filtrados.map((c) => (
                    <TableRow
                      key={c.id}
                      className="cursor-pointer"
                      onClick={() => navigate(`/clientes/${c.id}`)}
                    >
                      <TableCell className="py-3">
                        <div className="font-medium text-slate-900">{c.nombre}</div>
                        {c.industria && <div className="text-xs text-slate-500">{c.industria}</div>}
                      </TableCell>
                      <TableCell><SegmentoBadge segmento={c.segmento} /></TableCell>
                      <TableCell className="text-right">
                        <span className={cn("font-mono text-base font-semibold tabular-nums", scoreColor(c.score))}>
                          {c.score}
                        </span>
                      </TableCell>
                      <TableCell className="text-right font-mono text-sm tabular-nums text-slate-900">
                        {fmtUSD.format(c.monto_pendiente_total)}
                      </TableCell>
                      <TableCell className="text-right font-mono text-sm tabular-nums text-slate-600">
                        {c.dias_atraso_max}d
                      </TableCell>
                      <TableCell className="text-right text-sm text-slate-500">
                        {formatGestion(c.ultima_gestion_fecha)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
            <footer className="flex items-center justify-between text-xs text-slate-400">
              <span>Mostrando {filtrados.length} de {clientes.length}.</span>
              <span className="font-mono">score = urgencia × 0.6 + impacto × 0.4</span>
            </footer>
          </>
        )}
      </div>
    </main>
  )
}
