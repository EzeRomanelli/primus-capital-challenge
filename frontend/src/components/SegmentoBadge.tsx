import { cn } from "@/lib/utils"
import type { SegmentoNombre } from "@/api/types"

const STYLES: Record<SegmentoNombre, { bg: string; fg: string; label: string }> = {
  corporativo: { bg: "bg-blue-100",  fg: "text-blue-800",  label: "Corporativo" },
  pyme_sana:   { bg: "bg-green-100", fg: "text-green-800", label: "PyME sana"   },
  en_riesgo:   { bg: "bg-amber-100", fg: "text-amber-800", label: "En riesgo"   },
  zombi:       { bg: "bg-red-100",   fg: "text-red-800",   label: "Zombi"       },
}

export function SegmentoBadge({ segmento }: { segmento: SegmentoNombre }) {
  const s = STYLES[segmento]
  return (
    <span className={cn(
      "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium",
      s.bg, s.fg,
    )}>
      {s.label}
    </span>
  )
}
