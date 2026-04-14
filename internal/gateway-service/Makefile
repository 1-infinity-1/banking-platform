COMPOSE := docker compose
OGEN_VERSION := v1.18.0

BIN_DIR := $(CURDIR)/bin
OGEN := $(BIN_DIR)/ogen

.PHONY: install-tools generate clean-tools up down

# Install ogen if it is not found in PATH or version does not match
# Устанавливает ogen, если он не найден в PATH или версия не совпадает
install-tools:
	@mkdir -p $(BIN_DIR)
	@if [ ! -x "$(OGEN)" ] || [ "$$($(OGEN) --version 2>/dev/null | awk '{print $$3}')" != "$(OGEN_VERSION)" ]; then \
		echo "Installing ogen $(OGEN_VERSION) to $(BIN_DIR)..."; \
		GOBIN=$(BIN_DIR) go install -v github.com/ogen-go/ogen/cmd/ogen@$(OGEN_VERSION); \
	fi

# Generate code from OpenAPI spec using ogen
# Генерирует код из OpenAPI спецификации с помощью ogen
generate: install-tools
	$(OGEN) -target ./api/ogen -package api ./api/openapi.yaml

# Start services using docker compose in detached mode
# Запускает сервисы через docker compose в фоне
up:
	$(COMPOSE) up -d

# Stop and remove services using docker compose
# Останавливает и удаляет сервисы через docker compose
down:
	$(COMPOSE) down