.PHONY: build run test integration lint

COMPOSE ?= docker compose
TOOLS   := $(COMPOSE) --profile tools run --rm -e GOFLAGS=-buildvcs=false tools

build:
	$(COMPOSE) build goboxd

run:
	$(COMPOSE) up goboxd

test:
	$(TOOLS) go test ./...

integration:
	$(TOOLS) go test -tags=integration ./tests/...

lint:
	$(TOOLS) golangci-lint run ./...
