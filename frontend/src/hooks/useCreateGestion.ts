import { useMutation, useQueryClient } from "@tanstack/react-query"
import { createGestion } from "@/api/clientes"
import type { CrearGestionReq } from "@/api/types"

export function useCreateGestion(clienteId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: CrearGestionReq) => createGestion(clienteId, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["cliente", clienteId] })
      qc.invalidateQueries({ queryKey: ["clientes"] })
    },
  })
}
