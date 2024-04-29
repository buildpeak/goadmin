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

BIN_API = bin/api
BIN_PARSE_CONFIG = bin/tool/parse-config

$(BIN_API):
	@go build -o $(BIN_API) cmd/api/main.go

GO_MOD_NAME := $(shell go list -m)

test: ## run tests
	@go test -race -cover ./... | \
		sed "s/PASS/$$(printf '\033[32mPASS\033[0m')/" | \
		sed "s/FAIL/$$(printf '\033[31mFAIL\033[0m')/" | \
		sed "s/ok/$$(printf '\033[32mok\033[0m')/" | \
		sed -E "s/coverage: ([89]...%|100.0%)/$$(printf 'coverage: \033[32m')\1$$(printf '\033[0m')/" | \
		sed -E "s/coverage: ([6-7]...%)/$$(printf 'coverage: \033[33m')\1$$(printf '\033[0m')/" | \
		sed -E "s/coverage: ([0-5]...%|[0-9]..%)/$$(printf 'coverage: \033[31m')\1$$(printf '\033[0m')/"

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
	@(\
		if [ -z "$(name)" ]; then \
			echo "Usage: make migrate-create name=<migration_name>"; \
			exit 0; \
		else \
			migrate create -ext sql -dir ./database/migrations $(name); \
		fi; \
	)

.test:
	@echo "GO_MOD_NAME: $(GO_MOD_NAME);"
