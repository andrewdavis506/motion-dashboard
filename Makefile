.PHONY: build run clean test

# Default binary name
BINARY_NAME=task-dashboard

# Build the application
build:
	@echo "Building application..."
	go build -o $(BINARY_NAME) ./cmd/server

# Run the application
run:
	@echo "Running application..."
	go run ./cmd/server

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -f $(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download

# Build and run in one command
dev: build
	@echo "Starting application in development mode..."
	./$(BINARY_NAME)