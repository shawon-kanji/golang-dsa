# Go CLI Commands Cheatsheet

## Module Management

```bash
# Initialize module
go mod init github.com/user/project

# Add dependencies
go get package@version
go get package@latest
go get package@v1.2.3

# Update dependencies
go get -u              # Update all
go get -u ./...        # Update all in project
go get -u package      # Update specific package

# Remove unused dependencies
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify

# View module graph
go mod graph

# Vendor dependencies
go mod vendor

# View available versions
go list -m -versions package
```

## Building

```bash
# Build current package
go build

# Build specific file
go build main.go

# Build with output name
go build -o myapp

# Build for different OS/arch
GOOS=linux GOARCH=amd64 go build
GOOS=windows GOARCH=amd64 go build
GOOS=darwin GOARCH=arm64 go build

# Build with flags
go build -ldflags="-X main.version=1.0.0"
go build -tags=prod

# Install binary to $GOPATH/bin
go install

# Cross-compilation targets
GOOS=linux GOARCH=amd64
GOOS=windows GOARCH=amd64
GOOS=darwin GOARCH=amd64
GOOS=darwin GOARCH=arm64
```

## Running

```bash
# Run current package
go run .

# Run specific file
go run main.go

# Run with arguments
go run main.go arg1 arg2

# Run with build flags
go run -tags=dev main.go
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests in current package
go test

# Run specific test
go test -run TestName

# Run with verbose output
go test -v

# Run with coverage
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=.
go test -bench=BenchmarkName
go test -benchmem

# Run with race detection
go test -race

# Short mode (skip long tests)
go test -short

# Timeout
go test -timeout 30s

# Parallel execution
go test -parallel 4

# CPU profiling
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof
go tool pprof mem.prof
```

## Code Quality

```bash
# Format code
go fmt ./...
gofmt -w .

# Check for errors
go vet ./...

# Static analysis
go vet -all ./...

# Lint (requires golangci-lint)
golangci-lint run

# Check imports
goimports -w .
```

## Documentation

```bash
# View package documentation
go doc package
go doc package.Function

# Start documentation server
godoc -http=:6060

# Generate documentation
go doc -all > docs.txt
```

## Profiling

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Profile running program
go tool pprof http://localhost:6060/debug/pprof/profile

# Heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# View profile in browser
go tool pprof -http=:8080 cpu.prof

# Trace
go test -trace=trace.out
go tool trace trace.out
```

## Debugging

```bash
# Print AST
go build -gcflags="-m"

# Print optimizations
go build -gcflags="-m -m"

# Escape analysis
go build -gcflags="-m" 2>&1 | grep escape

# Assembly output
go tool compile -S file.go

# Build with debug info
go build -gcflags="all=-N -l"

# Use delve debugger
dlv debug
dlv test
dlv attach <pid>
```

## Package Management

```bash
# List packages
go list ./...
go list -m all        # List all modules
go list -u -m all     # List upgradable modules

# Show package info
go list -json package

# Clean build cache
go clean -cache
go clean -testcache
go clean -modcache

# Environment info
go env
go env GOPATH
go env GOOS GOARCH
```

## Code Generation

```bash
# Generate code
go generate ./...

# Generate with specific tool
//go:generate command args
```

## Workspace

```bash
# Initialize workspace
go work init

# Add module to workspace
go work use ./module

# Sync workspace
go work sync
```

## Common Flags

```bash
# Verbose output
-v

# Show commands being run
-x

# Number of parallel builds
-p n

# Build tags
-tags tag1,tag2

# Race detector
-race

# Memory sanitizer
-msan

# Work directory
-work

# Keep temporary files
-trimpath=false
```

## Environment Variables

```bash
# Go path
export GOPATH=$HOME/go

# Go binary path
export GOBIN=$GOPATH/bin

# Module proxy
export GOPROXY=https://proxy.golang.org,direct

# Private modules
export GOPRIVATE=github.com/private/*

# Module checksum database
export GOSUMDB=sum.golang.org

# Disable CGO
export CGO_ENABLED=0

# Number of CPUs
export GOMAXPROCS=4
```

## Useful Commands

```bash
# Update Go
go install golang.org/dl/go1.21.0@latest
go1.21.0 download

# Install tools
go install github.com/user/tool@latest

# Show version
go version

# Show build info
go version -m binary

# Bug report
go bug
```

## Development Workflow

```bash
# 1. Initialize project
go mod init project

# 2. Add dependencies
go get package

# 3. Format code
go fmt ./...

# 4. Check for issues
go vet ./...

# 5. Run tests
go test ./...

# 6. Build
go build -o myapp

# 7. Clean up
go mod tidy
```
