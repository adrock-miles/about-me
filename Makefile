.PHONY: build run dev clean frontend backend

# Build everything
build: frontend backend

# Build Go binary
backend:
	go build -o bin/server ./cmd/server

# Build React frontend
frontend:
	cd web && npm install && npm run build

# Run the Go server (serves built frontend)
run: build
	./bin/server

# Dev mode: run Go server (frontend uses vite dev server with proxy)
dev:
	go run ./cmd/server

# Clean build artifacts
clean:
	rm -rf bin/ web/dist/ web/node_modules/
