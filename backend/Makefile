.PHONY: help \
	test \
	lint \
	gci-format \
	migrate-up \
	migrate-down \
	migrate-create

help: ## show this help
	@echo "Usage: make [target]"
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

SHELL = /bin/bash

test: ## run tests
	@go test -v -race -cover ./... | sed "/PASS/s//$$(printf '\033[32mPASS\033[0m')/" | sed "/FAIL/s//$$(printf '\033[31mFAIL\033[0m')/"

bench:
	@go test -bench=. -benchmem ./...


########
# lint #
########

lint: ## lints the entire codebase
	@golangci-lint run ./... --config=./.golangci.toml

###########
#   GCI   #
###########

gci-format: ## format repo through gci linter
	gci write ./ --skip-generated -s standard -s default -s "Prefix(goadmin-backend)"

##############
# db-migrate #
##############
migrate-up: ## run db migration up
	@go run ./cmd/db-migrate/main.go up

migrate-down: ## run db migration down
	@go run ./cmd/db-migrate/main.go down

migrate-create: ## create new migration
	@go run ./cmd/db-migrate/main.go create
