.PHONY: generate
generate:
	go generate ./...

.PHONY: up
up:
	docker compose up -d

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	gofumpt -w .
	gci write . --skip-generated -s standard -s default

.PHONY: lint
lint: tidy fmt
	golangci-lint run

.DEFAULT_GOAL := lint