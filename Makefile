.PHONY: help dev build run test lint fmt tidy clean docker-build docker-run

# Default target — show available commands.
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ---- Development ----
dev: ## Run the server from source (no rebuild needed)
	go run ./cmd/server

# ---- Build ----
build: ## Build the production binary into bin/server
	go build -o bin/server ./cmd/server

run: build ## Build then run the binary
	./bin/server

# ---- Quality ----
test: ## Run the test suite
	go test ./...

lint: ## Run golangci-lint over all packages
	golangci-lint run ./...

fmt: ## Format every Go file (gofmt -s)
	gofmt -s -w .

tidy: ## Sync go.mod / go.sum
	go mod tidy

# ---- Docker ----
docker-build: ## Build the production container image
	docker build -t about-me:latest .

docker-run: docker-build ## Build then run the container locally on :8080
	docker run --rm -p 8080:8080 \
		-e PORTFOLIO_EMAIL_PROVIDER \
		-e PORTFOLIO_EMAIL_RESEND_API_KEY \
		-e PORTFOLIO_EMAIL_FROM \
		-e PORTFOLIO_EMAIL_TO \
		about-me:latest

# ---- Cleanup ----
clean: ## Remove build artifacts
	rm -rf bin/
