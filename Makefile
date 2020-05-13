# go option
GO        ?= go
TAGS      :=
TESTS     := ./...
TESTFLAGS := -race
LDFLAGS   :=
GOFLAGS   :=
BINARIES  := dancing-links
VERSION   := v0.1.0

# Required for globs to work correctly
SHELL=/bin/bash

.DEFAULT_GOAL := unit-test

.PHONY: release
release:
	@echo
	@echo "==> Preparing the release $(VERSION) <=="
	go mod tidy
	git tag ${VERSION}

.PHONY: bench
bench:
	$(GO) test -v -benchmem -bench=. ./sudoku

.PHONY: unit-test
unit-test:
	@echo
	@echo "==> Running unit tests <=="
	$(GO) test $(GOFLAGS) $(TESTS) $(TESTFLAGS)
