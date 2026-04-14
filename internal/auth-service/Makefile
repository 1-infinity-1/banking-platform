COMPOSE := docker compose

.PHONY: up down

# Start services using docker compose in detached mode
# Запускает сервисы через docker compose в фоне
up:
	$(COMPOSE) up -d

# Stop and remove services using docker compose
# Останавливает и удаляет сервисы через docker compose
down:
	$(COMPOSE) down