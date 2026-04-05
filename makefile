include .env

include api/.env

API_PATH = api
MIGRATE_PATH = $(API_PATH)/migrations
DB_URL = $(DB_DSN)

.PHONY: migrate migrate-up migrate-down migrate-back api-test client

migrate:
	cd $(API_PATH) && migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_URL) up

migrate-down:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_URL) down

migrate-back:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_URL) force $(version)

api-test:
	cd api && go test ./... -count=1

client:
	cd app && pnpm run dev
