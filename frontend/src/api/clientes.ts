import { apiFetch } from "./client"
import type { ClienteDetalleDTO, ClientePriorizadoDTO, CrearGestionReq, Gestion } from "./types"

export function fetchClientes(): Promise<ClientePriorizadoDTO[]> {
  return apiFetch<ClientePriorizadoDTO[]>("/api/clientes")
}

export function fetchClienteDetalle(id: string): Promise<ClienteDetalleDTO> {
  return apiFetch<ClienteDetalleDTO>(`/api/clientes/${id}`)
}

export function createGestion(id: string, body: CrearGestionReq): Promise<Gestion> {
  return apiFetch<Gestion>(`/api/clientes/${id}/gestiones`, { method: "POST", body })
}
