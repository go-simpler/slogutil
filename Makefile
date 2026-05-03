.POSIX:
.SUFFIXES:

help:
	@echo 'Available commands:'
	@echo '  build       Build the project'
	@echo '  fmt         Run formatters'
	@echo '  lint        Run linters'
	@echo '  test        Run tests'
	@echo '  test/bench  Run tests and benchmarks'
	@echo '  test/cover  Run tests and open coverage report'

build:
	@go build -o /dev/null ./...

fmt:
	@golangci-lint fmt

lint:
	@golangci-lint run --fix

test:
	@go test -race -shuffle=on -coverprofile=coverage.out ./...

test/bench: test
	@go test -bench ./...

test/cover: test
	@go tool cover -html=coverage.out
