.PHONY: help dev-backend dev-frontend db-up db-down migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make dev-backend    - Start backend server"
	@echo "  make dev-frontend   - Start frontend dev server"
	@echo "  make db-up          - Start PostgreSQL with Docker"
	@echo "  make db-down        - Stop PostgreSQL"
	@echo "  make migrate-up     - Run database migrations"

dev-backend:
	cd backend && go run main.go

dev-frontend:
	cd frontend && npm run dev

db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down

migrate-up:
	cd backend && psql $(DATABASE_URL) -f db/migrations/001_init.sql
