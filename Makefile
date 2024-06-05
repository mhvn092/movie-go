MAIN_PACKAGE_PATH := ./cmd/service
MIGRATION_PACKAGE_PATH := ./cmd/migration
BINARY_NAME := movie
MIGRATION_NAME := migration

build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

run: build
	/tmp/bin/${BINARY_NAME}

migrate:
	go build -o=/tmp/bin/${MIGRATION_NAME} ${MIGRATION_PACKAGE_PATH}

up: migrate
	/tmp/bin/${MIGRATION_NAME}
