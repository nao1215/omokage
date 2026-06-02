.PHONY: build test lint clean tools help

APP         = omokage
VERSION     = $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)
GO          = go
GO_BUILD    = $(GO) build
GO_TEST     = $(GO) test
GO_TOOL     = $(GO) tool
GO_INSTALL  = $(GO) install
GOOS        = ""
GOARCH      = ""
GO_PKGROOT  = ./...

build: ## Build binary
	env GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) -o $(APP) main.go

clean: ## Clean project
	-rm -rf $(APP) coverage.out coverage.html

test: ## Run tests with coverage output
	env GOOS=$(GOOS) $(GO_TEST) -cover -covermode=atomic -coverpkg=$(GO_PKGROOT) -coverprofile=coverage.out $(GO_PKGROOT)
	$(GO_TOOL) cover -html=coverage.out -o coverage.html

lint: ## Run golangci-lint
	golangci-lint run --config .golangci.yml

tools: ## Install developer tools
	$(GO_INSTALL) github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	$(GO_INSTALL) github.com/k1LoW/octocov@latest

.DEFAULT_GOAL := help
help:
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'
