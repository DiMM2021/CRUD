# Создание таблицы
MIGRATE_BIN=$(shell go env GOPATH)/bin/migrate
MIGRATIONS_DIR=db/migrations

create-migration:
	@read -p "Enter migration name: " name; \
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) $$name

# Запуск миграции
MIGRATE_BIN=migrate
MIGRATIONS_PATH=/migrations
DATABASE_URL=postgres://postgres:qwerty@db:5432/cruddb?sslmode=disable

migrate-up:
	docker-compose run --rm $(MIGRATE_BIN) -path=$(MIGRATIONS_PATH) -database=$(DATABASE_URL) up

migrate-down:
	docker-compose run --rm $(MIGRATE_BIN) -path=$(MIGRATIONS_PATH) -database=$(DATABASE_URL) down

# Запуск линтера
.PHONY: lint

lint:
	docker-compose run --rm golangci-lint
