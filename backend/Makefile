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

GOFLAGS=-ldflags=-extldflags=-Wl,-ld_classic

BIN_API = bin/api
BIN_PARSE_CONFIG = bin/tool/parse-config

$(BIN_API):
	@go build -o $(BIN_API) cmd/api/main.go

GO_MOD_NAME := $(shell go list -m)

test: ## run tests
	@go test -race -cover ./... | sed "/PASS/s//$$(printf '\033[32mPASS\033[0m')/" | sed "/FAIL/s//$$(printf '\033[31mFAIL\033[0m')/"

bench:
	@go test -bench=. -benchmem ./...

########
# lint #
########

lint: ## lints the entire codebase
	@golangci-lint run ./...

###########
#   GCI   #
###########

gci-format: ## format repo through gci linter
	gci write ./ --skip-generated -s standard -s default -s "Prefix($(GO_MOD_NAME))"

##############
# db-migrate #
##############
$(BIN_PARSE_CONFIG):
	@go build -o $(BIN_PARSE_CONFIG) cmd/tool/parse-config/main.go

migrate-up: $(BIN_PARSE_CONFIG) ## run db migration up
	@(\
		DBURL=$$($(BIN_PARSE_CONFIG) get database_url); \
		migrate -database $$DBURL -path ./database/migrations up; \
	)

migrate-down: $(BIN_PARSE_CONFIG) ## run db migration down
	@(\
		DBURL=$$($(BIN_PARSE_CONFIG) get database_url); \
		migrate -database $$DBURL -path ./database/migrations down; \
	)

migrate-create: ## create new migration
	@migrate create -ext sql -dir ./database/migrations $(name)

.test:
	@echo "GO_MOD_NAME: $(GO_MOD_NAME);"
