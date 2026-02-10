#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
NC='\033[0m' # No Color

# Check for golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    echo -e "${YELLOW}WARNING: golangci-lint is not installed.${NC}"
    echo -e "${CYAN}To install it, run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest${NC}"
    echo -e "${GRAY}Skipping linting step...${NC}"
    echo ""
else
    echo -e "${CYAN}Running Linter...${NC}"
    golangci-lint run
    if [ $? -ne 0 ]; then
        echo -e "${RED}ERROR: Linting failed!${NC}"
        exit 1
    fi
fi

echo -e "${CYAN}Running Unit Tests...${NC}"
# Check if CGO is enabled
CGO_ENABLED=$(go env CGO_ENABLED)
if [ "$CGO_ENABLED" = "1" ]; then
    go test -v -race ./...
else
    echo -e "${GRAY}INFO: CGO is disabled, running tests without -race flag${NC}"
    go test -v ./...
fi

if [ $? -ne 0 ]; then
    echo -e "${RED}ERROR: Tests failed!${NC}"
    exit 1
fi

echo -e "${GREEN}SUCCESS: All checks passed!${NC}"
