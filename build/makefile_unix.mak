.PHONY: all build run test test-pretty clean

APP_NAME=kitsune_bot
BIN_DIR=bin
CMD_PATH=./cmd/kitsune_bot

all: build run

build:
	@mkdir -p $(BIN_DIR)
	go mod tidy
	export CGO_ENABLED=1
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_PATH)/main.go

run:
	$(BIN_DIR)/$(APP_NAME)
	@rm -rf $(BIN_DIR)

test:
	go test ./tests/...

test-pretty:
	gest ./tests/...

clean:
	@rm -rf $(BIN_DIR)

