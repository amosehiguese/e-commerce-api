# App
APP_NAME := ecommerce-api
BUILD_DIR := bin

# Tools
GO := go

.PHONY: all build run test clean


all: build

build:
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) main.go
	@echo "Build completed."

run: build
	$(BUILD_DIR)/$(APP_NAME)

test:
	$(GO) test ./... -v
	@echo "Tests executed."

lint:
	golangci-lint run
	@echo "Linting completed."

clean:
	rm -rf $(BUILD_DIR)
	@echo "Cleaned up build artifacts."
