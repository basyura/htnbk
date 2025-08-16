# Code Style and Conventions

## Go Code Style
- Follows standard Go formatting (`go fmt`)
- Package structure: `internal/` for private packages
- Clear separation of concerns:
  - `fetcher/` - API interaction
  - `models/` - Data structures
  - `storage/` - File operations

## Naming Conventions
- Functions: PascalCase for exported, camelCase for private
- Variables: camelCase
- Constants: ALL_CAPS or PascalCase
- Files: lowercase with underscores

## Code Organization
- Single `main.go` entry point
- Business logic in `doMain()` function
- Clean error handling with descriptive Japanese messages
- XML struct tags for API parsing
- Time handling using RFC3339 format

## Documentation
- Comments in Japanese for user-facing messages
- English for code comments and documentation
- README in Japanese for end users
- CLAUDE.md in English for development

## File Structure
- Use `internal/` for implementation packages
- Group related functionality in packages
- Keep models simple with XML tags
- Separate concerns cleanly