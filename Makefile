NAME ?= "common-go-dependency-tree"
export PACKAGE_NAME ?= $(NAME)
export VERSION=$(shell ./scripts/workflows/current-version.sh -f VERSION)

COBERTURA = cobertura

GOX = gox

GOLANGCI_LINT = golangci-lint

DEVELOPMENT_TOOLS = $(GOX) $(COBERTURA) $(GOLANGCI_LINT)

.PHONY: help
help:
  # make version:
	# make test
	# make lint

.PHONY: version
version:
	@echo Version: $(VERSION)


.PHONY: test
test:
	@echo "Running tests..."
	@scripts/test -d ./src

.PHONY: coverage
coverage:
	@echo "Running coverage report..."
	@scripts/coverage -d ./src

.PHONY: lint
lint:
	@echo "Running linter..."
	@scripts/lint

.PHONY: build
build:
	@echo "Building..."
	@scripts/build -d ./src -p $(PACKAGE_NAME)

.PHONY: deps
deps: $(DEVELOPMENT_TOOLS)

$(COBERTURA):
	@echo "Installing cobertura..."
	@go install github.com/axw/gocov/gocov@latest
	@go install github.com/AlekSi/gocov-xml@latest
	@go install github.com/matm/gocov-html/cmd/gocov-html@latest

$(GOX):
	@echo "Installing gox..."
	@go install github.com/mitchellh/gox@latest


$(GOLANGCI_LINT):
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint