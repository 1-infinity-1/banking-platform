COMPOSE := docker compose

.PHONY: up down

# Start all ledger-service containers (grpc, consumer, db, kafka, migrate)
# Запускает все контейнеры ledger-service
up:
	$(COMPOSE) up -d

# Stop and remove all ledger-service containers
# Останавливает и удаляет все контейнеры ledger-service
down:
	$(COMPOSE) down
