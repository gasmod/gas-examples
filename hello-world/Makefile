.PHONY: help build run

# Default target
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Building
build: ## Build the binary
	go build ./...

# Running
run: ## Run the demo
	go run .
