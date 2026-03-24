.PHONY: all build run tests

ifeq ($(OS),Windows_NT)
    include build/makefile_windows.mak
else ifeq ($(UNAME_S),Linux)
    include build/makefile_linux.mak
endif
