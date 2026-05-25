import { apiFetch } from "./client"
import type { Segmento } from "./types"

export function fetchSegmentos(): Promise<Segmento[]> {
  return apiFetch<Segmento[]>("/api/segmentos")
}
