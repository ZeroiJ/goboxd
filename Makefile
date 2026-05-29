.PHONY: build run test integration lint load

COMPOSE ?= docker compose
TOOLS   := $(COMPOSE) --profile tools run --rm -e GOFLAGS=-buildvcs=false tools

build:
	$(COMPOSE) build goboxd

run:
	$(COMPOSE) up goboxd

test:
	$(TOOLS) go test ./...

integration:
	$(COMPOSE) up -d goboxd
	$(TOOLS) go test -tags=integration ./tests/...

lint:
	$(TOOLS) golangci-lint run ./...

load:
	@echo "Running k6 load test..."
	docker run --rm -i --network host -v $(PWD)/scripts:/scripts grafana/k6 run /scripts/load_test.js
