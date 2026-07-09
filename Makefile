.PHONY: all build build-daemon build-dashboard dev dev-daemon dev-dashboard clean check test docker-build docker-up docker-down

BINARY_NAME=vesseld
BUILD_DIR=bin

all: check build

check:
	@echo "🔍 Running Go checks and formatting..."
	go fmt ./...
	go vet ./...

test:
	@echo "🧪 Running full test suite..."
	go test ./... -v

build: build-dashboard build-daemon
	@echo "✅ Build complete! Binary available at $(BUILD_DIR)/$(BINARY_NAME) and GUI at dashboard/dist"

build-daemon:
	@echo "⚙️  Building Go daemon binary ($(BINARY_NAME))..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/vesseld

build-dashboard:
	@echo "💻 Building TanStack + Vite Dashboard GUI..."
	npm run build:dashboard

dev:
	@echo "🚀 Launching backend daemon and frontend GUI concurrently..."
	npx concurrently -k "make dev-daemon" "make dev-dashboard"

dev-daemon:
	@echo "🚀 Running Go daemon in dev mode..."
	go run ./cmd/vesseld

dev-dashboard:
	@echo "💻 Running Dashboard dev server on port 3000..."
	npm run dev:dashboard

dev-website:
	@echo "🌐 Running Astro Marketing site dev server..."
	npm run dev:website

docker-build:
	@echo "🐳 Building Docker image..."
	docker compose build

docker-up:
	@echo "🐳 Starting Vessel via Docker Compose..."
	docker compose up -d

docker-down:
	@echo "🐳 Stopping Vessel Docker stack..."
	docker compose down

clean:
	@echo "🧹 Cleaning builds and temporary binaries..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
