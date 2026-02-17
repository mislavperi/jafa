include server/.env
export DATABASE_URL

.PHONY: db migrate seed server web dev stop

## Ensure PostgreSQL is running and the jafa database exists
db:
	@until pg_isready -h localhost -p 5432 -U postgres >/dev/null 2>&1; do \
		echo "Waiting for PostgreSQL on localhost:5432..."; sleep 1; \
	done
	@echo "PostgreSQL is ready"
	@PGPASSWORD=password psql -h localhost -p 5432 -U postgres -tc \
		"SELECT 1 FROM pg_database WHERE datname = 'jafa'" | grep -q 1 || \
		(PGPASSWORD=password createdb -h localhost -p 5432 -U postgres jafa && echo "Created database 'jafa'")

## Run database migrations
migrate: db
	DATABASE_URL="$(DATABASE_URL)" dbmate -d $(MIGRATIONS_DIR) up

## Seed the database with sample data
seed: migrate
	cd server && go run ./scripts/seed.go

## Start the Go backend (blocking)
server:
	cd server && go run ./cmd/main.go

## Install frontend deps and start dev server (blocking)
web:
	cd web && npm install && npm run dev

## Start db + backend + frontend
dev: db
	@echo "Starting backend and frontend..."
	@trap 'kill 0' INT TERM; \
		(cd server && go run ./cmd/main.go) & \
		(cd web && npm install >/dev/null 2>&1 && npm run dev) & \
		wait

## Stop backend and frontend
stop:
	@pkill -f "go run ./cmd/main.go" 2>/dev/null || true
	@echo "Stopped"
