import { cn } from "@/lib/utils"
import type { SegmentoNombre } from "@/api/types"

export type SegmentoFiltro = SegmentoNombre | "todos"

const ORDEN: SegmentoFiltro[] = ["todos", "corporativo", "pyme_sana", "en_riesgo", "zombi"]

const LABELS: Record<SegmentoFiltro, string> = {
  todos: "Todos",
  corporativo: "Corporativo",
  pyme_sana: "PyME sana",
  en_riesgo: "En riesgo",
  zombi: "Zombi",
}

export function FiltroSegmentos({
  seleccionado,
  onSelect,
  counts,
}: {
  seleccionado: SegmentoFiltro
  onSelect: (s: SegmentoFiltro) => void
  counts: Record<SegmentoFiltro, number>
}) {
  return (
    <div className="flex flex-wrap gap-2">
      {ORDEN.map((s) => {
        const active = seleccionado === s
        return (
          <button
            key={s}
            type="button"
            onClick={() => onSelect(s)}
            className={cn(
              "inline-flex items-center rounded-full px-3.5 py-1.5 text-xs font-medium transition-colors",
              active
                ? "bg-slate-900 text-white"
                : "border border-slate-200 bg-white text-slate-700 hover:bg-slate-50",
            )}
          >
            {LABELS[s]} ({counts[s]})
          </button>
        )
      })}
    </div>
  )
}
