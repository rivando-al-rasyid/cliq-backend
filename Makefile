include .env

MIGRATION_PATH=database/migrations
DOCKER_MIGRATION_PATH=/migrations
SEED_FILE=database/seed.sql

DB_CONTAINER_NAME=postgres

DOCKER_DATABASE_URL=postgresql://$(DB_USER):$(DB_PASS)@$(DB_CONTAINER_NAME):5432/$(DB_NAME)?sslmode=disable
LOCAL_DATABASE_URL=postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate-create:
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

migrate-up:
	@docker compose run --rm migrate -database "$(DOCKER_DATABASE_URL)" -path $(DOCKER_MIGRATION_PATH) up

migrate-down:
	@docker compose run --rm migrate -database "$(DOCKER_DATABASE_URL)" -path $(DOCKER_MIGRATION_PATH) down

migrate-force:
	@docker compose run --rm migrate -database "$(DOCKER_DATABASE_URL)" -path $(DOCKER_MIGRATION_PATH) force $(VERSION)

migrate-status:
	@docker compose run --rm migrate -database "$(DOCKER_DATABASE_URL)" -path $(DOCKER_MIGRATION_PATH) version