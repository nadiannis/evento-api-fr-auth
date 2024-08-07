include .env

DB_DSN = ${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}

ifeq ($(strip $(PORT)),)
  PORT = 8080
endif

## help: list available commands
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run: run the application 
.PHONY: run
run:
	go run ./cmd -port=${PORT} -db-dsn=postgres://${DB_DSN}?sslmode=disable -jwt-secret=${JWT_SECRET}

## db/migrations/new name=$1: create new database migration files 
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Apply all migrations...'
	migrate -path ./migrations -database pgx://${DB_DSN} up

## db/migrations/down: revert all migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Rollback migrations...'
	migrate -path ./migrations -database pgx://${DB_DSN} down

## audit: tidy dependencies, format code, & vet code
.PHONY: audit
audit:
	@echo 'Tidying & verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...