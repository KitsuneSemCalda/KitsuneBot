.PHONY: all build run tests

ifeq ($(OS),Windows_NT)
    include build/makefile_windows.mak
endif