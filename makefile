include .env

API_PATH = api
MIGRATE_PATH = $(API_PATH)/migrations
DB_URL = $(DB_DSN)

.PHONY: migrate migrate-up migrate-down migrate-back client

migrate:
	cd $(API_PATH) && migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_URL) up

migrate-down:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_URL) down

migrate-back:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_URL) force $(version)

client:
	cd app && pnpm run dev
