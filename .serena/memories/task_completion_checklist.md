# Task Completion Checklist

## CRITICAL: Always run after making any code changes

### 1. Code Quality Checks
```bash
go fmt ./...      # Format code
go vet ./...      # Static analysis
go build          # MUST verify compilation
```

### 2. Testing
```bash
go test ./...     # Run all tests
```

### 3. Dependencies
```bash
go mod tidy       # Clean up dependencies if added/removed
```

### 4. Final Verification
- Code compiles without errors
- All tests pass
- Code is properly formatted
- No vet warnings

## Important Notes
- **NEVER commit code that doesn't compile**
- Always test incremental and full fetch modes if modifying core logic
- Verify output file format if changing storage logic
- Test with actual Hatena Blog API if modifying fetcher

## User Instructions (CLAUDE.md)
- Never include CLAUDE mentions in commit logs
- Run `git status` when user types just "s"
- Show full results without abbreviation