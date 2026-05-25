import { useEffect, useState } from "react"
import {
  Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from "@/components/ui/select"
import { useCrearGestion } from "@/hooks/useCrearGestion"
import { ApiError } from "@/api/client"
import type { ResultadoGestion, TipoGestion } from "@/api/types"

const TIPOS: { value: TipoGestion; label: string }[] = [
  { value: "llamada", label: "Llamada" },
  { value: "email", label: "Email" },
  { value: "whatsapp", label: "WhatsApp" },
  { value: "visita", label: "Visita" },
]

const RESULTADOS: { value: ResultadoGestion; label: string }[] = [
  { value: "sin_respuesta", label: "Sin respuesta" },
  { value: "promesa_pago", label: "Promesa de pago" },
  { value: "disputa", label: "Disputa" },
  { value: "pagado", label: "Pagado" },
  { value: "otro", label: "Otro" },
]

const NOTAS_MAX = 2000

export function GestionDialog({
  clienteId, open, onOpenChange,
}: { clienteId: string; open: boolean; onOpenChange: (open: boolean) => void }) {
  const [tipo, setTipo] = useState<TipoGestion>("llamada")
  const [resultado, setResultado] = useState<ResultadoGestion>("sin_respuesta")
  const [notas, setNotas] = useState("")
  const mutation = useCrearGestion(clienteId)

  useEffect(() => {
    if (open) {
      setTipo("llamada")
      setResultado("sin_respuesta")
      setNotas("")
      mutation.reset()
    }
  }, [open]) // eslint-disable-line react-hooks/exhaustive-deps

  const notasOver = notas.length > NOTAS_MAX

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (notasOver) return
    mutation.mutate({ tipo, resultado, notas }, { onSuccess: () => onOpenChange(false) })
  }

  const errorMsg = mutation.error instanceof ApiError ? mutation.error.message : mutation.error?.message

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <form onSubmit={handleSubmit} className="space-y-5">
          <DialogHeader>
            <DialogTitle>Registrar gestión</DialogTitle>
          </DialogHeader>

          <div className="space-y-2">
            <Label htmlFor="tipo">Tipo</Label>
            <Select value={tipo} onValueChange={(v) => setTipo(v as TipoGestion)}>
              <SelectTrigger id="tipo"><SelectValue /></SelectTrigger>
              <SelectContent>
                {TIPOS.map((t) => <SelectItem key={t.value} value={t.value}>{t.label}</SelectItem>)}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="resultado">Resultado</Label>
            <Select value={resultado} onValueChange={(v) => setResultado(v as ResultadoGestion)}>
              <SelectTrigger id="resultado"><SelectValue /></SelectTrigger>
              <SelectContent>
                {RESULTADOS.map((r) => <SelectItem key={r.value} value={r.value}>{r.label}</SelectItem>)}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label htmlFor="notas">Notas</Label>
              <span className={notasOver ? "text-xs text-red-600" : "text-xs text-slate-400"}>
                {notas.length}/{NOTAS_MAX}
              </span>
            </div>
            <Textarea
              id="notas"
              rows={4}
              value={notas}
              onChange={(e) => setNotas(e.target.value)}
              placeholder="¿Qué pasó en la gestión?"
            />
          </div>

          {errorMsg && (
            <p className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-800">
              {errorMsg}
            </p>
          )}

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={mutation.isPending}>
              Cancelar
            </Button>
            <Button type="submit" disabled={mutation.isPending || notasOver}>
              {mutation.isPending ? "Registrando…" : "Registrar"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
