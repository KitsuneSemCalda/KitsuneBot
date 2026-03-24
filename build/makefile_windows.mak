.PHONY: all build run tests

all: build run

build:
	if not exist bin mkdir bin
	go build -o bin\kitsune_bot.exe cmd\kitsune_bot\main.go
run:
	bin\kitsune_bot.exe
	rmdir /s /q bin

tests:
	cls
	go test ./...