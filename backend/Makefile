.PHONY: build run dev test lint clean docker-build docker-up migrate

APP_NAME=instaforge
BUILD_DIR=./build
MAIN_PATH=./cmd/server

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run:
	@go run $(MAIN_PATH)/main.go

dev:
	@air \
		--build.cmd "go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)" \
		--build.bin "$(BUILD_DIR)/$(APP_NAME)" \
		--build.delay "100" \
		--build.exclude_dir "vendor,kubernetes,scripts" \
		--build.include_ext "go,env" \
		--build.stop_on_error "true"

test:
	@go test -v -race -count=1 -coverprofile=coverage.out ./...

test-short:
	@go test -v -short -count=1 ./...

lint:
	@golangci-lint run ./...

clean:
	@rm -rf $(BUILD_DIR) tmp coverage.out
	@echo "Cleaned up"

docker-build:
	@docker build -t $(APP_NAME):latest .

docker-up:
	@docker-compose up -d --build

docker-down:
	@docker-compose down

migrate:
	@go run scripts/migrate/main.go

proto:
	@echo "No protobuf files yet"

bench:
	@go test -bench=. -benchmem ./...

.PHONY: coverage
coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
