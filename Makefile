include .envrc
MIGRATIONS_PATH=./migrations


.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@, $(MAKECMDGOALS))

.PHONY:migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY:migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@, $(MAKECMDGOALS))

.PHONY: test
test:
	@go test -v ./...

dev:
	@air