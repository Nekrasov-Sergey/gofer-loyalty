.DEFAULT_GOAL := run

.PHONY: run
run: build
	@./cmd/gophermart/main $(args)

.PHONY: build
build:
	@go build -o ./cmd/gophermart/main ./cmd/gophermart/main.go

.PHONY: accrual
accrual:
	@RUN_ADDRESS=localhost:8081 ./cmd/accrual/accrual_darwin_arm64

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: tidy
tidy:
	@go mod tidy
	@go mod verify

.PHONY: deps
deps:
	@go get -u ./...
	@go mod tidy
	@go mod verify
	@go build ./...
