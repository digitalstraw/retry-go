.PHONY: lint test fmt tidy

# Run all formatting and linting
lint:
	@echo "=> Running linters and formatters..."
	@go fmt ./...
	@gofmt -w .
	@golangci-lint run --fix --build-tags "tests"

# Run tests with gotestsum and test build tags
tests:
	@echo "=> Running tests..."
	@LOCAL_DEV=true gotestsum -- -tags=tests ./...

# Run only formatting (useful before commits)
fmt:
	@echo "=> Formatting code..."
	@go fmt ./...
	@gofmt -w .

# Clean up and organize go.mod / go.sum
tidy:
	@echo "=> Tidying Go modules..."
	@go mod tidy

# Run the application
run:
	@echo "=> Starting application..."
	@go run -tags tests .
