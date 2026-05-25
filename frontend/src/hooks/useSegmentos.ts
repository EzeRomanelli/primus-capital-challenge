import { useQuery } from "@tanstack/react-query"
import { fetchSegmentos } from "@/api/segmentos"

export function useSegmentos() {
  return useQuery({
    queryKey: ["segmentos"],
    queryFn: fetchSegmentos,
    staleTime: 10 * 60 * 1000, // 10 min: config rara vez cambia
  })
}
