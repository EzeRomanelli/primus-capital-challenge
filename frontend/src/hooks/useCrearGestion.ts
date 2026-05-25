import { useMutation, useQueryClient } from "@tanstack/react-query"
import { crearGestion } from "@/api/clientes"
import type { CrearGestionReq } from "@/api/types"

export function useCrearGestion(clienteId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: CrearGestionReq) => crearGestion(clienteId, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["cliente", clienteId] })
      qc.invalidateQueries({ queryKey: ["clientes"] })
    },
  })
}
