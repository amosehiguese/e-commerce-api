# App
APP_NAME := ecommerce-api
BUILD_DIR := bin

# Tools
GO := go

# Environment variables
ENV ?= development

# Paths to .env files
ENV_FILE = .env
ENV_FILE_DEV = .env.dev
ENV_FILE_PROD = .env.prod

# Load environment variables based on the environment
ifneq ("$(wildcard $(ENV_FILE))","")
    include $(ENV_FILE)
    export
endif

ifeq ($(ENV),development)
    ifneq ("$(wildcard $(ENV_FILE_DEV))","")
        include $(ENV_FILE_DEV)
        export
    endif
endif

ifeq ($(ENV),production)
    ifneq ("$(wildcard $(ENV_FILE_PROD))","")
        include $(ENV_FILE_PROD)
        export
    endif
endif

.PHONY: all build run test lint clean init_db clean_db init_and_build init_and_test init_and_run

all: init_and_build

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

clean_db:
	-docker stop postgresql
	sleep 3
	@echo "Cleaned up PostgreSQL Docker Container"

init_db:clean_db
	./scripts/init-db.sh
	@echo "Database initialized."

init_and_build: init_db build

init_and_test: init_db test

init_and_run: init_db run
