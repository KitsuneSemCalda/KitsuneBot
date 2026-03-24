.PHONY: all build run test

all: build run

build:
	if not exist bin mkdir bin
	go build -o bin\kitsune_bot.exe cmd\kitsune_bot\main.go
run:
	bin\kitsune_bot.exe
	rmdir /s /q bin

test:
	cls
	go test ./...

test-pretty:
	cls
	gest ./...