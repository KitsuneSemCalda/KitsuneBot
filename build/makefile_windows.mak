.PHONY: all build run test test-pretty clean

APP_NAME = kitsune_bot
BIN_DIR  = bin
CMD_PATH = cmd\kitsune_bot

all: build run

build:
	if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	go mod tidy
	@cmd /C "set CGO_ENABLED=1 && go build -o $(BIN_DIR)\$(APP_NAME).exe $(CMD_PATH)\main.go" || powershell -Command "$env:CGO_ENABLED=1; go build -o $(BIN_DIR)\$(APP_NAME).exe $(CMD_PATH)\main.go"
run:
	$(BIN_DIR)\$(APP_NAME).exe

test:
	go test ./...

test-pretty:
	gest ./...

clean:
	if exist $(BIN_DIR) rmdir /s /q $(BIN_DIR)