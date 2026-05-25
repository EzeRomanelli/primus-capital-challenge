import { useQuery } from "@tanstack/react-query"
import { fetchClientes } from "@/api/clientes"

export function useClientes() {
  return useQuery({ queryKey: ["clientes"], queryFn: fetchClientes })
}
