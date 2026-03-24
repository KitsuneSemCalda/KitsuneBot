.PHONY: all build run tests

all: build run

build:
	mkdir bin
	go build -o bin/kitsune_bot cmd/kitsune_bot/main.go

run:
	cls
	./bin/kitsune_bot

tests:
	cls
	gest ./...
