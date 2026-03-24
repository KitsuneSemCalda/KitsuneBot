.PHONY: all build run tests

UNAME_S := $(shell uname -s)

ifeq ($(OS),Windows_NT)
    include build/makefile_windows.mak
else ifeq ($(UNAME_S),Linux)
    include build/makefile_linux.mak
else ifeq ($(UNAME_S),Darwin)
    include build/makefile_unix.mak  # macOS usa mesma lógica do Linux
else
    $(error Unsupported OS: $(UNAME_S))
endif
