import { useQuery } from "@tanstack/react-query"
import { fetchClienteDetalle } from "@/api/clientes"

export function useCliente(id: string | undefined) {
  return useQuery({
    queryKey: ["cliente", id],
    queryFn: () => fetchClienteDetalle(id!),
    enabled: Boolean(id),
  })
}
