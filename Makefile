BIN_DIR := $(CURDIR)/bin
PROTO_DIR := pkg/proto
PROTO_FILE := $(PROTO_DIR)/auth/auth.proto
OUT_DIR := pkg/proto/generated/go

PROTOC := protoc
PROTOC_GEN_GO := $(BIN_DIR)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(BIN_DIR)/protoc-gen-go-grpc

PROTOC_GEN_GO_VERSION := v1.36.11
PROTOC_GEN_GO_GRPC_VERSION := v1.6.1

GOLANGCI_LINT_VERSION := v2.11.4

.PHONY: tools proto clean install-lint lint lint-fix

# Install required protobuf code generation tools
# Устанавливает необходимые инструменты для генерации protobuf кода
tools:
	GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	GOBIN=$(BIN_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

# Generate Go code from .proto files
# Генерирует Go код из .proto файлов
proto: tools
	PATH="$(BIN_DIR):$$PATH" $(PROTOC) \
		-I $(PROTO_DIR) \
		$(PROTO_FILE) \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative

# Install golangci-lint locally into project bin directory
# Устанавливает golangci-lint в локальную директорию bin проекта
install-lint:
	mkdir -p $(BIN_DIR)
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(BIN_DIR) $(GOLANGCI_LINT_VERSION)

# Run linters on the entire project
# Запускает линтеры для всего проекта
lint:
	$(BIN_DIR)/golangci-lint run ./...

# Run linters with automatic fixes where possible
# Запускает линтеры с автоматическим исправлением проблем (где возможно)
lint-fix:
	$(BIN_DIR)/golangci-lint run --fix ./...