
GO ?= go

TAG  ?= $(shell git describe --tags --abbrev=0 HEAD)
DATE_FMT = +"%Y-%m-%dT%H:%M:%S%z"
BUILD_DATE ?= $(shell date "$(DATE_FMT)")

VERSION_LDFLAGS=\
  -X go.szostok.io/version.version=$(TAG) \
  -X go.szostok.io/version.buildDate=$(BUILD_DATE) \
  -X go.szostok.io/version.commit=$(shell git rev-parse --short HEAD) \
  -X go.szostok.io/version.commitDate=$(shell git log -1 --date=format:"%Y-%m-%dT%H:%M:%S%z" --format=%cd) \
  -X go.szostok.io/version.dirtyBuild=false

COBRA_CLI_PACKAGE ?= github.com/spf13/cobra-cli@v1.3.0
GOLANGCI_LINT_PACKAGE ?= github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.1
MISSPELL_PACKAGE ?= github.com/client9/misspell/cmd/misspell@v0.3.4
GOVULNCHECK_PACKAGE ?= golang.org/x/vuln/cmd/govulncheck@v1.0.4

TEST_TAGS ?=
GOTESTFLAGS ?=

GO_DIRS := build cmd pkg tests

ifeq ($(IS_WINDOWS),yes)
	GOFLAGS := -v -buildmode=exe
	EXECUTABLE ?= teachart.exe
else
	GOFLAGS := -v
	EXECUTABLE ?= teachart
endif

.PHONY: help
help:
	@echo "Make Help:"
	@echo " - \"\"                               equivalent to \"build\""
	@echo " - build                            build everything"
	@echo " - checks                           checks go files"
	@echo " - deps                             install dependencies"
	@echo " - deps-mod                         install go mod dependencies"
	@echo " - deps-tools                       install tool dependencies"
	@echo " - lint                             check lint"
	@echo " - lint-fix                         check lint and fix"
	@echo " - test                             run go test"
	@echo " - tidy                             run go tidy"

.PHONY: build
build:
	$(GO) build -ldflags="$(VERSION_LDFLAGS)"

.PHONY: deps
deps: deps-mod deps-tools

.PHONY: deps-mod
deps:
	$(GO) mod download

.PHONY: deps-tools
deps-tools:
	$(GO) install $(COBRA_CLI_PACKAGE)
	$(GO) install $(GOLANGCI_LINT_PACKAGE)
	$(GO) install $(GOVULNCHECK_PACKAGE)
	$(GO) install $(MISSPELL_PACKAGE)

.PHONY: lint
lint:
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run --timeout=5m

.PHONY: lint-fix
lint-fix:
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run --fix

.PHONY: test
test:
	@echo "Running go test with -tags '$(TEST_TAGS)'"
	@$(GO) test $(GOTESTFLAGS) -tags='$(TEST_TAGS)' ./... .

.PHONY: checks
checks: tidy-check misspell-check security-check

.PHONY: tidy
tidy:
	$(eval MIN_GO_VERSION := $(shell grep -Eo '^go\s+[0-9]+\.[0-9.]+' go.mod | cut -d' ' -f2))
	$(GO) mod tidy -compat=$(MIN_GO_VERSION)

.PHONY: tidy-check
tidy-check: tidy
	@diff=$$(git diff --color=always go.mod go.sum); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make tidy' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi

.PHONY: misspell-check
misspell-check:
	go run $(MISSPELL_PACKAGE) -error $(GO_DIRS)

.PHONY: security-check
security-check:
	go run $(GOVULNCHECK_PACKAGE) ./...
