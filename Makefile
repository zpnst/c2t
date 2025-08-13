BIN_DIR = bin

build-client:
	@go build -o $(BIN_DIR)/c2t-client ./cmd/c2t-client

build-instance:
	@go build -o $(BIN_DIR)/c2t-instance ./cmd/c2t-instance

build-all: build-client build-instance

run-client: build-client
	@./$(BIN_DIR)/c2t-client

run-instance: build-instance
	@./$(BIN_DIR)/c2t-instance

clean:
	@rm -rf $(BIN_DIR)

.PHONY: build-client build-instance build-all run-client run-instance clean