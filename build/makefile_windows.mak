.PHONY: all build run test

all: build

build:
	if not exist bin mkdir bin
	go mod tidy
	go build -o bin\kitsune_bot.exe cmd\kitsune_bot\main.go
run:
	bin\kitsune_bot.exe
	rmdir /s /q bin

test:
	go test ./...

test-pretty:
	gest ./...
