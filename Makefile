SHELL=/bin/bash
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))
name := $(shell head -1 $(current_dir)/go.mod|sed -e 's,^.*/,,g')

.DEFAULT_GOAL := run

depends_cmds := go
check:
	@for cmd in ${depends_cmds}; do command -v $$cmd >&/dev/null || (echo "No $$cmd command" && exit 1); done

clean:
	@for d in $(name); do if [[ -e $${d} ]]; then echo "==> Removing $${d}.." && rm -rf $${d}; fi done

run: check clean
	@LOG_LEVEL=debug go run . ./test/test.xlsx

help:
	@go run ./main.go -h

build-linux:
	@make GOOS=linux _build
build-mac:
	@make GOOS=darwin _build

_build: check clean
	@env GOOS=$(GOOS) go build -ldflags="-s -w"

deps:
	@go list -m all

tidy:
	@go mod tidy
