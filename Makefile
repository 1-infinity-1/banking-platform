BIN_DIR := $(CURDIR)/bin
PROTO_DIR := pkg/proto
PROTO_FILE := $(PROTO_DIR)/auth/auth.proto
OUT_DIR := pkg/proto/generated/go

PROTOC := protoc
PROTOC_GEN_GO := $(BIN_DIR)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(BIN_DIR)/protoc-gen-go-grpc

PROTOC_GEN_GO_VERSION := v1.36.11
PROTOC_GEN_GO_GRPC_VERSION := v1.6.1

.PHONY: tools proto clean

tools:
	GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	GOBIN=$(BIN_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

proto: tools
	PATH="$(BIN_DIR):$$PATH" $(PROTOC) \
		-I $(PROTO_DIR) \
		$(PROTO_FILE) \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative