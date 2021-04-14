VENDORDIR = vendor

GO ?= go
GOLANGCI_LINT ?= golangci-lint

.PHONY: $(VENDORDIR) lint test test-unit

$(VENDORDIR):
	@mkdir -p $(VENDORDIR)
	@$(GO) mod vendor
	@$(GO) mod tidy

lint:
	@$(GOLANGCI_LINT) run

test: test-unit test-integration

## Run unit tests
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

## Run integration tests
test-integration:
	@echo ">> integration test"
	@$(GO) test -gcflags=-l -coverprofile=integration.coverprofile -covermode=atomic -race ./... --tags integration
