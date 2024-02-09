NAME ?= common-go-dependency-tree
export PACKAGE_NAME ?= $(NAME)
ifeq ($(OS),Windows_NT)
	export VERSION=$(shell .\scripts\workflows\current-version.bat -f VERSION)
else
	export VERSION=$(shell ./scripts/workflows/current-version.sh -f VERSION)
endif

COBERTURA = cobertura

GOX = gox

GOLANGCI_LINT = golangci-lint

DEVELOPMENT_TOOLS = $(GOX) $(COBERTURA) $(GOLANGCI_LINT)

GET_CURRENT_LINT_CONTAINER = $(shell docker ps -a -q -f "name=$(PACKAGE_NAME)-linter")


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
ifeq ($(OS),Windows_NT)
	@scripts\test.bat -d ./pkg
else
	@scripts/test -d ./pkg
endif

.PHONY: coverage
coverage:
	@echo "Running coverage report..."
ifeq ($(OS),Windows_NT)
	@scripts\coverage.bat -d ./pkg
else
	@scripts/coverage -d ./pkg
endif

.PHONY: lint
lint:
	@echo "Running linter..."
ifeq ($(GET_CURRENT_LINT_CONTAINER),)
	@echo "Linter container does not exist, creating it..."
	@-docker run --name $(PACKAGE_NAME)-linter -e RUN_LOCAL=true -e VALIDATE_ALL_CODEBASE=true -e VALIDATE_JSCPD=false -e CREATE_LOG_FILE=true -e LOG_FILE=lint.log -v .:/tmp/lint ghcr.io/super-linter/super-linter:slim-v5
else
	@echo "Linter container already exists, starting it..."
	@-docker start $(PACKAGE_NAME)-linter --attach
endif
	@docker cp $(PACKAGE_NAME)-linter:/tmp/lint/lint.log ./lint-report.log
	@echo "Linter report saved to lint-report.log"
	@echo "Linter finished."


.PHONY: build
build:
	@echo "Building..."
ifeq ($(OS),Windows_NT)
	@scripts\build.bat -d ./pkg -p $(PACKAGE_NAME)
else
	@scripts/build -d ./pkg -p $(PACKAGE_NAME)
endif

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

$(START_CURRENT_LINT_CONTAINER):
	@echo "Linter container already exists, starting it..."
	@docker start $(PACKAGE_NAME)-linter --attach