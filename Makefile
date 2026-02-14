# --- Миграции ---
.PHONY: migrate-add
migrate-add: ## Create new migration file, usage: migrate-add [name=<migration_name>]
	@echo ">> Creating new migration: $(NAME)"
	goose -dir database/migrations create $(name) sql

# --- Линтер ---
.PHONY: lint
lint:
	docker compose --profile devtools run --rm --no-deps -T lint