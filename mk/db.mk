MIGRATION_NAME ?= unnamed

.PHONY: migrate
migrate: ## Run database migration
	go run ./cmd/migrate

.PHONY: generate-migration
MIGRATION_NAME?=unnamed
generate-migration: ## Generate new migration file, use MIGRATION_NAME=name
	tern new -m internal/db/migrations $(MIGRATION_NAME)