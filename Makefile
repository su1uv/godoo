include .envrc

.PHONY: confirm
confirm:
	@echo "Are you sure? [y/N]" && read ans && [ $${ans:-N} = y ]


.PHONY: run/api
# run/api: confirm
run/api:
	@echo "Starting application"
	go run ./cmd/api -db-dsn=${GODOO_DB_DSN}

.PHONY: db/psql
db/psql:
	@psql ${GODOO_DB_DSN}


# Goose migrations
.PHONY: migrate/up
migrate/up:
	@echo "Running migrations"
	@goose postgres "${GODOO_DB_DSN}" up -dir ./migrations

.PHONY: migrate/down
migrate/down:
	@echo "Rolling back migrations"
	@goose postgres "${GODOO_DB_DSN}" down -dir ./migrations

.PHONY: migrate/reset
migrate/reset:
	@echo "Resetting migrations"
	@goose postgres "${GODOO_DB_DSN}" reset -dir ./migrations

.PHONY: migrate/status
migrate/status:
	@echo "Checking migration status"
	@goose postgres "${GODOO_DB_DSN}" status -dir ./migrations
