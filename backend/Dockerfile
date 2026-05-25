# syntax=docker/dockerfile:1.7
# Multi-stage. Buildea los 2 binarios (server + seed) en una sola imagen.
# El compose elige cual correr via command override.

FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./cmd/server \
 && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/seed   ./cmd/seed

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata wget
WORKDIR /app
COPY --from=builder /out/server /app/server
COPY --from=builder /out/seed   /app/seed

EXPOSE 8080
CMD ["/app/server"]
