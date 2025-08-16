# Project Overview: htnblg-export

## Purpose
A Hatena Blog Export Tool written in Go that extracts articles from Hatena Blog and outputs each entry to separate files for backup purposes.

## Tech Stack
- **Language**: Go 1.24.5
- **Architecture**: Simple command-line tool with clean separation of concerns
- **API**: Hatena Blog AtomPub API for fetching blog entries
- **Output Format**: Markdown files with YAML frontmatter

## Key Features
- Incremental fetching (only new articles since last run)
- Full backup mode with --all option
- Automatic directory organization by year/month
- Pagination support for large blogs
- File naming: YYYY-MM-DD_title.md

## Project Structure
```
├── main.go                     # Main entry point
├── go.mod                      # Go module (1.24.5)
├── internal/
│   ├── fetcher/fetcher.go      # API fetching logic
│   ├── models/                 # Data structures
│   │   ├── entry.go           # Blog entry model
│   │   ├── feed.go            # Feed model
│   │   └── link.go            # Link model
│   └── storage/storage.go      # File I/O operations
├── README.md                   # Japanese documentation
└── CLAUDE.md                   # Development instructions
```

## API Integration
- Uses Hatena Blog AtomPub API
- Basic authentication with API key
- XML parsing for entry data
- Link extraction for entry URLs