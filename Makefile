APP_NAME := LCA
BUILD_DIR := build
SRC := ./cmd/LCA/main.go
Go := go

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

.PHONY: test
test:
	@echo "Runing tests..."
	@${Go} test tests/..

.PHONY: clean

ifeq ($(OS),Windows_NT)
RM := del /Q /S
RMDIR := rmdir /S /Q
else
RM := rm -f
RMDIR := rm -rf
endif

clean:
	@echo "🧹 Cleaning up..."
	@$(RMDIR) build


