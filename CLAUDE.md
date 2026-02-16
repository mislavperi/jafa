# CLAUDE.md

## Project Overview

Bamako (JAFA) is a personal finance/expense tracking application with a Go backend and Vue 3 frontend, backed by PostgreSQL.

## Tech Stack

- **Backend:** Go 1.25, Gin framework, pgx/v5 (PostgreSQL driver), SQLC (code generation from SQL)
- **Frontend:** Vue 3, TypeScript, Vite, PrimeVue, Tailwind CSS, TanStack Vue Query, Pinia
- **Database:** PostgreSQL
- **Testing:** Vitest (unit), Nightwatch (e2e)

## Project Structure

```
server/                          # Go backend
  cmd/main.go                    # Entry point
  cmd/server/bootstrap/server.go # DI and server setup
  internal/domain/               # Business logic (models, services, mappers)
  internal/infrastructure/psql/  # Database layer (queries, migrations, repositories)
  internal/server/               # HTTP layer (controllers, routing)
  utils/                         # Shared utilities
web/                             # Vue 3 frontend
  src/core/views/                # Layout components
  src/modules/                   # Feature modules (expense/)
  src/router/                    # Vue Router
  src/stores/                    # Pinia stores
```

## Build & Run

### Backend

```bash
cd server && go build ./cmd/...
cd server && go run ./cmd/main.go     # runs on :8080
```

Database connection: `postgres://postgres:password@localhost:5432/jafa`

### Frontend

```bash
cd web && npm install
cd web && npm run dev       # dev server with proxy to :8080
cd web && npm run build     # production build
cd web && npm run test:unit # vitest
cd web && npm run lint      # eslint
cd web && npm run format    # prettier
```

Vite proxies `/api/*` to `http://localhost:8080`, stripping the `/api` prefix.

## Code Conventions

### Backend (Go)

- **Layered architecture:** domain (models/services/mappers) → infrastructure (psql/repositories) → server (controllers)
- **SQLC** generates type-safe Go from SQL in `internal/infrastructure/psql/query.sql`
- **Mapper pattern:** convert between SQLC-generated DB types and domain models
- **Constructor injection:** `NewExpenseService(queries, mapper)` style
- **Controllers** return `gin.HandlerFunc`
- Migrations live in `internal/infrastructure/psql/migrations/`

### Frontend (Vue 3/TypeScript)

- **Feature-based modules:** `modules/<feature>/{api,composables,models,components,views}`
- **Composition API** with `<script setup>` syntax
- **TanStack Vue Query** for server state; Pinia for client state
- **Fetch-based** API clients in `api/` directories
- PrimeVue components are auto-imported
