.PHONY: all build run test test-pretty clean

ifeq ($(OS),Windows_NT)
    include build/makefile_windows.mak
else
    UNAME_S := $(shell uname -s 2>/dev/null || echo Unknown)
    ifeq ($(UNAME_S),Linux)
        include build/makefile_unix.mak
    else ifeq ($(UNAME_S),Darwin)
        include build/makefile_unix.mak
    else
        $(warning Unsupported OS: $(UNAME_S) - defaulting to Unix makefile)
        include build/makefile_unix.mak
    endif
endif
