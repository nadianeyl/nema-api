## help: list available commands
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/api: run the cmd/api application 
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${NEMA_DB_DSN}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${NEMA_DB_DSN}

## db/migrations/new name=$1: create new database migration files
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path=./migrations -database=${NEMA_DB_DSN} up

## db/migrations/down: revert migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Rollback migrations...'
	migrate -path=./migrations -database=${NEMA_DB_DSN} down

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
