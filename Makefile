# Default target: build the project
all: build

# Build the Go project and output the binary to bin/quantum
build:
    go build -o bin/quantum ./...

# Run all unit tests
test:
    go test ./...

# Format all Go source files
fmt:
    go fmt ./...

# Run static code analysis (requires golangci-lint)
lint:
    golangci-lint run

# Remove build artifacts
clean:
    rm -rf bin/

# Manage Go module dependencies
tidy:
    go mod tidy