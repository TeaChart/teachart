
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
GOLANGCI_LINT_PACKAGE ?= github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

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
	@echo " - watch                            watch everything and continuously rebuild"
	@echo " - deps                             install dependencies"
	@echo " - deps-mod                         install go mod dependencies"
	@echo " - deps-tools                       install tool dependencies"

.PHONY: build
build:
	$(GO) build -ldflags="$(VERSION_LDFLAGS)"

.PHONY: watch
watch:
	$(GO) run $(AIR_PACKAGE) -c .air.toml

.PHONY: deps
deps: deps-mod deps-tools

.PHONY: deps-mod
deps:
	$(GO) mod download

.PHONY: deps-tools
deps-tools:
	$(GO) install $(COBRA_CLI_PACKAGE)
	$(GO) install $(GOLANGCI_LINT_PACKAGE)

.PHONY: lint
lint:
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run

.PHONY: lint-fix
lint-fix:
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run --fix