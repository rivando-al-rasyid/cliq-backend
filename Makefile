include ./.env
export # This automatically exports all variables from .env to sub-shells

MIGRATION_PATH=database/migrations
SEED_FILE=database/seed.sql

# ── Docker Config ──────────────────────────────────────────────────────────
# Change this to match your DB service name in docker-compose.yml
DB_CONTAINER_NAME=db 

# Internal Docker URL (used when running commands via docker compose exec)
DOCKER_DATABASE_URL=postgresql://$(DB_USER):$(DB_PASS)@$(DB_CONTAINER_NAME):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Local URL (used if running migrate tool directly from your host machine)
LOCAL_DATABASE_URL=postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# ── Migration ──────────────────────────────────────────────────────────────
migrate-create:
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

# Runs migrate using a temporary docker container inside the network
migrate-up:
	@docker compose run --rm migrate -database $(DOCKER_DATABASE_URL) -path $(MIGRATION_PATH) up

migrate-down:
	@docker compose run --rm migrate -database $(DOCKER_DATABASE_URL) -path $(MIGRATION_PATH) down

migrate-force:
	@docker compose run --rm migrate -database $(DOCKER_DATABASE_URL) -path $(MIGRATION_PATH) force $(VERSION)

migrate-status:
	@docker compose run --rm migrate -database $(DOCKER_DATABASE_URL) -path $(MIGRATION_PATH) version

# ── Seed ───────────────────────────────────────────────────────────────────
# Executes psql directly inside the already running DB container
seed:
	@docker compose exec -T $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) < $(SEED_FILE)

seed-reset:
	@docker compose exec -T $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) -c \
		"TRUNCATE TABLE withdrawals, expenses, transfers, transactions, topups, wallets, user_pins, favorites, profiles, users RESTART IDENTITY CASCADE;"
	@$(MAKE) seed

# ── Docker ─────────────────────────────────────────────────────────────────
docker-up: 
	@docker compose up -d --build

docker-down: 
	@docker compose down -v