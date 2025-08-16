# Suggested Commands for htnblg-export

## Development Commands

### Building and Running
- `go run main.go` - Run the application directly
- `go build` - Build the binary
- `go build -o htnblg-export` - Build with specific binary name

### Testing and Quality (ALWAYS run after changes)
- `go test ./...` - Run all tests
- `go vet ./...` - Run Go's static analysis tool
- `go fmt ./...` - Format all Go files
- `go build` - Build to verify code compiles (CRITICAL: run after any code changes)

### Dependencies
- `go mod tidy` - Clean up module dependencies
- `go mod download` - Download dependencies

### Running the Tool
```bash
# Incremental fetch (default)
./htnblg-export <hatenaID> <blogID> <apiKey>

# Full fetch
./htnblg-export --all <hatenaID> <blogID> <apiKey>

# Example
./htnblg-export basyura blog.basyura.org your_api_key
```

### System Commands (Darwin)
- `ls` - List files
- `find` - Search files
- `grep` - Search content
- `git` - Version control

## Important Notes
- Always run `go build` after making code changes to verify compilation
- Use `go fmt ./...` to maintain consistent code formatting
- The tool creates `entries/YYYY/MM/` directory structure automatically