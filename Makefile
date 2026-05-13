SHELL := /bin/bash

.PHONY: help install run dev dev-backend dev-frontend prod prod-backend build fmt test test-unit lint migrate migrate-up migrate-down migrate-status migrate-create env-up env-down

GO ?= go
NPM ?= npm
WEB_DIR := apps/web
MIGRATIONS_DIR := ./migrations

APP_ENV ?= development
ENV_FILE := .env

ifeq ($(APP_ENV),production)
ENV_FILE := .env.prod
endif

help:
	@echo "Cryplio targets:"
	@echo "  make run              Start backend API server"
	@echo "  make dev              Start backend and frontend in development mode"
	@echo "  make dev-backend      Run backend with Air and APP_ENV=development"
	@echo "  make dev-frontend     Run Next.js frontend"
	@echo "  make prod             Start backend with APP_ENV=production"
	@echo "  make build            Compile backend binary to ./bin/api"
	@echo "  make test             Run Go tests"
	@echo "  make lint             Run Go and frontend lint checks"
	@echo "  make migrate-up       Apply migrations"
	@echo "  make migrate-down     Roll back the last migration"
	@echo "  make migrate-status   Show migration status"
	@echo "  make migrate-force v=7 Force migration to version v"
	@echo "  make env-up           Start Docker Compose services"
	@echo "  make env-down         Stop Docker Compose services"
	@echo "  make seed             Seed the database with initial/dummy data"

install:
	cd $(WEB_DIR) && $(NPM) install

run:
	$(GO) run ./cmd/api

dev:
	@trap 'kill $$backend_pid $$frontend_pid 2>/dev/null; wait 2>/dev/null' INT TERM; \
	$(MAKE) dev-backend & backend_pid=$$!; \
	$(MAKE) dev-frontend & frontend_pid=$$!; \
	wait -n $$backend_pid $$frontend_pid; \
	status=$$?; \
	kill $$backend_pid $$frontend_pid 2>/dev/null; \
	wait 2>/dev/null; \
	exit $$status

dev-backend:
	@command -v air >/dev/null 2>&1 || { echo "air is required: go install github.com/air-verse/air@latest"; exit 1; }
	APP_ENV=development air

dev-frontend:
	cd $(WEB_DIR) && $(NPM) run dev

prod: prod-backend

prod-backend:
	APP_ENV=production $(GO) run ./cmd/api

build:
	@mkdir -p bin
	$(GO) build -o ./bin/api ./cmd/api

fmt:
	$(GO) fmt ./...

test: test-unit

test-unit:
	$(GO) test ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is required"; exit 1; }
	golangci-lint run ./...
	cd $(WEB_DIR) && $(NPM) run lint

migrate: migrate-up

migrate-up:
	@set -a; \
	[ -f "$(ENV_FILE)" ] && source "$(ENV_FILE)"; \
	set +a; \
	APP_ENV=$(APP_ENV) $(GO) run ./cmd/migrate -create-db -dir="$(MIGRATIONS_DIR)" \
		-host="$${DB_HOST}" -port="$${DB_PORT}" -user="$${DB_USER}" -password="$${DB_PASSWORD}" -dbname="$${DB_NAME}"

migrate-down:
	@set -a; \
	[ -f "$(ENV_FILE)" ] && source "$(ENV_FILE)"; \
	set +a; \
	APP_ENV=$(APP_ENV) $(GO) run ./cmd/migrate -down -dir="$(MIGRATIONS_DIR)" \
		-host="$${DB_HOST}" -port="$${DB_PORT}" -user="$${DB_USER}" -password="$${DB_PASSWORD}" -dbname="$${DB_NAME}"

migrate-status:
	@set -a; \
	[ -f "$(ENV_FILE)" ] && source "$(ENV_FILE)"; \
	set +a; \
	APP_ENV=$(APP_ENV) $(GO) run ./cmd/migrate -status -dir="$(MIGRATIONS_DIR)" \
		-host="$${DB_HOST}" -port="$${DB_PORT}" -user="$${DB_USER}" -password="$${DB_PASSWORD}" -dbname="$${DB_NAME}"

migrate-force:
	@test -n "$(v)" || { echo "usage: make migrate-force v=7"; exit 1; }
	@set -a; \
	[ -f "$(ENV_FILE)" ] && source "$(ENV_FILE)"; \
	set +a; \
	APP_ENV=$(APP_ENV) $(GO) run ./cmd/migrate -force=$(v) -dir="$(MIGRATIONS_DIR)" \
		-host="$${DB_HOST}" -port="$${DB_PORT}" -user="$${DB_USER}" -password="$${DB_PASSWORD}" -dbname="$${DB_NAME}"

migrate-create:
	@test -n "$(name)" || { echo "usage: make migrate-create name=add_feature"; exit 1; }
	@next=$$(printf "%03d" $$(( $$(find $(MIGRATIONS_DIR) -maxdepth 1 -name '*.up.sql' | wc -l) )) ); \
	touch "$(MIGRATIONS_DIR)/$${next}_$(name).up.sql" "$(MIGRATIONS_DIR)/$${next}_$(name).down.sql"; \
	echo "created $(MIGRATIONS_DIR)/$${next}_$(name).up.sql"; \
	echo "created $(MIGRATIONS_DIR)/$${next}_$(name).down.sql"

env-up:
	docker compose up -d

env-down:
	docker compose down

seed:
	@set -a; \
	[ -f "$(ENV_FILE)" ] && source "$(ENV_FILE)"; \
	set +a; \
	$(GO) run ./cmd/seed
