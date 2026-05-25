import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { SegmentoBadge } from "@/components/SegmentoBadge"
import { useSegmentos } from "@/hooks/useSegmentos"
import type { Segmento } from "@/api/types"

// De más permisivo a menos: corporativo 30 → pyme_sana 15 → en_riesgo 5 → zombi 0.
function sortByToleranciaDesc(segs: Segmento[]) {
  return [...segs].sort((a, b) => b.tolerancia_dias - a.tolerancia_dias)
}

export function AyudaDialog() {
  const { data: segmentos } = useSegmentos()
  const ordered = segmentos ? sortByToleranciaDesc(segmentos) : []

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">¿Cómo funciona?</Button>
      </DialogTrigger>
      <DialogContent className="max-w-xl">
        <DialogHeader>
          <DialogTitle>Cómo se prioriza la cartera</DialogTitle>
        </DialogHeader>

        <div className="space-y-5 pt-2">
          <section className="space-y-2">
            <h3 className="text-sm font-semibold text-slate-900">El score</h3>
            <p className="text-sm text-slate-600">
              Cada cliente tiene un score 0-100 que combina dos factores. La tabla
              está ordenada por este score (descendente) — el top es a quién deberías
              atender hoy.
            </p>
            <div className="rounded-md bg-slate-50 px-3 py-2 font-mono text-xs text-slate-700">
              score = urgencia × 0.6 + impacto × 0.4
            </div>
            <ul className="space-y-1 text-xs text-slate-600">
              <li><strong>Urgencia:</strong> días de atraso de la factura más vencida vs la tolerancia de su segmento.</li>
              <li><strong>Impacto:</strong> monto pendiente total normalizado contra el ticket máximo (USD 15.000).</li>
            </ul>
          </section>

          <section className="space-y-3">
            <h3 className="text-sm font-semibold text-slate-900">Los 4 segmentos</h3>
            <p className="text-sm text-slate-600">
              Cada cliente se clasifica automáticamente según su MRR, payment terms
              acordados y patrón de pago. Cada segmento tiene una <strong>tolerancia
              en días</strong> post-vencimiento: hasta ese plazo no pesa en la urgencia.
            </p>
            <div className="space-y-3">
              {ordered.map((s) => (
                <div key={s.nombre} className="flex items-start gap-3 rounded-md border bg-white p-3">
                  <div className="pt-0.5">
                    <SegmentoBadge segmento={s.nombre} />
                  </div>
                  <div className="flex-1 space-y-1">
                    <p className="font-mono text-xs text-slate-500">
                      tolerancia: {s.tolerancia_dias}d
                    </p>
                    <p className="text-sm text-slate-700">{s.descripcion}</p>
                  </div>
                </div>
              ))}
            </div>
          </section>

          <section className="rounded-md border border-slate-200 bg-slate-50 px-3 py-2 text-xs text-slate-600">
            <strong>Tip:</strong> cuando registrás una gestión con resultado{" "}
            <code className="font-mono">pagado</code>, el segmento del cliente se
            recalcula en background — un zombi puede dejar de ser zombi.
          </section>
        </div>
      </DialogContent>
    </Dialog>
  )
}
