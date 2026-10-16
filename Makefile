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
