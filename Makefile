APP_NAME := LCA
BUILD_DIR := build
SRC := ./cmd/LCA/main.go
Go := go

.PHONY: lint fmt
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...
fmt:
	@echo "Formatting code with goimports..."
	@goimports -w .

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building ${APP_NAME}..."
	@${Go} build -o ${BUILD_DIR}/${APP_NAME} ${SRC}

.PHONY: run
run: build
	@echo "Running ${APP_NAME}..."
	@${BUILD_DIR}/${APP_NAME}

# .PHONY: test
# test:
# 	@echo "Runing tests..."
# 	@${Go} test tests/..

.PHONY: clean
clean:
	@echo "🧹 Cleaning up..."
	@rmdir /S /Q build


