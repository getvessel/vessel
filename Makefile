.PHONY: all build-daemon build-dashboard dev-daemon dev-dashboard clean check

BINARY_NAME=vesseld
BUILD_DIR=bin

all: check build-daemon build-dashboard

check:
	@echo "🔍 Running Go checks and formatting..."
	go fmt ./...
	go vet ./...

build-daemon:
	@echo "⚙️  Building Go daemon binary ($(BINARY_NAME))..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/vesseld

build-dashboard:
	@echo "💻 Building TanStack + Vite Dashboard GUI..."
	cd dashboard && npm run build

dev-daemon:
	@echo "🚀 Running Go daemon in dev mode..."
	go run ./cmd/vesseld

dev-dashboard:
	@echo "💻 Running Dashboard dev server on port 3000..."
	cd dashboard && npm run dev

dev-website:
	@echo "🌐 Running Astro Marketing site dev server..."
	cd website && npm run dev

clean:
	@echo "🧹 Cleaning builds and temporary binaries..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
