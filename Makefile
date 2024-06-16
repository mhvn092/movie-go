MAIN_PACKAGE_PATH := ./cmd/service
MIGRATION_PACKAGE_PATH := ./cmd/migration
BINARY_NAME := movie
MIGRATION_NAME := migration

.PHONY: build run migrate create up down

build:
	go build -o /tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

run: build
	/tmp/bin/${BINARY_NAME}

migrate:
	go build -o /tmp/bin/${MIGRATION_NAME} ${MIGRATION_PACKAGE_PATH}

create: migrate
ifndef NAME
	@echo "NAME is not set. Usage: make create NAME=\"<name>\""
	exit 1
endif
	@echo "Provided NAME: $(NAME). You Should wrap your name inside double quotes"
	/tmp/bin/${MIGRATION_NAME} create -name $(NAME)

up: migrate
	/tmp/bin/${MIGRATION_NAME} up

down: migrate
	/tmp/bin/${MIGRATION_NAME} down
