.PHONY: run tidy watch build clean test lint

# Run the application
run:
	@go run cmd/main.go

# Tidy up Go modules
tidy:
	@go mod tidy

# Watch for changes and restart (requires 'air' installed)
watch:
	@air

# Build the binary
build:
	@go build -o bin/app cmd/main.go
	@echo "✅ Build successful!"

# Remove compiled binary and temporary files
clean:
	@rm -rf bin/
	@echo "🧹 Cleaned up!"

# Run tests
test:
	@go test ./... -cover

# Run linter (requires 'golangci-lint' installed)
lint:
	@golangci-lint run

# Default target
.DEFAULT_GOAL := run
