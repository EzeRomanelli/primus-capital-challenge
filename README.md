# Northwind — Cobranza

MVP de 3 días para el challenge de Primus Capital. Herramienta web que ayuda al equipo de cobranza de Northwind a saber **a quién llamar hoy** y por qué.

## El problema

Northwind facturó USD 380K en el último mes. La mora pasó de 6% a 14% en un año. Dos analistas gestionan 420 clientes con un Sheet manual. La CEO pidió textual *"foco"*.

El MVP entrega ese foco: una **tabla priorizada por score** + la capacidad de **registrar gestiones** que persisten en la DB. Cuando una gestión cambia algo relevante (ej: "pagado"), el segmento del cliente se recalcula automáticamente.

## Stack

- **Backend:** Go 1.25 + Chi v5 + pgx v5. SQL plano, sin ORM.
- **Frontend:** React 18 + TypeScript + Vite + Tailwind 3 + shadcn/ui + TanStack Query v5 + React Router v6.
- **DB:** PostgreSQL 16.
- **Orquestación:** Docker Compose (todo en containers).

## Quick start

**Único requisito:** Docker (con `docker compose`). Nada más en el host.

```bash
cp .env.example .env
docker compose up
```

Eso levanta 5 servicios en orden: Postgres → migraciones → seed (420 clientes sintéticos, semilla fija) → backend → frontend. Cuando termina el build, abrí http://localhost:5173.

Para apagar: `docker compose down` (preserva DB) o `docker compose down -v` (resetea desde cero).

### Cómo probar el flujo principal

1. **Abrí http://localhost:5173.** Ves la tabla priorizada: los clientes con score más alto arriba (urgencia × 0.6 + impacto × 0.4). Los zombis suelen estar primero.
2. **Click en cualquier fila.** Ves el detalle del cliente: score con desglose, facturas pendientes, gestiones.
3. **Click "+ Registrar gestión".** Elegís tipo (llamada/email/whatsapp/visita), resultado (sin_respuesta/promesa_pago/disputa/pagado/otro), notas opcionales. Submit.
4. **La gestión queda persistida.** El segmento del cliente se recalcula en background (ej: si registrás `pagado`, un cliente zombi puede dejar de ser zombi).
5. **Volvé al dashboard.** Recargá. El orden y el badge pueden haber cambiado.

### Modo desarrollo (HMR)

Para iterar código local sin Docker para backend/frontend:

```bash
cp .env.example .env
docker compose up postgres -d
make migrate-up
make seed
make backend-run                                # terminal 1
cd frontend && npm install && make frontend-run # terminal 2
```

Pre-requisitos extra para modo dev: Go 1.25+, Node 22+, `golang-migrate` (`brew install golang-migrate`).

## Tests

```bash
make test           # unit + integration, ~2s
make db-test-up     # (re)crear DB de tests
```

3 niveles: `scoring.Calcular` (tabla de casos), `segments.Sugerir` (tabla de casos), `GET /api/clientes` contra Postgres real.

## Documentación

- [`DECISIONS.md`](DECISIONS.md) — las decisiones más importantes que tomé y por qué.
- [`AI_LOG.md`](AI_LOG.md) — cómo usé Claude Code.
- [`docs/API.md`](docs/API.md) — contrato HTTP de los 5 endpoints.
- **Swagger UI:** http://localhost:8080/swagger/ (cuando el backend está corriendo). Spec en YAML: http://localhost:8080/openapi.yaml.
- [`docs/diagrams/data-model.png`](docs/diagrams/data-model.png) — diagrama ER del schema.
