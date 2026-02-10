.PHONY: test lint build tidy

# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Linting
LINTCMD=golangci-lint run

all: lint test build

tidy:
	$(GOMOD) tidy

test:
	$(GOTEST) -v -race ./...

lint:
	$(LINTCMD)

build:
	$(GOCMD) build ./...

# Example command to run only unit tests in a specific package
# make test-pkg PKG=./internal/utils
test-pkg:
	$(GOTEST) -v $(PKG)
