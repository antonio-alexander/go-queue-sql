## ----------------------------------------------------------------------
## This makefile can be used to execute common functions to interact with
## the source code, these functions ease local development and can also be
## used in CI/CD pipelines.
## ----------------------------------------------------------------------

# REFERENCE: https://stackoverflow.com/questions/16931770/makefile4-missing-separator-stop
help: ## - Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

check-lint: ## - validate/install golangci-lint installation
	@which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1)

lint: check-lint ## - lint the source with verbose output
	@golangci-lint run --verbose

check-godoc: ## - validate/install godoc
	@which godoc || (go install golang.org/x/tools/cmd/godoc@v0.1.10)

serve-godoc: check-godoc ## - serve (web) the godocs
	@godoc -http :8080

test: ## - test the source
	@go test --count=1 -p 1 -cover --count=1 ./...

test-verbose: ## - test the source with verbose output
	@go test --count=1 -p 1 -cover -v --count=1 ./...

run: ## - run the dependencies
	@docker compose up -d

stop: ## - stop the dependencies
	@docker compose down