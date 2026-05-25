import type { ErrorResponse } from "./types"

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080"

export class ApiError extends Error {
  status: number
  code: string

  constructor(status: number, code: string, message: string) {
    super(message)
    this.name = "ApiError"
    this.status = status
    this.code = code
  }
}

type RequestOptions = Omit<RequestInit, "body"> & { body?: unknown }

export async function apiFetch<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { body, headers, ...rest } = opts
  const init: RequestInit = {
    ...rest,
    headers: {
      "Accept": "application/json",
      ...(body !== undefined ? { "Content-Type": "application/json" } : {}),
      ...headers,
    },
    body: body !== undefined ? JSON.stringify(body) : undefined,
  }
  const res = await fetch(`${BASE_URL}${path}`, init)
  if (res.status === 204) return undefined as T
  const text = await res.text()
  const data = text ? (JSON.parse(text) as unknown) : undefined
  if (!res.ok) {
    const err = (data ?? {}) as Partial<ErrorResponse>
    throw new ApiError(res.status, err.code ?? "unknown", err.error ?? `HTTP ${res.status}`)
  }
  return data as T
}
