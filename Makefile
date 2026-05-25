.PHONY: help up down logs db-shell migrate-up migrate-down migrate-status seed backend-run frontend-run test db-test-up

# Cargar .env si existe y exportar variables (no duplicamos config en cada target).
ifneq (,$(wildcard .env))
  include .env
  export
endif

help:
	@echo "Targets disponibles:"
	@echo "  up             - Levantar Postgres (docker compose)"
	@echo "  down           - Apagar Postgres"
	@echo "  logs           - Ver logs de Postgres"
	@echo "  db-shell       - Abrir psql en el container"
	@echo "  migrate-up     - Aplicar migraciones pendientes"
	@echo "  migrate-down   - Revertir la ultima migracion"
	@echo "  migrate-status - Ver version actual del schema"
	@echo "  seed           - Cargar datos sinteticos"
	@echo "  backend-run    - Levantar backend Go"
	@echo "  frontend-run   - Levantar frontend Vite"
	@echo "  test           - Correr tests del backend"
	@echo "  db-test-up     - (Re)crear DB de tests y aplicar migraciones"

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f postgres

db-shell:
	docker compose exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

migrate-up:
	migrate -path backend/internal/db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path backend/internal/db/migrations -database "$(DATABASE_URL)" down 1

migrate-status:
	migrate -path backend/internal/db/migrations -database "$(DATABASE_URL)" version

seed:
	cd backend && go run ./cmd/seed

backend-run:
	cd backend && go run ./cmd/server

frontend-run:
	cd frontend && npm run dev

test:
	cd backend && go test ./...

db-test-up:
	docker compose exec -T postgres psql -U $(POSTGRES_USER) -d postgres -c "DROP DATABASE IF EXISTS northwind_test;"
	docker compose exec -T postgres psql -U $(POSTGRES_USER) -d postgres -c "CREATE DATABASE northwind_test;"
	migrate -path backend/internal/db/migrations -database "$(TEST_DATABASE_URL)" up
