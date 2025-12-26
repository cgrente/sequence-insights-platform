SHELL := /bin/bash

.PHONY: help
help:
	@echo "Targets:"
	@echo "  dev		Run api locally"
	@echo "  test		Run tests"
	@echo "  fmt		gofmt"
	@echo "  docker-up  Start local stack with Docker Compose"
	@echo "  docker-down Stop local stack with Docker Compose"

.PHONY: dev
dev:
	cd server && go run ./cmd/api

.PHONY: test
test:
	cd server && go test ./...

.PHONY: fmt
fmt:
	cd server && gofmt -w .

.PHONY: docker-up
docker-up:
	docker compose up --build

.PHONY: docker-down
docker-down:
	docker compose down -v

check:
	cd server && gofmt -w .
	cd server && go vet ./...
	cd server && go test ./...