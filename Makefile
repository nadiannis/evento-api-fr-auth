DB_DSN = pgx://postgres:pass1234@localhost:5432/eventodb

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
	go run ./cmd

## db/migrations/new name=$1: create new database migration files 
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Apply all migrations...'
	migrate -path ./migrations -database ${DB_DSN} up

## db/migrations/down: revert all migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Rollback migrations...'
	migrate -path ./migrations -database ${DB_DSN} down

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