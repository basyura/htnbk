# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is "htnblg-export" (Hatena Blog Export) - a Go tool for extracting articles from Hatena Blog and outputting each entry to separate files for backup purposes.

## Development Commands

### Building and Running
- `go run main.go` - Run the application directly
- `go build` - Build the binary
- `go build -o htnblg-export` - Build with specific binary name

### Testing and Quality
- `go test ./...` - Run all tests
- `go vet ./...` - Run Go's static analysis tool
- `go fmt ./...` - Format all Go files
- `go build` - Build to verify code compiles (ALWAYS run after changes)

### Dependencies
- `go mod tidy` - Clean up module dependencies
- `go mod download` - Download dependencies

## Architecture

This is a simple Go command-line tool with a single main.go file. The tool uses Hatena Blog's AtomPub API to fetch blog entries and display them in a formatted list.

### Features
- Fetches blog entries from Hatena Blog using AtomPub API
- Uses Basic authentication with API key
- Displays entries in format: "連番 : 公開日 - タイトル"
- Parses XML response and extracts entry metadata

### Configuration
- Requires `HATENA_API_KEY` environment variable
- Currently configured for:
  - Hatena ID: basyura
  - Blog ID: blog.basyura.org

### API Reference
- Hatena Blog AtomPub API documentation: https://developer.hatena.ne.jp/ja/documents/blog/apis/atom/

## Project Structure

- `main.go` - Main entry point
- `go.mod` - Go module definition (Go 1.24.5)
- `README.md` - Project documentation in Japanese